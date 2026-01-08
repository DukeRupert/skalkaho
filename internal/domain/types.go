package domain

import "time"

// SurchargeMode defines how surcharges are calculated.
type SurchargeMode string

const (
	SurchargeModeStacking SurchargeMode = "stacking"
	SurchargeModeOverride SurchargeMode = "override"
)

// LineItemType distinguishes materials, labor, and equipment.
type LineItemType string

const (
	LineItemTypeMaterial  LineItemType = "material"
	LineItemTypeLabor     LineItemType = "labor"
	LineItemTypeEquipment LineItemType = "equipment"
)

// StandardTypes returns the list of built-in line item types.
var StandardTypes = []LineItemType{
	LineItemTypeMaterial,
	LineItemTypeLabor,
	LineItemTypeEquipment,
}

// IsStandardType returns true if the type is one of the built-in types.
func IsStandardType(t LineItemType) bool {
	return t == LineItemTypeMaterial || t == LineItemTypeLabor || t == LineItemTypeEquipment
}

// JobItemType represents a custom line item type for a specific job.
type JobItemType struct {
	ID               string    `json:"id"`
	JobID            string    `json:"job_id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Color            string    `json:"color"`
	SortOrder        int       `json:"sort_order"`
	SurchargePercent *float64  `json:"surcharge_percent,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// Settings holds application-wide defaults.
type Settings struct {
	ID                      string        `json:"id"`
	DefaultSurchargeMode    SurchargeMode `json:"default_surcharge_mode"`
	DefaultSurchargePercent float64       `json:"default_surcharge_percent"`
}

// Job is the top-level container for a quote.
type Job struct {
	ID                        string        `json:"id"`
	Name                      string        `json:"name"`
	CustomerName              *string       `json:"customer_name,omitempty"`
	SurchargePercent          float64       `json:"surcharge_percent"`
	MaterialSurchargePercent  *float64      `json:"material_surcharge_percent,omitempty"`
	LaborSurchargePercent     *float64      `json:"labor_surcharge_percent,omitempty"`
	EquipmentSurchargePercent *float64      `json:"equipment_surcharge_percent,omitempty"`
	SurchargeMode             SurchargeMode `json:"surcharge_mode"`
	CreatedAt                 time.Time     `json:"created_at"`
}

// Category represents an organizational grouping within a job.
type Category struct {
	ID               string   `json:"id"`
	JobID            string   `json:"job_id"`
	ParentID         *string  `json:"parent_id,omitempty"`
	Name             string   `json:"name"`
	SurchargePercent *float64 `json:"surcharge_percent,omitempty"`
	SortOrder        int      `json:"sort_order"`
}

// LineItem represents an individual material or labor entry.
type LineItem struct {
	ID               string       `json:"id"`
	CategoryID       string       `json:"category_id"`
	Type             LineItemType `json:"type"`
	Name             string       `json:"name"`
	Description      *string      `json:"description,omitempty"`
	Quantity         float64      `json:"quantity"`
	Unit             string       `json:"unit"`
	UnitPrice        float64      `json:"unit_price"`
	SurchargePercent *float64     `json:"surcharge_percent,omitempty"`
	SortOrder        int          `json:"sort_order"`
	Tag              *string      `json:"tag,omitempty"`
}

// BasePrice calculates quantity * unit_price.
func (li *LineItem) BasePrice() float64 {
	return li.Quantity * li.UnitPrice
}

// CommonUnits returns suggested units for the UI.
var CommonUnits = struct {
	Material []string
	Labor    []string
}{
	Material: []string{"ea", "sqft", "lnft", "bundle", "box", "bag", "gal", "sheet"},
	Labor:    []string{"hr", "day", "job", "sqft"},
}

// EstimateStatus defines the lifecycle states for an estimate.
type EstimateStatus string

const (
	EstimateStatusDraft    EstimateStatus = "draft"
	EstimateStatusSent     EstimateStatus = "sent"
	EstimateStatusAccepted EstimateStatus = "accepted"
	EstimateStatusRejected EstimateStatus = "rejected"
)

// Estimate represents a client-facing snapshot of a quote.
type Estimate struct {
	ID          string         `json:"id"`
	JobID       string         `json:"job_id"`
	Version     int            `json:"version"`
	Status      EstimateStatus `json:"status"`
	GrandTotal  float64        `json:"grand_total"`
	Notes       *string        `json:"notes,omitempty"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	RespondedAt *time.Time     `json:"responded_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// EstimateCategory represents a category snapshot in an estimate.
type EstimateCategory struct {
	ID               string  `json:"id"`
	EstimateID       string  `json:"estimate_id"`
	CategoryID       string  `json:"category_id"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
	Tier             int     `json:"tier"`
	Name             string  `json:"name"`
	Description      *string `json:"description,omitempty"`
	Total            float64 `json:"total"`
	SortOrder        int     `json:"sort_order"`
}

// CompanyProfile holds contractor/business information.
type CompanyProfile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Address   *string   `json:"address,omitempty"`
	City      *string   `json:"city,omitempty"`
	State     *string   `json:"state,omitempty"`
	Zip       *string   `json:"zip,omitempty"`
	LogoPath  *string   `json:"logo_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignatureRequestStatus defines lifecycle states for signature requests.
type SignatureRequestStatus string

const (
	SignatureRequestStatusPending   SignatureRequestStatus = "pending"
	SignatureRequestStatusSigned    SignatureRequestStatus = "signed"
	SignatureRequestStatusExpired   SignatureRequestStatus = "expired"
	SignatureRequestStatusCancelled SignatureRequestStatus = "cancelled"
)

// SignatureRequest tracks outbound signature requests.
type SignatureRequest struct {
	ID              string                 `json:"id"`
	EstimateID      string                 `json:"estimate_id"`
	RecipientEmail  string                 `json:"recipient_email"`
	RecipientName   string                 `json:"recipient_name"`
	Token           string                 `json:"-"` // Never expose token in JSON
	DocumentHash    string                 `json:"document_hash"`
	QuoteSnapshot   string                 `json:"-"` // Large JSON, exclude from default serialization
	Message         *string                `json:"message,omitempty"`
	Status          SignatureRequestStatus `json:"status"`
	ExpiresAt       time.Time              `json:"expires_at"`
	SenderIP        *string                `json:"sender_ip,omitempty"`
	SenderUserAgent *string                `json:"sender_user_agent,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// Signature captures the actual signing event (immutable after creation).
type Signature struct {
	ID                 string    `json:"id"`
	RequestID          string    `json:"request_id"`
	LegalName          string    `json:"legal_name"`
	ConsentText        string    `json:"consent_text"`
	DocumentHash       string    `json:"document_hash"`
	SignedAt           time.Time `json:"signed_at"`
	SignerIP           string    `json:"signer_ip"`
	SignerUserAgent    string    `json:"signer_user_agent"`
	SignerEmail        string    `json:"signer_email"`
	CertificatePDFPath *string   `json:"certificate_pdf_path,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

// QuoteSnapshot captures the quote state at send time for immutable reference.
type QuoteSnapshot struct {
	JobName     string             `json:"job_name"`
	ClientName  string             `json:"client_name"`
	ClientEmail string             `json:"client_email"`
	Categories  []CategorySnapshot `json:"categories"`
	GrandTotal  float64            `json:"grand_total"`
	Notes       *string            `json:"notes,omitempty"`
	GeneratedAt time.Time          `json:"generated_at"`
	EstimateID  string             `json:"estimate_id"`
	Version     int                `json:"version"`
}

// CategorySnapshot captures category state for the quote snapshot.
type CategorySnapshot struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Total       float64 `json:"total"`
	Tier        int     `json:"tier"`
	SortOrder   int     `json:"sort_order"`
}
