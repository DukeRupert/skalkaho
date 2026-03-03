package router

import (
	"net/http"

	"github.com/dukerupert/skalkaho/internal/auth"
	authhandler "github.com/dukerupert/skalkaho/internal/handler/auth"
	"github.com/dukerupert/skalkaho/internal/handler/keyboard"
)

// Register sets up all routes.
func Register(mux *http.ServeMux, h *keyboard.Handler, authH *authhandler.Handler, sm *auth.SessionManager) {
	// Health check (public)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static files (public)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ==================
	// Auth Routes (public - no auth required)
	// ==================
	mux.HandleFunc("GET /login", authH.GetLogin)
	mux.HandleFunc("POST /login", authH.PostLogin)
	mux.HandleFunc("GET /register", authH.GetRegister)
	mux.HandleFunc("POST /register", authH.PostRegister)
	mux.HandleFunc("GET /logout", authH.Logout)
	mux.HandleFunc("POST /logout", authH.Logout)

	// Public Signature Pages (no auth)
	mux.HandleFunc("GET /sign/{token}", h.GetSignaturePage)
	mux.HandleFunc("POST /sign/{token}", h.SubmitSignature)
	mux.HandleFunc("GET /sign/{token}/complete", h.GetSignatureComplete)

	// ==================
	// Protected Routes (require auth)
	// ==================
	// Wrap all protected handlers with session middleware and RequireAuth

	// Jobs
	mux.Handle("GET /", protect(sm, http.HandlerFunc(h.ListJobs)))
	mux.Handle("GET /jobs/{id}", protect(sm, http.HandlerFunc(h.GetJob)))
	mux.Handle("POST /jobs", protect(sm, http.HandlerFunc(h.CreateJob)))
	mux.Handle("PUT /jobs/{id}", protect(sm, http.HandlerFunc(h.UpdateJob)))
	mux.Handle("DELETE /jobs/{id}", protect(sm, http.HandlerFunc(h.DeleteJob)))
	mux.Handle("GET /job-form", protect(sm, http.HandlerFunc(h.GetJobForm)))
	mux.Handle("GET /jobs/{id}/markup", protect(sm, http.HandlerFunc(h.GetMarkupForm)))
	mux.Handle("PUT /jobs/{id}/markup", protect(sm, http.HandlerFunc(h.UpdateMarkup)))
	mux.Handle("GET /jobs/{id}/rename", protect(sm, http.HandlerFunc(h.GetJobRenameForm)))
	mux.Handle("PUT /jobs/{id}/name", protect(sm, http.HandlerFunc(h.UpdateJobName)))
	mux.Handle("GET /jobs/{id}/order-list", protect(sm, http.HandlerFunc(h.GetOrderList)))
	mux.Handle("GET /jobs/{id}/site-materials", protect(sm, http.HandlerFunc(h.GetSiteMaterials)))
	mux.Handle("GET /jobs/{id}/client", protect(sm, http.HandlerFunc(h.GetJobClientForm)))
	mux.Handle("PUT /jobs/{id}/client", protect(sm, http.HandlerFunc(h.UpdateJobClient)))

	// Estimates
	mux.Handle("GET /jobs/{jobID}/estimates", protect(sm, http.HandlerFunc(h.ListEstimates)))
	mux.Handle("GET /jobs/{jobID}/estimates/new", protect(sm, http.HandlerFunc(h.GetNewEstimateForm)))
	mux.Handle("POST /jobs/{jobID}/estimates", protect(sm, http.HandlerFunc(h.CreateEstimate)))
	mux.Handle("GET /estimates/{id}", protect(sm, http.HandlerFunc(h.GetEstimate)))
	mux.Handle("DELETE /estimates/{id}", protect(sm, http.HandlerFunc(h.DeleteEstimate)))
	mux.Handle("GET /estimates/{id}/preview", protect(sm, http.HandlerFunc(h.GetEstimatePreview)))
	mux.Handle("POST /estimates/{id}/send", protect(sm, http.HandlerFunc(h.SendEstimate)))
	mux.Handle("POST /estimates/{id}/status", protect(sm, http.HandlerFunc(h.UpdateEstimateStatus)))
	mux.Handle("PUT /estimate-categories/{id}/description", protect(sm, http.HandlerFunc(h.UpdateEstimateCategoryDescription)))

	// E-Signatures
	mux.Handle("GET /estimates/{id}/send-for-signature", protect(sm, http.HandlerFunc(h.GetSendSignatureForm)))
	mux.Handle("POST /estimates/{id}/send-for-signature", protect(sm, http.HandlerFunc(h.SendForSignature)))
	mux.Handle("POST /estimates/{id}/cancel-signature", protect(sm, http.HandlerFunc(h.CancelSignatureRequest)))

	// Categories
	mux.Handle("GET /categories/{id}", protect(sm, http.HandlerFunc(h.GetCategory)))
	mux.Handle("POST /jobs/{jobID}/categories", protect(sm, http.HandlerFunc(h.CreateCategory)))
	mux.Handle("POST /categories/{parentID}/subcategories", protect(sm, http.HandlerFunc(h.CreateSubcategory)))
	mux.Handle("DELETE /categories/{id}", protect(sm, http.HandlerFunc(h.DeleteCategory)))
	mux.Handle("GET /category-form", protect(sm, http.HandlerFunc(h.GetCategoryForm)))
	mux.Handle("GET /categories/{id}/markup", protect(sm, http.HandlerFunc(h.GetCategoryMarkupForm)))
	mux.Handle("PUT /categories/{id}/markup", protect(sm, http.HandlerFunc(h.UpdateCategoryMarkup)))
	mux.Handle("GET /categories/{id}/rename", protect(sm, http.HandlerFunc(h.GetCategoryRenameForm)))
	mux.Handle("PUT /categories/{id}/name", protect(sm, http.HandlerFunc(h.UpdateCategoryName)))

	// Job Item Types (Custom Item Types)
	mux.Handle("GET /jobs/{jobID}/item-types", protect(sm, http.HandlerFunc(h.ListJobItemTypes)))
	mux.Handle("POST /jobs/{jobID}/item-types", protect(sm, http.HandlerFunc(h.CreateJobItemType)))
	mux.Handle("GET /jobs/{jobID}/item-types/new", protect(sm, http.HandlerFunc(h.GetJobItemTypeForm)))
	mux.Handle("PUT /item-types/{id}", protect(sm, http.HandlerFunc(h.UpdateJobItemType)))
	mux.Handle("DELETE /item-types/{id}", protect(sm, http.HandlerFunc(h.DeleteJobItemType)))

	// JSON API
	mux.Handle("GET /api/jobs/{id}", protect(sm, http.HandlerFunc(h.GetJobJSON)))
	mux.Handle("PATCH /api/jobs/{id}/items/{itemId}", protect(sm, http.HandlerFunc(h.PatchLineItemJSON)))
	mux.Handle("POST /api/jobs/{id}/items", protect(sm, http.HandlerFunc(h.CreateLineItemJSON)))
	mux.Handle("DELETE /api/jobs/{id}/items/{itemId}", protect(sm, http.HandlerFunc(h.DeleteLineItemJSON)))
	mux.Handle("PATCH /api/jobs/{id}/items/reorder", protect(sm, http.HandlerFunc(h.ReorderItemsJSON)))
	mux.Handle("GET /api/items/search", protect(sm, http.HandlerFunc(h.SearchItemsJSON)))

	// Line Items
	mux.Handle("GET /categories/{categoryID}/spreadsheet-form", protect(sm, http.HandlerFunc(h.GetSpreadsheetInlineForm)))
	mux.Handle("GET /items/{id}/spreadsheet-edit", protect(sm, http.HandlerFunc(h.GetSpreadsheetEditForm)))
	mux.Handle("GET /categories/{categoryID}/batch-form", protect(sm, http.HandlerFunc(h.GetBatchForm)))
	mux.Handle("POST /categories/{categoryID}/batch-items", protect(sm, http.HandlerFunc(h.BatchCreateLineItems)))
	mux.Handle("POST /categories/{categoryID}/items", protect(sm, http.HandlerFunc(h.CreateLineItem)))
	mux.Handle("GET /categories/{categoryID}/form", protect(sm, http.HandlerFunc(h.GetInlineForm)))
	mux.Handle("GET /items/search", protect(sm, http.HandlerFunc(h.SearchItems)))
	mux.Handle("GET /items/{id}/edit", protect(sm, http.HandlerFunc(h.GetEditForm)))
	mux.Handle("PUT /items/{id}", protect(sm, http.HandlerFunc(h.UpdateLineItem)))
	mux.Handle("DELETE /items/{id}", protect(sm, http.HandlerFunc(h.DeleteLineItem)))

	// Item Templates
	mux.Handle("GET /items", protect(sm, http.HandlerFunc(h.ListItemTemplates)))
	mux.Handle("POST /items", protect(sm, http.HandlerFunc(h.CreateItemTemplate)))
	mux.Handle("GET /items/new", protect(sm, http.HandlerFunc(h.GetItemTemplateForm)))
	mux.Handle("GET /item-templates/{id}/edit", protect(sm, http.HandlerFunc(h.GetItemTemplateEditForm)))
	mux.Handle("PUT /item-templates/{id}", protect(sm, http.HandlerFunc(h.UpdateItemTemplate)))
	mux.Handle("DELETE /item-templates/{id}", protect(sm, http.HandlerFunc(h.DeleteItemTemplate)))
	mux.Handle("POST /item-templates/bulk-family", protect(sm, http.HandlerFunc(h.BulkUpdateFamily)))

	// Clients
	mux.Handle("GET /clients", protect(sm, http.HandlerFunc(h.ListClients)))
	mux.Handle("GET /clients/{id}", protect(sm, http.HandlerFunc(h.GetClient)))
	mux.Handle("POST /clients", protect(sm, http.HandlerFunc(h.CreateClient)))
	mux.Handle("PUT /clients/{id}", protect(sm, http.HandlerFunc(h.UpdateClient)))
	mux.Handle("DELETE /clients/{id}", protect(sm, http.HandlerFunc(h.DeleteClient)))
	mux.Handle("GET /client-form", protect(sm, http.HandlerFunc(h.GetClientForm)))
	mux.Handle("GET /clients/{id}/edit", protect(sm, http.HandlerFunc(h.GetClientEditForm)))

	// Settings
	mux.Handle("GET /settings", protect(sm, http.HandlerFunc(h.GetSettings)))
	mux.Handle("PUT /settings", protect(sm, http.HandlerFunc(h.UpdateSettings)))

	// Price Import
	mux.Handle("GET /price-import", protect(sm, http.HandlerFunc(h.GetPriceImportPage)))
	mux.Handle("POST /price-import/auth", protect(sm, http.HandlerFunc(h.ValidatePriceImportToken)))
	mux.Handle("POST /price-import/upload", protect(sm, http.HandlerFunc(h.UploadPriceFile)))
	mux.Handle("GET /price-import/{id}/review", protect(sm, http.HandlerFunc(h.GetImportReview)))
	mux.Handle("PUT /price-import/matches/{id}", protect(sm, http.HandlerFunc(h.UpdateMatchStatus)))
	mux.Handle("POST /price-import/matches/{id}/create-template", protect(sm, http.HandlerFunc(h.CreateTemplateFromMatch)))
	mux.Handle("POST /price-import/{id}/bulk-approve", protect(sm, http.HandlerFunc(h.BulkApproveMatches)))
	mux.Handle("POST /price-import/{id}/bulk-create", protect(sm, http.HandlerFunc(h.BulkCreateTemplates)))
	mux.Handle("POST /price-import/{id}/apply", protect(sm, http.HandlerFunc(h.ApplyPriceUpdates)))
}

// protect wraps a handler with session loading and authentication requirement.
func protect(sm *auth.SessionManager, handler http.Handler) http.Handler {
	return sm.SessionMiddleware(auth.RequireAuth(handler))
}
