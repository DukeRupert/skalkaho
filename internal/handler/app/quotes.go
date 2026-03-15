package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/service/email"
)

// QuoteListData holds data for the quote list partial.
type QuoteListData struct {
	ProjectID string
	Quotes    []QuoteView
}

// QuoteView is a template-friendly view of a quote.
type QuoteView struct {
	ID        string
	Version   int64
	Status    string
	Token     string
	SentAt    string
	CreatedAt string
	ExpiresAt string
	Signed    bool
	HasNotes  bool
}

// SendModalData holds data for the send quote modal.
type SendModalData struct {
	QuoteID  string
	QuoteURL string
	Status   string
	Email    string // pre-filled from client
}

// NotesModalData holds data for the quote notes modal.
type NotesModalData struct {
	QuoteID   string
	ProjectID string
	Version   int64
	Notes     string
}

// ListQuotes returns the quote version list partial for the overview page.
func (h *Handler) ListQuotes(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	quotes, err := h.queries.ListQuotesByProject(r.Context(), projectID)
	if err != nil {
		h.logger.Error("listing quotes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := QuoteListData{
		ProjectID: projectID,
		Quotes:    toQuoteViews(quotes),
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "quote-list", data); err != nil {
		h.logger.Error("rendering quote list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateQuote creates a new quote by snapshotting the current estimate.
func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	ctx := r.Context()

	project, err := h.queries.GetProject(ctx, projectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Get next version number
	latestVersion, err := h.queries.GetLatestQuoteVersion(ctx, projectID)
	if err != nil {
		h.logger.Error("getting latest version", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newVersion := latestVersion + 1

	// Supersede any active quotes
	if err := h.queries.SupersedeActiveQuotes(ctx, projectID); err != nil {
		h.logger.Error("superseding quotes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Build estimate snapshot
	sections, err := h.loadEstimateSections(r, projectID)
	if err != nil {
		h.logger.Error("loading estimate", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	globals := domain.MarkupGlobals{
		MaterialsMarkup: project.MaterialsMarkup,
		LaborMarkup:     project.LaborMarkup,
		EquipmentMarkup: project.EquipmentMarkup,
		SubsMarkup:      project.SubsMarkup,
		OtherMarkup:     project.OtherMarkup,
	}

	costSummary := domain.CalculateProjectCosts(sections, globals)

	estimateJSON, err := json.Marshal(sections)
	if err != nil {
		h.logger.Error("marshaling estimate", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalsJSON, err := json.Marshal(costSummary)
	if err != nil {
		h.logger.Error("marshaling totals", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()[:20]
	if _, err := h.queries.CreateQuote(ctx, repository.CreateQuoteParams{
		ID:               id,
		ProjectID:        projectID,
		Version:          newVersion,
		Status:           "draft",
		EstimateSnapshot: pqtype.NullRawMessage{RawMessage: estimateJSON, Valid: true},
		TotalsSnapshot:   pqtype.NullRawMessage{RawMessage: totalsJSON, Valid: true},
	}); err != nil {
		h.logger.Error("creating quote", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Move project to "In Review" when a quote is generated
	if err := h.queries.UpdateProjectStatus(ctx, repository.UpdateProjectStatusParams{
		ID:     projectID,
		Status: "In Review",
	}); err != nil {
		h.logger.Error("updating project status", "error", err)
	}

	// Return updated quote list
	h.ListQuotes(w, r)
}

// SendQuote generates a token and transitions the quote to 'sent'.
func (h *Handler) SendQuote(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	token, err := generateToken(10)
	if err != nil {
		h.logger.Error("generating token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour) // 30 days

	if err := h.queries.UpdateQuoteSent(ctx, repository.UpdateQuoteSentParams{
		ID:        quoteID,
		Token:     sql.NullString{String: token, Valid: true},
		SentAt:    sql.NullTime{Time: now, Valid: true},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: true},
	}); err != nil {
		h.logger.Error("updating quote sent", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Try to send email if configured and email provided
	if err := r.ParseForm(); err == nil {
		recipientEmail := r.FormValue("email")
		if recipientEmail != "" && h.emailClient != nil && h.emailClient.Enabled() {
			h.sendQuoteEmail(r, quoteID, token, recipientEmail)
		}
	}

	// Return the send modal with the link
	quoteURL := fmt.Sprintf("/q/%s", token)
	data := SendModalData{
		QuoteID:  quoteID,
		QuoteURL: quoteURL,
		Status:   "sent",
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "send-modal-result", data); err != nil {
		h.logger.Error("rendering send modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ResendQuote resets sent_at and expires_at without creating a new version.
func (h *Handler) ResendQuote(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)

	if err := h.queries.UpdateQuoteSent(ctx, repository.UpdateQuoteSentParams{
		ID:        quoteID,
		Token:     quote.Token,
		SentAt:    sql.NullTime{Time: now, Valid: true},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: true},
	}); err != nil {
		h.logger.Error("resending quote", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Try to send email if provided
	if err := r.ParseForm(); err == nil {
		recipientEmail := r.FormValue("email")
		if recipientEmail != "" && h.emailClient != nil && h.emailClient.Enabled() {
			h.sendQuoteEmail(r, quoteID, quote.Token.String, recipientEmail)
		}
	}

	quoteURL := fmt.Sprintf("/q/%s", quote.Token.String)
	data := SendModalData{
		QuoteID:  quoteID,
		QuoteURL: quoteURL,
		Status:   "sent",
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "send-modal-result", data); err != nil {
		h.logger.Error("rendering send modal result", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetSendModal returns the modal for sending/sharing a quote.
func (h *Handler) GetSendModal(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	// Try to pre-fill client email
	var clientEmail string
	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err == nil && project.ClientID.Valid {
		client, err := h.queries.GetClient(ctx, project.ClientID.String)
		if err == nil && client.Email.Valid {
			clientEmail = client.Email.String
		}
	}

	quoteURL := ""
	if quote.Token.Valid {
		quoteURL = fmt.Sprintf("/q/%s", quote.Token.String)
	}

	data := SendModalData{
		QuoteID:  quoteID,
		QuoteURL: quoteURL,
		Status:   quote.Status,
		Email:    clientEmail,
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "send-modal", data); err != nil {
		h.logger.Error("rendering send modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handler) sendQuoteEmail(r *http.Request, quoteID, token, recipientEmail string) {
	ctx := r.Context()

	// Get company name
	companyName := "Contractor"
	profile, err := h.queries.GetCompanyProfile(ctx)
	if err == nil {
		companyName = profile.Name
	}

	// Get project name
	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote for email", "error", err)
		return
	}
	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("getting project for email", "error", err)
		return
	}

	// Build full URL (use request host)
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	quoteURL := fmt.Sprintf("%s://%s/q/%s", scheme, r.Host, token)

	messageID, err := h.emailClient.SendQuoteEmail(recipientEmail, quoteURL, project.Name, companyName)
	if err != nil {
		h.logger.Error("sending quote email", "error", err, "to", recipientEmail)
		return
	}

	// Record email sent
	emailID := uuid.New().String()[:20]
	if _, err := h.queries.CreateQuoteEmail(ctx, repository.CreateQuoteEmailParams{
		ID:         emailID,
		QuoteID:    quoteID,
		Recipient:  recipientEmail,
		ProviderID: sql.NullString{String: messageID, Valid: messageID != ""},
	}); err != nil {
		h.logger.Error("recording quote email", "error", err)
	}
}

func toQuoteViews(quotes []repository.Quote) []QuoteView {
	views := make([]QuoteView, 0, len(quotes))
	for _, q := range quotes {
		view := QuoteView{
			ID:        q.ID,
			Version:   q.Version,
			Status:    q.Status,
			CreatedAt: q.CreatedAt.Format("Jan 2, 2006"),
		}
		if q.Token.Valid {
			view.Token = q.Token.String
		}
		if q.SentAt.Valid {
			view.SentAt = q.SentAt.Time.Format("Jan 2, 2006")
		}
		if q.ExpiresAt.Valid {
			view.ExpiresAt = q.ExpiresAt.Time.Format("Jan 2, 2006")
		}
		if q.Status == "signed" {
			view.Signed = true
		}
		if q.Notes != "" {
			view.HasNotes = true
		}
		views = append(views, view)
	}
	return views
}

// GetNotesModal returns the modal for viewing/editing quote notes.
func (h *Handler) GetNotesModal(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")

	quote, err := h.queries.GetQuote(r.Context(), quoteID)
	if err != nil {
		h.logger.Error("getting quote for notes", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	data := NotesModalData{
		QuoteID:   quoteID,
		ProjectID: quote.ProjectID,
		Version:   quote.Version,
		Notes:     quote.Notes,
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "notes-modal", data); err != nil {
		h.logger.Error("rendering notes modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateQuoteNotes saves notes on a quote.
func (h *Handler) UpdateQuoteNotes(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	notes := r.FormValue("notes")
	if err := h.queries.UpdateQuoteNotes(r.Context(), repository.UpdateQuoteNotesParams{
		ID:    quoteID,
		Notes: notes,
	}); err != nil {
		h.logger.Error("updating quote notes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the project ID to re-render the quote list
	quote, err := h.queries.GetQuote(r.Context(), quoteID)
	if err != nil {
		h.logger.Error("getting quote after notes update", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return updated quote list by setting the path value and delegating
	// We need to render the notes-saved confirmation modal
	data := NotesModalData{
		QuoteID:   quoteID,
		ProjectID: quote.ProjectID,
		Version:   quote.Version,
		Notes:     notes,
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "notes-modal-saved", data); err != nil {
		h.logger.Error("rendering notes saved", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func generateToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}

// SetEmailClient sets the email client for quote sending.
func (h *Handler) SetEmailClient(client *email.Client) {
	h.emailClient = client
}
