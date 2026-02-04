package testutil

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Fixtures provides helper methods for creating test data.
type Fixtures struct {
	t  *testing.T
	db *sql.DB
}

// NewFixtures creates a new Fixtures instance for creating test data.
func NewFixtures(t *testing.T, db *sql.DB) *Fixtures {
	t.Helper()
	return &Fixtures{t: t, db: db}
}

// JobParams holds parameters for creating a test job.
type JobParams struct {
	ID               string
	OrgID            string // Will be used once multi-tenancy is implemented
	Name             string
	CustomerName     *string
	SurchargePercent float64
	SurchargeMode    string
	Status           string
	ClientID         *string
}

// CreateJob creates a job in the test database.
// Returns the job ID.
func (f *Fixtures) CreateJob(params JobParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Name == "" {
		params.Name = "Test Job"
	}
	if params.SurchargeMode == "" {
		params.SurchargeMode = "stacking"
	}
	if params.Status == "" {
		params.Status = "draft"
	}

	// Note: This query doesn't include org_id yet - it will fail when multi-tenancy is added
	// and will need to be updated to include org_id
	query := `
		INSERT INTO jobs (id, name, customer_name, surcharge_percent, surcharge_mode, status, created_at, client_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var customerName sql.NullString
	if params.CustomerName != nil {
		customerName = sql.NullString{String: *params.CustomerName, Valid: true}
	}

	var clientID sql.NullString
	if params.ClientID != nil {
		clientID = sql.NullString{String: *params.ClientID, Valid: true}
	}

	_, err := f.db.Exec(query,
		params.ID,
		params.Name,
		customerName,
		params.SurchargePercent,
		params.SurchargeMode,
		params.Status,
		time.Now().UTC().Format(time.RFC3339),
		clientID,
	)
	if err != nil {
		f.t.Fatalf("failed to create job: %v", err)
	}

	return params.ID
}

// CategoryParams holds parameters for creating a test category.
type CategoryParams struct {
	ID               string
	JobID            string
	ParentID         *string
	Name             string
	SurchargePercent *float64
	SortOrder        int64
}

// CreateCategory creates a category in the test database.
// Returns the category ID.
func (f *Fixtures) CreateCategory(params CategoryParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Name == "" {
		params.Name = "Test Category"
	}

	query := `
		INSERT INTO categories (id, job_id, parent_id, name, surcharge_percent, sort_order)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	var parentID sql.NullString
	if params.ParentID != nil {
		parentID = sql.NullString{String: *params.ParentID, Valid: true}
	}

	var surchargePercent sql.NullFloat64
	if params.SurchargePercent != nil {
		surchargePercent = sql.NullFloat64{Float64: *params.SurchargePercent, Valid: true}
	}

	_, err := f.db.Exec(query,
		params.ID,
		params.JobID,
		parentID,
		params.Name,
		surchargePercent,
		params.SortOrder,
	)
	if err != nil {
		f.t.Fatalf("failed to create category: %v", err)
	}

	return params.ID
}

// LineItemParams holds parameters for creating a test line item.
type LineItemParams struct {
	ID               string
	CategoryID       string
	Type             string
	Name             string
	Description      *string
	Quantity         float64
	Unit             string
	UnitPrice        float64
	SurchargePercent *float64
	SortOrder        int64
	Tag              *string
}

// CreateLineItem creates a line item in the test database.
// Returns the line item ID.
func (f *Fixtures) CreateLineItem(params LineItemParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Name == "" {
		params.Name = "Test Item"
	}
	if params.Type == "" {
		params.Type = "material"
	}
	if params.Unit == "" {
		params.Unit = "ea"
	}
	if params.Quantity == 0 {
		params.Quantity = 1
	}
	if params.UnitPrice == 0 {
		params.UnitPrice = 100
	}

	query := `
		INSERT INTO line_items (id, category_id, type, name, description, quantity, unit, unit_price, surcharge_percent, sort_order, tag)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var description sql.NullString
	if params.Description != nil {
		description = sql.NullString{String: *params.Description, Valid: true}
	}

	var surchargePercent sql.NullFloat64
	if params.SurchargePercent != nil {
		surchargePercent = sql.NullFloat64{Float64: *params.SurchargePercent, Valid: true}
	}

	var tag sql.NullString
	if params.Tag != nil {
		tag = sql.NullString{String: *params.Tag, Valid: true}
	}

	_, err := f.db.Exec(query,
		params.ID,
		params.CategoryID,
		params.Type,
		params.Name,
		description,
		params.Quantity,
		params.Unit,
		params.UnitPrice,
		surchargePercent,
		params.SortOrder,
		tag,
	)
	if err != nil {
		f.t.Fatalf("failed to create line item: %v", err)
	}

	return params.ID
}

// SettingsParams holds parameters for creating test settings.
type SettingsParams struct {
	ID                      string
	OrgID                   string // Will be used once multi-tenancy is implemented
	DefaultSurchargeMode    string
	DefaultSurchargePercent float64
}

// CreateSettings creates settings in the test database.
// Returns the settings ID.
func (f *Fixtures) CreateSettings(params SettingsParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.DefaultSurchargeMode == "" {
		params.DefaultSurchargeMode = "stacking"
	}

	query := `
		INSERT INTO settings (id, default_surcharge_mode, default_surcharge_percent)
		VALUES (?, ?, ?)
	`

	_, err := f.db.Exec(query,
		params.ID,
		params.DefaultSurchargeMode,
		params.DefaultSurchargePercent,
	)
	if err != nil {
		f.t.Fatalf("failed to create settings: %v", err)
	}

	return params.ID
}

// ClientParams holds parameters for creating a test client.
type ClientParams struct {
	ID    string
	OrgID string // Will be used once multi-tenancy is implemented
	Name  string
	Email *string
	Phone *string
}

// CreateClient creates a client in the test database.
// Returns the client ID.
func (f *Fixtures) CreateClient(params ClientParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Name == "" {
		params.Name = "Test Client"
	}

	query := `
		INSERT INTO clients (id, name, email, phone, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	var email sql.NullString
	if params.Email != nil {
		email = sql.NullString{String: *params.Email, Valid: true}
	}

	var phone sql.NullString
	if params.Phone != nil {
		phone = sql.NullString{String: *params.Phone, Valid: true}
	}

	_, err := f.db.Exec(query,
		params.ID,
		params.Name,
		email,
		phone,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		f.t.Fatalf("failed to create client: %v", err)
	}

	return params.ID
}

// EstimateParams holds parameters for creating a test estimate.
type EstimateParams struct {
	ID         string
	JobID      string
	Version    int64
	Status     string
	GrandTotal float64
	Notes      *string
}

// CreateEstimate creates an estimate in the test database.
// Returns the estimate ID.
func (f *Fixtures) CreateEstimate(params EstimateParams) string {
	f.t.Helper()

	if params.ID == "" {
		params.ID = uuid.New().String()
	}
	if params.Status == "" {
		params.Status = "draft"
	}
	if params.Version == 0 {
		params.Version = 1
	}

	query := `
		INSERT INTO estimates (id, job_id, version, status, grand_total, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var notes sql.NullString
	if params.Notes != nil {
		notes = sql.NullString{String: *params.Notes, Valid: true}
	}

	_, err := f.db.Exec(query,
		params.ID,
		params.JobID,
		params.Version,
		params.Status,
		params.GrandTotal,
		notes,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		f.t.Fatalf("failed to create estimate: %v", err)
	}

	return params.ID
}

// Helper functions for common patterns

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// Float64Ptr returns a pointer to the given float64.
func Float64Ptr(f float64) *float64 {
	return &f
}

// Int64Ptr returns a pointer to the given int64.
func Int64Ptr(i int64) *int64 {
	return &i
}
