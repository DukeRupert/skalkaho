package keyboard

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// GetSendSignatureForm shows the form for sending an estimate for signature.
func (h *Handler) GetSendSignatureForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	estimateID := r.PathValue("id")

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    estimate.JobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Get client info to pre-fill form
	var client *repository.Client
	if job.ClientID.Valid {
		c, err := h.queries.GetClient(ctx, repository.GetClientParams{
			ID:    job.ClientID.String,
			OrgID: orgID,
		})
		if err == nil {
			client = &c
		}
	}

	// Get company profile
	company, _ := h.queries.GetCompanyProfile(ctx, orgID.UUID)

	data := map[string]any{
		"Estimate": estimate,
		"Job":      job,
		"Client":   client,
		"Company":  company,
	}

	if err := h.renderer.Render(w, "signature_send", data); err != nil {
		logger.Error("failed to render signature send form", "error", err)
	}
}

// SendForSignature creates a signature request and sends the signing link.
func (h *Handler) SendForSignature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	estimateID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Validate input
	input := domain.SignatureRequestInput{
		EstimateID:     estimateID,
		RecipientEmail: strings.TrimSpace(r.FormValue("recipient_email")),
		RecipientName:  strings.TrimSpace(r.FormValue("recipient_name")),
		Message:        strings.TrimSpace(r.FormValue("message")),
	}

	if errs := input.Validate(); len(errs) > 0 {
		http.Error(w, errs[0].Message, http.StatusBadRequest)
		return
	}

	// Get estimate
	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	// Check for existing pending request
	existingRequest, err := h.queries.GetPendingSignatureRequestByEstimate(ctx, repository.GetPendingSignatureRequestByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	if err == nil && existingRequest.ID != "" {
		http.Error(w, "This estimate already has a pending signature request. Cancel it first.", http.StatusBadRequest)
		return
	}

	// Get job and client for snapshot
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    estimate.JobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Get estimate categories for snapshot
	categories, err := h.queries.ListEstimateCategoriesByEstimate(ctx, repository.ListEstimateCategoriesByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	if err != nil {
		logger.Error("failed to get estimate categories", "error", err)
		http.Error(w, "Failed to load estimate categories", http.StatusInternalServerError)
		return
	}

	// Build category snapshots
	categorySnapshots := make([]domain.CategorySnapshot, len(categories))
	for i, cat := range categories {
		var desc *string
		if cat.Description.Valid {
			desc = &cat.Description.String
		}
		categorySnapshots[i] = domain.CategorySnapshot{
			Name:        cat.Name,
			Description: desc,
			Total:       cat.Total,
			Tier:        int(cat.Tier),
			SortOrder:   int(cat.SortOrder),
		}
	}

	// Get notes
	var notes *string
	if estimate.Notes.Valid {
		notes = &estimate.Notes.String
	}

	// Create snapshot
	snapshot := domain.CreateQuoteSnapshot(
		estimateID,
		int(estimate.Version),
		job.Name,
		input.RecipientName,
		input.RecipientEmail,
		estimate.GrandTotal,
		notes,
		categorySnapshots,
	)

	// Serialize and hash
	snapshotBytes, err := domain.SerializeSnapshot(snapshot)
	if err != nil {
		logger.Error("failed to serialize snapshot", "error", err)
		http.Error(w, "Failed to create document snapshot", http.StatusInternalServerError)
		return
	}
	documentHash := domain.HashContent(snapshotBytes)

	// Generate token
	token, err := domain.GenerateSecureToken()
	if err != nil {
		logger.Error("failed to generate token", "error", err)
		http.Error(w, "Failed to generate signing link", http.StatusInternalServerError)
		return
	}

	// Calculate expiry
	expiresAt := domain.CalculateExpiryTime()

	// Create signature request
	var message sql.NullString
	if input.Message != "" {
		message = sql.NullString{String: input.Message, Valid: true}
	}

	var senderIP, senderUA sql.NullString
	if ip := r.RemoteAddr; ip != "" {
		senderIP = sql.NullString{String: ip, Valid: true}
	}
	if ua := r.UserAgent(); ua != "" {
		senderUA = sql.NullString{String: ua, Valid: true}
	}

	_, err = h.queries.CreateSignatureRequest(ctx, repository.CreateSignatureRequestParams{
		ID:              uuid.New().String(),
		OrgID:           orgID,
		EstimateID:      estimateID,
		RecipientEmail:  input.RecipientEmail,
		RecipientName:   input.RecipientName,
		Token:           token,
		DocumentHash:    documentHash,
		QuoteSnapshot:   string(snapshotBytes),
		Message:         message,
		Status:          string(domain.SignatureRequestStatusPending),
		ExpiresAt:       expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		SenderIp:        senderIP,
		SenderUserAgent: senderUA,
	})
	if err != nil {
		logger.Error("failed to create signature request", "error", err)
		http.Error(w, "Failed to create signature request", http.StatusInternalServerError)
		return
	}

	// Update estimate status to 'sent'
	_, err = h.queries.MarkEstimateSent(ctx, repository.MarkEstimateSentParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to update estimate status", "error", err)
		// Non-fatal: signature request was created, just log the error
	}

	// Build signing URL
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	signingURL := fmt.Sprintf("%s://%s/sign/%s", scheme, r.Host, token)

	// Log email stub (for development)
	logger.Info("EMAIL STUB: Would send signature request",
		"estimate_id", estimateID,
		"recipient_email", input.RecipientEmail,
		"recipient_name", input.RecipientName,
		"signing_url", signingURL,
		"expires_at", expiresAt,
	)

	// For development: show the signing URL in response
	if h.config.Environment == "development" {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<div class="p-4 bg-green-50 border border-green-200 rounded-lg">
			<p class="font-medium text-green-800">Signing link generated!</p>
			<p class="text-sm text-green-700 mt-2">In production, this would be emailed to %s</p>
			<p class="text-sm text-green-700 mt-1">Development signing URL:</p>
			<a href="%s" target="_blank" class="text-copper-600 hover:text-copper-700 underline break-all">%s</a>
		</div>`, input.RecipientEmail, signingURL, signingURL)
		return
	}

	// Redirect back to estimate
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/estimates/"+estimateID)
		return
	}
	http.Redirect(w, r, "/estimates/"+estimateID, http.StatusSeeOther)
}

// GetSignaturePage shows the public signing page for a signature request.
func (h *Handler) GetSignaturePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	token := r.PathValue("token")

	// Look up request by token
	request, err := h.queries.GetSignatureRequestByToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			h.renderSignatureError(w, "Invalid signing link", "This signing link is invalid or has expired.")
			return
		}
		logger.Error("failed to get signature request", "error", err)
		h.renderSignatureError(w, "Error", "Failed to load signing page.")
		return
	}

	// Check if expired
	if request.Status == string(domain.SignatureRequestStatusExpired) {
		h.renderSignatureError(w, "Link Expired", "This signing link has expired. Please request a new one from the sender.")
		return
	}

	// Check if cancelled
	if request.Status == string(domain.SignatureRequestStatusCancelled) {
		h.renderSignatureError(w, "Link Cancelled", "This signing link has been cancelled.")
		return
	}

	// Check if already signed
	if request.Status == string(domain.SignatureRequestStatusSigned) {
		// Redirect to complete page
		http.Redirect(w, r, "/sign/"+token+"/complete", http.StatusSeeOther)
		return
	}

	// Deserialize snapshot
	snapshot, err := domain.DeserializeSnapshot([]byte(request.QuoteSnapshot))
	if err != nil {
		logger.Error("failed to deserialize snapshot", "error", err)
		h.renderSignatureError(w, "Error", "Failed to load document.")
		return
	}

	// Get company profile
	company, _ := h.queries.GetCompanyProfile(ctx, orgID.UUID)

	data := map[string]any{
		"Request":     request,
		"Snapshot":    snapshot,
		"Company":     company,
		"Token":       token,
		"ConsentText": domain.ConsentText,
	}

	if err := h.renderer.Render(w, "signature_page", data); err != nil {
		logger.Error("failed to render signature page", "error", err)
	}
}

// SubmitSignature handles the signature submission.
func (h *Handler) SubmitSignature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	token := r.PathValue("token")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Validate input
	input := domain.SignatureInput{
		LegalName: strings.TrimSpace(r.FormValue("legal_name")),
		Agreed:    r.FormValue("agreed") == "on" || r.FormValue("agreed") == "true",
	}

	if errs := input.Validate(); len(errs) > 0 {
		http.Error(w, errs[0].Message, http.StatusBadRequest)
		return
	}

	// Look up request by token
	request, err := h.queries.GetSignatureRequestByToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid signing link", http.StatusNotFound)
			return
		}
		logger.Error("failed to get signature request", "error", err)
		http.Error(w, "Failed to process signature", http.StatusInternalServerError)
		return
	}

	// Verify request is still pending
	if request.Status != string(domain.SignatureRequestStatusPending) {
		http.Error(w, "This signature request is no longer valid", http.StatusBadRequest)
		return
	}

	// Create signature record
	_, err = h.queries.CreateSignature(ctx, repository.CreateSignatureParams{
		ID:              uuid.New().String(),
		OrgID:           orgID,
		RequestID:       request.ID,
		LegalName:       input.LegalName,
		ConsentText:     domain.ConsentText,
		DocumentHash:    request.DocumentHash,
		SignerIp:        r.RemoteAddr,
		SignerUserAgent: r.UserAgent(),
		SignerEmail:     request.RecipientEmail,
	})
	if err != nil {
		logger.Error("failed to create signature", "error", err)
		http.Error(w, "Failed to save signature", http.StatusInternalServerError)
		return
	}

	// Update request status to signed
	_, err = h.queries.UpdateSignatureRequestStatus(ctx, repository.UpdateSignatureRequestStatusParams{
		Status: string(domain.SignatureRequestStatusSigned),
		ID:     request.ID,
		OrgID:  orgID,
	})
	if err != nil {
		logger.Error("failed to update signature request status", "error", err)
	}

	// Update estimate status to 'accepted'
	_, err = h.queries.MarkEstimateAccepted(ctx, repository.MarkEstimateAcceptedParams{
		ID:    request.EstimateID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to update estimate status to accepted", "error", err)
		// Non-fatal: signature was captured, just log the error
	}

	// Log email notifications (stub)
	logger.Info("EMAIL STUB: Would send signature confirmation",
		"signer_email", request.RecipientEmail,
		"signer_name", input.LegalName,
	)

	// Redirect to confirmation page
	http.Redirect(w, r, "/sign/"+token+"/complete", http.StatusSeeOther)
}

// GetSignatureComplete shows the signature confirmation page.
func (h *Handler) GetSignatureComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	token := r.PathValue("token")

	// Look up request by token
	request, err := h.queries.GetSignatureRequestByToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			h.renderSignatureError(w, "Invalid Link", "This signing link is invalid.")
			return
		}
		logger.Error("failed to get signature request", "error", err)
		h.renderSignatureError(w, "Error", "Failed to load confirmation.")
		return
	}

	// Must be signed
	if request.Status != string(domain.SignatureRequestStatusSigned) {
		http.Redirect(w, r, "/sign/"+token, http.StatusSeeOther)
		return
	}

	// Get signature details
	signature, err := h.queries.GetSignatureByRequest(ctx, repository.GetSignatureByRequestParams{
		RequestID: request.ID,
		OrgID:     orgID,
	})
	if err != nil {
		logger.Error("failed to get signature", "error", err)
		h.renderSignatureError(w, "Error", "Failed to load signature details.")
		return
	}

	// Deserialize snapshot
	snapshot, err := domain.DeserializeSnapshot([]byte(request.QuoteSnapshot))
	if err != nil {
		logger.Error("failed to deserialize snapshot", "error", err)
	}

	// Get company profile
	company, _ := h.queries.GetCompanyProfile(ctx, orgID.UUID)

	data := map[string]any{
		"Request":   request,
		"Signature": signature,
		"Snapshot":  snapshot,
		"Company":   company,
		"Token":     token,
	}

	if err := h.renderer.Render(w, "signature_complete", data); err != nil {
		logger.Error("failed to render signature complete", "error", err)
	}
}

// CancelSignatureRequest cancels a pending signature request.
func (h *Handler) CancelSignatureRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	estimateID := r.PathValue("id")

	// Get the pending request
	request, err := h.queries.GetPendingSignatureRequestByEstimate(ctx, repository.GetPendingSignatureRequestByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No pending signature request found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get signature request", "error", err)
		http.Error(w, "Failed to cancel request", http.StatusInternalServerError)
		return
	}

	// Cancel it
	_, err = h.queries.CancelSignatureRequest(ctx, repository.CancelSignatureRequestParams{
		ID:    request.ID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to cancel signature request", "error", err)
		http.Error(w, "Failed to cancel request", http.StatusInternalServerError)
		return
	}

	// Redirect back to estimate
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/estimates/"+estimateID)
		return
	}
	http.Redirect(w, r, "/estimates/"+estimateID, http.StatusSeeOther)
}

// renderSignatureError renders a friendly error page for public signature endpoints.
func (h *Handler) renderSignatureError(w http.ResponseWriter, title, message string) {
	data := map[string]any{
		"Title":   title,
		"Message": message,
	}
	if err := h.renderer.Render(w, "signature_error", data); err != nil {
		http.Error(w, message, http.StatusBadRequest)
	}
}
