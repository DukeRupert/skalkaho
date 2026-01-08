package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// TokenLength is the number of random bytes for token generation (256 bits).
	TokenLength = 32
	// TokenExpiryDays is the default expiration period for signature tokens.
	TokenExpiryDays = 30
)

// ConsentText is the exact text shown to signers and stored in the signature record.
const ConsentText = `By checking this box and clicking "Accept Quote," I agree to the terms and pricing in this quote. I confirm that the name entered above is my legal name and that this action constitutes my electronic signature.`

// GenerateSecureToken creates a cryptographically secure random token.
// Returns the URL-safe base64 encoded token.
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashContent creates a SHA-256 hash of the given content.
// Used to create document hashes for signature verification.
func HashContent(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// CreateQuoteSnapshot builds a QuoteSnapshot from the provided data.
func CreateQuoteSnapshot(
	estimateID string,
	version int,
	jobName string,
	clientName string,
	clientEmail string,
	grandTotal float64,
	notes *string,
	categories []CategorySnapshot,
) QuoteSnapshot {
	return QuoteSnapshot{
		EstimateID:  estimateID,
		Version:     version,
		JobName:     jobName,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		GrandTotal:  grandTotal,
		Notes:       notes,
		Categories:  categories,
		GeneratedAt: time.Now().UTC(),
	}
}

// SerializeSnapshot converts a QuoteSnapshot to JSON bytes.
func SerializeSnapshot(snapshot QuoteSnapshot) ([]byte, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("serializing snapshot: %w", err)
	}
	return data, nil
}

// DeserializeSnapshot converts JSON bytes back to a QuoteSnapshot.
func DeserializeSnapshot(data []byte) (QuoteSnapshot, error) {
	var snapshot QuoteSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return QuoteSnapshot{}, fmt.Errorf("deserializing snapshot: %w", err)
	}
	return snapshot, nil
}

// CalculateExpiryTime returns the expiration time for a new signature request.
func CalculateExpiryTime() time.Time {
	return time.Now().UTC().AddDate(0, 0, TokenExpiryDays)
}

// SignatureRequestInput represents input for creating a signature request.
type SignatureRequestInput struct {
	EstimateID     string
	RecipientEmail string
	RecipientName  string
	Message        string
}

// Validate checks the signature request input for errors.
func (i *SignatureRequestInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.EstimateID == "" {
		errors = append(errors, ValidationError{
			Field:   "estimate_id",
			Message: "Estimate ID is required",
		})
	}

	if i.RecipientEmail == "" {
		errors = append(errors, ValidationError{
			Field:   "recipient_email",
			Message: "Recipient email is required",
		})
	}

	if i.RecipientName == "" {
		errors = append(errors, ValidationError{
			Field:   "recipient_name",
			Message: "Recipient name is required",
		})
	}

	return errors
}

// SignatureInput represents input for submitting a signature.
type SignatureInput struct {
	LegalName string
	Agreed    bool
}

// Validate checks the signature input for errors.
func (i *SignatureInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.LegalName == "" {
		errors = append(errors, ValidationError{
			Field:   "legal_name",
			Message: "Legal name is required",
		})
	} else if len(i.LegalName) < 2 {
		errors = append(errors, ValidationError{
			Field:   "legal_name",
			Message: "Legal name must be at least 2 characters",
		})
	} else if len(i.LegalName) > 200 {
		errors = append(errors, ValidationError{
			Field:   "legal_name",
			Message: "Legal name must be less than 200 characters",
		})
	}

	if !i.Agreed {
		errors = append(errors, ValidationError{
			Field:   "agreed",
			Message: "You must agree to the terms",
		})
	}

	return errors
}
