package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// PublicQuoteData holds data for the public quote page.
type PublicQuoteData struct {
	// Quote info
	QuoteID   string
	Version   int64
	Status    string // "active", "signed", "expired", "not_found"
	SentAt    string
	ExpiresAt string

	// Project + client
	ProjectName string
	ClientName  string

	// Company (contractor)
	CompanyName string

	// Cost breakdown from snapshot
	CostSummary domain.ProjectCostSummary

	// Signature (if signed)
	SignerName string
	SignedAt   string

	// Error message
	Error string
}

// GetQuotePage renders the public quote page at /q/{token}.
func (h *Handler) GetQuotePage(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	row, err := h.queries.GetQuoteByToken(r.Context(), sql.NullString{String: token, Valid: true})
	if err != nil {
		data := PublicQuoteData{Status: "not_found"}
		if renderErr := h.renderer.Render(w, "quote_public.html", data); renderErr != nil {
			h.logger.Error("rendering not found quote", "error", renderErr)
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}

	// Load company profile
	companyName := "Contractor"
	profile, err := h.queries.GetCompanyProfile(r.Context())
	if err == nil {
		companyName = profile.Name
	}

	// Parse totals snapshot
	var costSummary domain.ProjectCostSummary
	if row.TotalsSnapshot.Valid {
		_ = json.Unmarshal(row.TotalsSnapshot.RawMessage, &costSummary)
	}

	// Determine display status
	status := "active"
	switch row.Status {
	case "signed":
		status = "signed"
	case "expired", "superseded":
		status = "expired"
	case "sent":
		if row.ExpiresAt.Valid && time.Now().After(row.ExpiresAt.Time) {
			status = "expired"
		}
	case "draft":
		status = "expired" // draft shouldn't be accessible via token, treat as expired
	}

	data := PublicQuoteData{
		QuoteID:     row.ID,
		Version:     row.Version,
		Status:      status,
		ProjectName: row.ProjectName,
		CompanyName: companyName,
		CostSummary: costSummary,
	}

	if row.ProjectClientName.Valid {
		data.ClientName = row.ProjectClientName.String
	}
	if row.ClientCompany.Valid {
		data.ClientName = row.ClientCompany.String
	}
	if row.SentAt.Valid {
		data.SentAt = row.SentAt.Time.Format("January 2, 2006")
	}
	if row.ExpiresAt.Valid {
		data.ExpiresAt = row.ExpiresAt.Time.Format("January 2, 2006")
	}

	// Load signature if signed
	if status == "signed" {
		sig, err := h.queries.GetQuoteSignature(r.Context(), row.ID)
		if err == nil {
			data.SignerName = sig.SignerName
			data.SignedAt = sig.SignedAt.Format("January 2, 2006 at 3:04 PM")
		}
	}

	if err := h.renderer.Render(w, "quote_public.html", data); err != nil {
		h.logger.Error("rendering public quote", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// SubmitSignature handles POST /q/{token} — records the signature.
func (h *Handler) SubmitSignature(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	signerName := r.FormValue("signer_name")
	if signerName == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	row, err := h.queries.GetQuoteByToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	// Validate quote can be signed
	if row.Status != "sent" {
		http.Error(w, "Quote cannot be signed", http.StatusBadRequest)
		return
	}
	if row.ExpiresAt.Valid && time.Now().After(row.ExpiresAt.Time) {
		http.Error(w, "Quote has expired", http.StatusBadRequest)
		return
	}

	// Record signature
	sigID := uuid.New().String()[:20]
	signerIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		signerIP = forwarded
	}

	if _, err := h.queries.CreateQuoteSignature(ctx, repository.CreateQuoteSignatureParams{
		ID:         sigID,
		QuoteID:    row.ID,
		SignerName: signerName,
		SignerIp:   sql.NullString{String: signerIP, Valid: true},
	}); err != nil {
		h.logger.Error("creating signature", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update quote status to signed
	if err := h.queries.UpdateQuoteStatus(ctx, repository.UpdateQuoteStatusParams{
		ID:     row.ID,
		Status: "signed",
	}); err != nil {
		h.logger.Error("updating quote status", "error", err)
	}

	// Update project status to "In Review"
	if err := h.queries.UpdateProjectStatus(ctx, repository.UpdateProjectStatusParams{
		ID:     row.ProjectID,
		Status: "In Review",
	}); err != nil {
		h.logger.Error("updating project status", "error", err)
	}

	// Redirect back to the same page to show signed state
	http.Redirect(w, r, fmt.Sprintf("/q/%s", token), http.StatusSeeOther)
}
