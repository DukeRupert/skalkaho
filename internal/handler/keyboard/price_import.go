package keyboard

import (
	"bytes"
	"context"
	"crypto/subtle"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/service/excel"
	"github.com/google/uuid"
)

const priceImportCookieName = "price_import_auth"

// checkPriceImportAuth checks if the user has valid authentication for price import.
func (h *Handler) checkPriceImportAuth(r *http.Request) bool {
	// If no token is configured, allow access (for development)
	if h.config.PriceImportToken == "" {
		return true
	}

	cookie, err := r.Cookie(priceImportCookieName)
	if err != nil {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(h.config.PriceImportToken)) == 1
}

// GetPriceImportPage renders the price import upload page.
func (h *Handler) GetPriceImportPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Check if token authentication is required and valid
	requiresToken := h.config.PriceImportToken != ""
	isAuthenticated := h.checkPriceImportAuth(r)

	// Check if Claude API is configured
	hasAPI := h.matcher != nil

	// Get list of imports
	imports, err := h.queries.ListPriceImports(ctx, repository.ListPriceImportsParams{
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		logger.Error("failed to list imports", "error", err)
		imports = []repository.PriceImport{}
	}

	// Check if any imports are still processing (for auto-refresh)
	hasProcessing := false
	for _, imp := range imports {
		if imp.Status == "processing" {
			hasProcessing = true
			break
		}
	}

	// Check for success message
	successCount := r.URL.Query().Get("success")

	data := map[string]interface{}{
		"HasClaudeAPI":    hasAPI,
		"RequiresToken":   requiresToken,
		"IsAuthenticated": isAuthenticated,
		"Imports":         imports,
		"HasProcessing":   hasProcessing,
		"SuccessCount":    successCount,
	}

	if err := h.renderer.Render(w, "price_import", data); err != nil {
		logger.Error("failed to render price import page", "error", err)
	}
}

// ValidatePriceImportToken validates the token and sets auth cookie.
func (h *Handler) ValidatePriceImportToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")

	// Validate token using constant-time comparison
	if subtle.ConstantTimeCompare([]byte(token), []byte(h.config.PriceImportToken)) != 1 {
		logger.Warn("invalid price import token attempt")
		// Return the page with error
		data := map[string]interface{}{
			"HasClaudeAPI":    h.matcher != nil,
			"RequiresToken":   true,
			"IsAuthenticated": false,
			"TokenError":      "Invalid token. Please try again.",
		}
		if err := h.renderer.Render(w, "price_import", data); err != nil {
			logger.Error("failed to render price import page", "error", err)
		}
		return
	}

	// Set authentication cookie (expires in 24 hours)
	http.SetCookie(w, &http.Cookie{
		Name:     priceImportCookieName,
		Value:    h.config.PriceImportToken,
		Path:     "/price-import",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	logger.Info("price import token validated successfully")

	// Redirect to show the upload form
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/price-import")
		return
	}
	http.Redirect(w, r, "/price-import", http.StatusSeeOther)
}

// UploadPriceFile handles Excel file upload.
func (h *Handler) UploadPriceFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Verify authentication
	if !h.checkPriceImportAuth(r) {
		logger.Warn("unauthorized price import upload attempt")
		http.Error(w, "Unauthorized. Please authenticate first.", http.StatusUnauthorized)
		return
	}

	// Check if Claude API is configured
	if h.matcher == nil {
		http.Error(w, "Claude API not configured. Set CLAUDE_API_KEY environment variable.", http.StatusServiceUnavailable)
		return
	}

	// Parse multipart form (10MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Error("failed to parse multipart form", "error", err)
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Error("no file uploaded", "error", err)
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		http.Error(w, "Invalid file type. Please upload .xlsx or .xls file", http.StatusBadRequest)
		return
	}

	// Read file into memory so we can process in background
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("failed to read file", "error", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	filename := header.Filename

	// Create import record immediately with "processing" status
	importID := uuid.New().String()
	_, err = h.queries.CreatePriceImport(ctx, repository.CreatePriceImportParams{
		ID:        importID,
		Filename:  filename,
		Status:    "processing",
		TotalRows: 0, // Will be updated after processing
	})
	if err != nil {
		logger.Error("failed to create import record", "error", err)
		http.Error(w, "Failed to create import", http.StatusInternalServerError)
		return
	}

	logger.Info("starting background price import processing", "import_id", importID, "filename", filename)

	// Process in background goroutine
	go h.processImportInBackground(importID, filename, fileBytes, logger)

	// Return immediately to the imports list page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/price-import")
		return
	}
	http.Redirect(w, r, "/price-import", http.StatusSeeOther)
}

// processImportInBackground handles the Claude API call and match storage.
func (h *Handler) processImportInBackground(importID, filename string, fileBytes []byte, logger *slog.Logger) {
	// Use background context since the request context is gone
	ctx := context.Background()

	// Convert Excel file to text for Claude to parse
	parser := excel.NewParser()
	spreadsheet, err := parser.ParseToText(bytes.NewReader(fileBytes), filename)
	if err != nil {
		logger.Error("failed to parse excel file", "error", err, "import_id", importID)
		h.updateImportError(ctx, importID, "Failed to parse Excel file: "+err.Error())
		return
	}

	// Get all item templates for matching
	templates, err := h.queries.ListItemTemplates(ctx)
	if err != nil {
		logger.Error("failed to list templates", "error", err, "import_id", importID)
		h.updateImportError(ctx, importID, "Failed to load item templates")
		return
	}

	// Call Claude API to extract items and match them
	extractResult, err := h.matcher.ExtractAndMatchItems(ctx, spreadsheet, templates)
	if err != nil {
		logger.Error("failed to extract and match items with Claude", "error", err, "import_id", importID)
		h.updateImportError(ctx, importID, "AI extraction/matching failed: "+err.Error())
		return
	}

	// Store matches in database
	matchedCount := 0
	autoApproveThreshold := h.config.AutoApproveThreshold

	for _, item := range extractResult.Items {
		status := "pending"
		if item.Confidence >= autoApproveThreshold && item.TemplateID != nil {
			status = "auto_approved"
		}

		var templateID sql.NullInt64
		if item.TemplateID != nil {
			templateID = sql.NullInt64{Int64: *item.TemplateID, Valid: true}
		}

		var sourceUnit sql.NullString
		if item.Unit != "" {
			sourceUnit = sql.NullString{String: item.Unit, Valid: true}
		}

		var matchReason sql.NullString
		if item.Reason != "" {
			matchReason = sql.NullString{String: item.Reason, Valid: true}
		}

		_, err = h.queries.CreatePriceImportMatch(ctx, repository.CreatePriceImportMatchParams{
			ImportID:          importID,
			RowNumber:         int64(item.RowNumber),
			SourceName:        item.Name,
			SourceUnit:        sourceUnit,
			SourcePrice:       item.Price,
			MatchedTemplateID: templateID,
			Confidence:        item.Confidence,
			MatchReason:       matchReason,
			Status:            status,
		})
		if err != nil {
			logger.Error("failed to create match", "error", err, "row", item.RowNumber, "import_id", importID)
			continue
		}

		if item.TemplateID != nil {
			matchedCount++
		}
	}

	// Update import status to ready
	_, err = h.queries.UpdatePriceImportStatus(ctx, repository.UpdatePriceImportStatusParams{
		ID:          importID,
		Status:      "ready",
		MatchedRows: int64(matchedCount),
		TotalRows:   int64(len(extractResult.Items)),
	})
	if err != nil {
		logger.Error("failed to update import status", "error", err, "import_id", importID)
		return
	}

	logger.Info("completed price import processing", "import_id", importID, "total_items", len(extractResult.Items), "matched", matchedCount)
}

// updateImportError marks an import as failed with an error message.
func (h *Handler) updateImportError(ctx context.Context, importID string, errMsg string) {
	_, _ = h.queries.UpdatePriceImportStatus(ctx, repository.UpdatePriceImportStatusParams{
		ID:           importID,
		Status:       "failed",
		ErrorMessage: sql.NullString{String: errMsg, Valid: true},
		TotalRows:    0,
		MatchedRows:  0,
	})
}

// GetImportReview shows the review page for matched items.
func (h *Handler) GetImportReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	importID := r.PathValue("id")
	if importID == "" {
		http.Error(w, "Import ID required", http.StatusBadRequest)
		return
	}

	// Get import record
	priceImport, err := h.queries.GetPriceImport(ctx, importID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Import not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get import", "error", err)
		http.Error(w, "Failed to load import", http.StatusInternalServerError)
		return
	}

	// Get matches
	matches, err := h.queries.ListMatchesByImport(ctx, importID)
	if err != nil {
		logger.Error("failed to list matches", "error", err)
		http.Error(w, "Failed to load matches", http.StatusInternalServerError)
		return
	}

	// Get match counts by status
	statusCounts, err := h.queries.CountMatchesByStatus(ctx, importID)
	if err != nil {
		logger.Error("failed to count matches", "error", err)
	}

	counts := map[string]int64{
		"pending":       0,
		"approved":      0,
		"rejected":      0,
		"auto_approved": 0,
		"created":       0,
	}
	for _, sc := range statusCounts {
		counts[sc.Status] = sc.Count
	}

	// Count unmatched items (pending with no template)
	unmatched, err := h.queries.ListUnmatchedItems(ctx, importID)
	if err != nil {
		logger.Error("failed to count unmatched", "error", err)
	}
	unmatchedCount := int64(len(unmatched))

	data := map[string]interface{}{
		"Import":         priceImport,
		"Matches":        matches,
		"StatusCounts":   counts,
		"Threshold":      h.config.AutoApproveThreshold,
		"UnmatchedCount": unmatchedCount,
	}

	if err := h.renderer.Render(w, "price_import_review", data); err != nil {
		logger.Error("failed to render review page", "error", err)
	}
}

// UpdateMatchStatus approves or rejects a single match, optionally with a new name.
func (h *Handler) UpdateMatchStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	matchID := r.PathValue("id")
	if matchID == "" {
		http.Error(w, "Match ID required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if status != "approved" && status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(matchID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Check if a new name was provided
	newName := strings.TrimSpace(r.FormValue("new_name"))
	var match repository.PriceImportMatch

	if newName != "" {
		match, err = h.queries.UpdateMatchWithName(ctx, repository.UpdateMatchWithNameParams{
			ID:      id,
			Status:  status,
			NewName: sql.NullString{String: newName, Valid: true},
		})
	} else {
		match, err = h.queries.UpdateMatchStatus(ctx, repository.UpdateMatchStatusParams{
			ID:     id,
			Status: status,
		})
	}
	if err != nil {
		logger.Error("failed to update match status", "error", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	// Return updated row partial
	if r.Header.Get("HX-Request") == "true" {
		var buf bytes.Buffer
		if err := h.renderer.RenderPartial(&buf, "match_row", match); err != nil {
			logger.Error("failed to render match row", "error", err)
			http.Error(w, "Failed to render", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
		return
	}

	// Redirect back to review page
	http.Redirect(w, r, "/price-import/"+match.ImportID+"/review", http.StatusSeeOther)
}

// CreateTemplateFromMatch creates a new item template from an unmatched import row.
func (h *Handler) CreateTemplateFromMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	matchID := r.PathValue("id")
	if matchID == "" {
		http.Error(w, "Match ID required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(matchID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Get form values
	name := strings.TrimSpace(r.FormValue("name"))
	unit := strings.TrimSpace(r.FormValue("unit"))
	category := strings.TrimSpace(r.FormValue("category"))
	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = "material" // default
	}

	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create the new template
	template, err := h.queries.CreateItemTemplate(ctx, repository.CreateItemTemplateParams{
		Type:         itemType,
		Category:     category,
		Name:         name,
		DefaultUnit:  unit,
		DefaultPrice: price,
	})
	if err != nil {
		logger.Error("failed to create template", "error", err)
		http.Error(w, "Failed to create template", http.StatusInternalServerError)
		return
	}

	// Mark the match as created and link to the new template
	match, err := h.queries.MarkMatchAsCreated(ctx, repository.MarkMatchAsCreatedParams{
		ID:                id,
		MatchedTemplateID: sql.NullInt64{Int64: template.ID, Valid: true},
	})
	if err != nil {
		logger.Error("failed to update match", "error", err)
		http.Error(w, "Failed to update match", http.StatusInternalServerError)
		return
	}

	logger.Info("created template from import", "template_id", template.ID, "name", name)

	// Return updated row partial
	if r.Header.Get("HX-Request") == "true" {
		var buf bytes.Buffer
		if err := h.renderer.RenderPartial(&buf, "match_row", match); err != nil {
			logger.Error("failed to render match row", "error", err)
			http.Error(w, "Failed to render", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
		return
	}

	// Redirect back to review page
	http.Redirect(w, r, "/price-import/"+match.ImportID+"/review", http.StatusSeeOther)
}

// BulkApproveMatches approves all pending matches above threshold.
func (h *Handler) BulkApproveMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	importID := r.PathValue("id")
	if importID == "" {
		http.Error(w, "Import ID required", http.StatusBadRequest)
		return
	}

	// Get threshold from form or use config default
	threshold := h.config.AutoApproveThreshold
	if t := r.FormValue("threshold"); t != "" {
		if parsed, err := strconv.ParseFloat(t, 64); err == nil {
			threshold = parsed
		}
	}

	// Bulk approve
	if err := h.queries.BulkAutoApproveMatches(ctx, repository.BulkAutoApproveMatchesParams{
		ImportID:   importID,
		Confidence: threshold,
	}); err != nil {
		logger.Error("failed to bulk approve", "error", err)
		http.Error(w, "Failed to bulk approve", http.StatusInternalServerError)
		return
	}

	// Redirect back to review page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/price-import/"+importID+"/review")
		return
	}
	http.Redirect(w, r, "/price-import/"+importID+"/review", http.StatusSeeOther)
}

// BulkCreateTemplates creates new item templates from all unmatched items.
func (h *Handler) BulkCreateTemplates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	importID := r.PathValue("id")
	if importID == "" {
		http.Error(w, "Import ID required", http.StatusBadRequest)
		return
	}

	// Get default type from form or use material
	itemType := r.FormValue("type")
	if itemType == "" {
		itemType = "material"
	}

	// Get all unmatched items
	unmatched, err := h.queries.ListUnmatchedItems(ctx, importID)
	if err != nil {
		logger.Error("failed to list unmatched items", "error", err)
		http.Error(w, "Failed to load unmatched items", http.StatusInternalServerError)
		return
	}

	// Create templates for each unmatched item
	createdCount := 0
	for _, item := range unmatched {
		// Create the new template
		template, err := h.queries.CreateItemTemplate(ctx, repository.CreateItemTemplateParams{
			Type:         itemType,
			Category:     "", // No category for bulk create
			Name:         item.SourceName,
			DefaultUnit:  item.SourceUnit.String,
			DefaultPrice: item.SourcePrice,
		})
		if err != nil {
			logger.Error("failed to create template", "error", err, "name", item.SourceName)
			continue
		}

		// Mark the match as created and link to the new template
		_, err = h.queries.MarkMatchAsCreated(ctx, repository.MarkMatchAsCreatedParams{
			ID:                item.ID,
			MatchedTemplateID: sql.NullInt64{Int64: template.ID, Valid: true},
		})
		if err != nil {
			logger.Error("failed to update match", "error", err, "match_id", item.ID)
			continue
		}

		createdCount++
	}

	logger.Info("bulk created templates from import", "import_id", importID, "created", createdCount)

	// Redirect back to review page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/price-import/"+importID+"/review")
		return
	}
	http.Redirect(w, r, "/price-import/"+importID+"/review", http.StatusSeeOther)
}

// ApplyPriceUpdates applies approved matches to item templates.
func (h *Handler) ApplyPriceUpdates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	importID := r.PathValue("id")
	if importID == "" {
		http.Error(w, "Import ID required", http.StatusBadRequest)
		return
	}

	// Get approved matches
	matches, err := h.queries.ListApprovedMatches(ctx, importID)
	if err != nil {
		logger.Error("failed to list approved matches", "error", err)
		http.Error(w, "Failed to load matches", http.StatusInternalServerError)
		return
	}

	// Apply price updates
	updatedCount := 0
	for _, match := range matches {
		if !match.MatchedTemplateID.Valid {
			continue
		}

		// If a new name was specified, update both name and price
		if match.NewName.Valid && match.NewName.String != "" {
			if err := h.queries.UpdateItemTemplatePriceAndName(ctx, repository.UpdateItemTemplatePriceAndNameParams{
				ID:           match.MatchedTemplateID.Int64,
				DefaultPrice: match.SourcePrice,
				Name:         match.NewName.String,
			}); err != nil {
				logger.Error("failed to update template price and name", "error", err, "template_id", match.MatchedTemplateID.Int64)
				continue
			}
		} else {
			if err := h.queries.UpdateItemTemplatePrice(ctx, repository.UpdateItemTemplatePriceParams{
				ID:           match.MatchedTemplateID.Int64,
				DefaultPrice: match.SourcePrice,
			}); err != nil {
				logger.Error("failed to update template price", "error", err, "template_id", match.MatchedTemplateID.Int64)
				continue
			}
		}
		updatedCount++
	}

	// Mark import as applied
	_, err = h.queries.MarkPriceImportApplied(ctx, importID)
	if err != nil {
		logger.Error("failed to mark import applied", "error", err)
	}

	logger.Info("applied price updates", "import_id", importID, "updated", updatedCount)

	// Redirect with success message
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/price-import?success="+strconv.Itoa(updatedCount))
		return
	}
	http.Redirect(w, r, "/price-import?success="+strconv.Itoa(updatedCount), http.StatusSeeOther)
}
