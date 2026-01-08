# Skalkaho E-Signature Implementation Plan

## Overview

Simple "click to accept" signature capture with legal name, audit trail, and PDF certificate. No drawn signatures, no multi-signer workflows, no third-party dependencies.

---

## User Flows

### Sender Flow (Contractor)

1. From quote view, click "Send for Signature"
2. Enter recipient email and name
3. Optionally add a message
4. Click "Send"
5. Quote status changes to "Awaiting Signature"
6. Receive email notification when signed

### Signer Flow (Customer)

1. Receive email with link to review quote
2. Click link → lands on public review page
3. See quote details (read-only)
4. Enter full legal name in text field
5. Check consent checkbox
6. Click "Accept Quote"
7. See confirmation page
8. Receive copy of signed quote via email

---

## Data Model

### signature_requests

Tracks outbound signature requests.

```sql
CREATE TABLE signature_requests (
    id TEXT PRIMARY KEY,  -- UUID
    quote_id TEXT NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    
    -- Recipient info
    recipient_email TEXT NOT NULL,
    recipient_name TEXT NOT NULL,
    
    -- Security
    token TEXT NOT NULL UNIQUE,  -- Secure random token for URL
    
    -- Document snapshot (immutable after sending)
    document_hash TEXT NOT NULL,  -- SHA-256 of quote content at send time
    quote_snapshot TEXT NOT NULL, -- JSON snapshot of quote data
    
    -- Optional message from sender
    message TEXT,
    
    -- Lifecycle
    status TEXT NOT NULL DEFAULT 'pending',  -- pending, signed, expired, cancelled
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,  -- e.g., 30 days from sent_at
    
    -- Sender context (for audit)
    sender_ip TEXT,
    sender_user_agent TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_signature_requests_token ON signature_requests(token);
CREATE INDEX idx_signature_requests_quote_id ON signature_requests(quote_id);
CREATE INDEX idx_signature_requests_status ON signature_requests(status);
```

### signatures

Captures the actual signature event. Immutable after creation.

```sql
CREATE TABLE signatures (
    id TEXT PRIMARY KEY,  -- UUID
    request_id TEXT NOT NULL REFERENCES signature_requests(id) ON DELETE RESTRICT,
    
    -- What they signed
    legal_name TEXT NOT NULL,  -- As typed by signer
    consent_text TEXT NOT NULL,  -- Exact text of the consent checkbox
    document_hash TEXT NOT NULL,  -- Must match request's document_hash
    
    -- Audit trail
    signed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signer_ip TEXT NOT NULL,
    signer_user_agent TEXT NOT NULL,
    signer_email TEXT NOT NULL,  -- Captured at sign time (may differ from recipient_email)
    
    -- Generated artifacts
    certificate_pdf_path TEXT,  -- Path to stored certificate PDF
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_signatures_request_id ON signatures(request_id);
```

### Quote status integration

Add to existing quotes table or handle via join:

```sql
-- Option A: Add column to quotes
ALTER TABLE quotes ADD COLUMN signature_status TEXT DEFAULT 'draft';
-- Values: draft, pending_signature, signed

-- Option B: Derive from signature_requests (no schema change)
-- Query: SELECT status FROM signature_requests WHERE quote_id = ? ORDER BY created_at DESC LIMIT 1
```

---

## Security Considerations

### Token Generation

```go
// Generate 32-byte cryptographically secure random token
// Encode as URL-safe base64 (43 characters)
func generateSignatureToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}
```

### Document Hash

```go
// Create SHA-256 hash of quote content at send time
// This proves the document wasn't altered after sending
func hashQuoteContent(snapshot []byte) string {
    hash := sha256.Sum256(snapshot)
    return hex.EncodeToString(hash[:])
}
```

### Quote Snapshot

Store a complete JSON representation of the quote at send time. This is what the signer sees and agrees to—even if the original quote is later modified or deleted.

```go
type QuoteSnapshot struct {
    QuoteName      string            `json:"quote_name"`
    CustomerName   string            `json:"customer_name"`
    Categories     []CategorySnapshot `json:"categories"`
    Subtotal       float64           `json:"subtotal"`
    MarkupAmount   float64           `json:"markup_amount"`
    Total          float64           `json:"total"`
    MarkupPercent  float64           `json:"markup_percent"`
    MarkupMode     string            `json:"markup_mode"`
    GeneratedAt    time.Time         `json:"generated_at"`
}
```

### Rate Limiting

- Limit signature request sends (e.g., 10 per hour per user)
- Limit signature page views (e.g., 30 per hour per token)
- Limit signature attempts (e.g., 5 per hour per token)

### Expiration

- Default expiration: 30 days from send
- Expired requests show "This quote has expired" page
- Sender can resend (creates new request with fresh token)

---

## Consent Language

The consent checkbox text must be:

1. Clear about what they're agreeing to
2. Stored exactly as shown in the signatures table
3. Legally defensible

**Recommended text:**

> By checking this box and clicking "Accept Quote," I agree to the terms and pricing in this quote. I confirm that the name entered above is my legal name and that this action constitutes my electronic signature.

---

## Email Templates

### Signature Request Email (to customer)

**Subject:** Quote from [Company Name] ready for your review

```
Hi [Recipient Name],

[Sender Name] has sent you a quote for your review.

[If message exists:]
Message from [Sender Name]:
"[Message]"

Review and accept your quote:
[Button: View Quote]

This link expires on [Expiration Date].

---
Sent via Skalkaho
```

### Signature Confirmation Email (to customer)

**Subject:** Quote accepted - [Quote Name]

```
Hi [Signer Name],

You accepted the quote "[Quote Name]" on [Date] at [Time].

A copy of the signed quote is attached for your records.

---
Signed electronically via Skalkaho
```

### Signature Notification Email (to contractor)

**Subject:** Quote signed by [Signer Name]

```
Hi [Sender Name],

[Signer Name] has accepted your quote "[Quote Name]".

Signed: [Date] at [Time]

[Button: View Quote]

---
Sent via Skalkaho
```

---

## PDF Certificate

Append a signature certificate page to the quote PDF when signed. This page documents the electronic signature for legal records.

### Certificate Content

```
───────────────────────────────────────────────────────────
ELECTRONIC SIGNATURE CERTIFICATE
───────────────────────────────────────────────────────────

Document:       [Quote Name]
Document ID:    [Quote ID]
Document Hash:  [SHA-256 hash]

───────────────────────────────────────────────────────────
SIGNATURE DETAILS
───────────────────────────────────────────────────────────

Signer Name:    [Legal Name as entered]
Signer Email:   [Email address]
Signed At:      [Timestamp with timezone]
IP Address:     [IP address]

───────────────────────────────────────────────────────────
CONSENT
───────────────────────────────────────────────────────────

The signer agreed to the following statement:

"By checking this box and clicking 'Accept Quote,' I agree 
to the terms and pricing in this quote. I confirm that the 
name entered above is my legal name and that this action 
constitutes my electronic signature."

───────────────────────────────────────────────────────────
VERIFICATION
───────────────────────────────────────────────────────────

This document was electronically signed using Skalkaho.
The document hash above can be used to verify that this
document has not been altered since signing.

Signature ID:   [Signature UUID]
Request ID:     [Request UUID]

───────────────────────────────────────────────────────────
```

---

## UI Components Needed

### Send for Signature Modal

```
┌─────────────────────────────────────────────────────────┐
│ Send for Signature                                   ✕  │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Recipient Email                                        │
│  ┌───────────────────────────────────────────────────┐  │
│  │ john.customer@email.com                           │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  Recipient Name                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │ John Customer                                     │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  Message (optional)                                     │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Thanks for the opportunity to bid on your        │  │
│  │ project. Let me know if you have questions.      │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  Link expires in 30 days.                               │
│                                                         │
├─────────────────────────────────────────────────────────┤
│                              [Cancel]  [Send Quote]     │
└─────────────────────────────────────────────────────────┘
```

### Public Signature Page

```
┌─────────────────────────────────────────────────────────┐
│  ◢◣ Skalkaho                                            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Quote for: John Customer                               │
│  From: Williams Construction                            │
│  Date: January 8, 2026                                  │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │                                                   │  │
│  │  [Quote content rendered read-only]               │  │
│  │                                                   │  │
│  │  Categories, line items, totals...                │  │
│  │                                                   │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  ─────────────────────────────────────────────────────  │
│                                                         │
│  TOTAL: $14,317.50                                      │
│                                                         │
│  ─────────────────────────────────────────────────────  │
│                                                         │
│  To accept this quote, enter your full legal name:      │
│                                                         │
│  Full Legal Name                                        │
│  ┌───────────────────────────────────────────────────┐  │
│  │                                                   │  │
│  └───────────────────────────────────────────────────┘  │
│                                                         │
│  ☐ By checking this box and clicking "Accept Quote,"   │
│    I agree to the terms and pricing in this quote.     │
│    I confirm that the name entered above is my legal   │
│    name and that this action constitutes my            │
│    electronic signature.                               │
│                                                         │
│              [Accept Quote]                             │
│                                                         │
│  ─────────────────────────────────────────────────────  │
│  Questions? Contact [sender email]                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Confirmation Page

```
┌─────────────────────────────────────────────────────────┐
│  ◢◣ Skalkaho                                            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│                         ✓                               │
│                                                         │
│              Quote Accepted                             │
│                                                         │
│  You accepted this quote on January 8, 2026             │
│  at 3:42 PM MST.                                        │
│                                                         │
│  A confirmation email with a copy of the signed         │
│  quote has been sent to john.customer@email.com.        │
│                                                         │
│              [Download Signed Quote]                    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Quote Status Badge (in contractor's view)

```
Draft              →  (no badge, or subtle "Draft")
Awaiting Signature →  [🕐 Awaiting Signature] (yellow/amber)
Signed             →  [✓ Signed Jan 8]        (green)
Expired            →  [✕ Expired]             (gray)
```

---

## Implementation Steps

### Phase 1: Database & Models (Day 1)

1. Create migration for `signature_requests` table
2. Create migration for `signatures` table
3. Create Go structs for both entities
4. Create repository functions:
   - `CreateSignatureRequest`
   - `GetSignatureRequestByToken`
   - `GetSignatureRequestsByQuoteID`
   - `UpdateSignatureRequestStatus`
   - `CreateSignature`
   - `GetSignatureByRequestID`

**Test:** Write unit tests for repository functions.

---

### Phase 2: Quote Snapshot & Hashing (Day 2)

1. Create `QuoteSnapshot` struct with all quote data
2. Implement `CreateQuoteSnapshot(quoteID) ([]byte, error)`
3. Implement `HashQuoteContent(snapshot []byte) string`
4. Implement `RenderQuoteFromSnapshot(snapshot []byte) (HTML, error)`

**Test:** Snapshot a quote, modify original, verify snapshot unchanged. Verify hash changes if snapshot changes.

---

### Phase 3: Send for Signature (Day 3)

1. Add "Send for Signature" button to quote view
2. Create send modal component
3. Implement `POST /quotes/{id}/send-for-signature`:
   - Validate recipient email and name
   - Generate secure token
   - Create quote snapshot
   - Hash snapshot
   - Create signature_request record
   - Send email to recipient
   - Return success response
4. Update quote status display

**Test:** Send signature request, verify email received with valid link.

---

### Phase 4: Public Signature Page (Day 4-5)

1. Create public route `GET /sign/{token}`
2. Implement token validation:
   - Token exists
   - Not expired
   - Not already signed
   - Not cancelled
3. Render quote from snapshot (read-only)
4. Create signature form:
   - Legal name input
   - Consent checkbox
   - Accept button
5. Add client-side validation:
   - Name required, minimum length
   - Checkbox required
   - Button disabled until valid

**Test:** Access signature page, verify quote displays correctly, form validates.

---

### Phase 5: Signature Capture (Day 6)

1. Implement `POST /sign/{token}`:
   - Validate token (same as GET)
   - Validate form data
   - Verify document hash matches
   - Capture audit data (IP, user agent, timestamp)
   - Create signature record
   - Update request status to 'signed'
   - Update quote signature_status
2. Redirect to confirmation page
3. Trigger async tasks:
   - Generate signed PDF with certificate
   - Send confirmation to signer
   - Send notification to sender

**Test:** Complete signature flow, verify all records created correctly.

---

### Phase 6: PDF Certificate (Day 7)

1. Create certificate page template
2. Implement `GenerateSignedQuotePDF`:
   - Render quote from snapshot
   - Append certificate page
   - Store PDF
3. Implement `GET /sign/{token}/download`:
   - Validate token and signature exists
   - Return signed PDF

**Test:** Download signed PDF, verify certificate page present with correct data.

---

### Phase 7: Email Notifications (Day 8)

1. Create email templates:
   - Signature request (to customer)
   - Signature confirmation (to customer, with PDF attachment)
   - Signature notification (to contractor)
2. Implement email sending for each trigger
3. Add sender email/company name to settings (if not already present)

**Test:** Verify all emails sent at correct times with correct content.

---

### Phase 8: Status & History (Day 9)

1. Add signature status badge to quote list
2. Add signature history section to quote view:
   - Sent to [email] on [date]
   - Signed by [name] on [date]
   - Or: Expired on [date] / Cancelled on [date]
3. Add "Resend" action for pending/expired requests
4. Add "Cancel" action for pending requests

**Test:** Verify status displays correctly, resend creates new request.

---

### Phase 9: Edge Cases & Polish (Day 10)

1. Handle expired signature page (show friendly message)
2. Handle already-signed page (show confirmation)
3. Handle cancelled request (show message)
4. Rate limiting
5. Error handling and user feedback
6. Keyboard navigation for signature form
7. Mobile-responsive signature page

**Test:** Test all edge cases, verify graceful handling.

---

## Future Considerations (Not in v1)

- Signature request reminders (email after X days if not signed)
- Decline with reason
- Multiple signers
- Signature request templates
- Bulk send
- Drawn signature option
- SMS notifications
- Webhook notifications for integrations

---

## Open Questions

1. **Sender identity:** Do you need a company/contractor profile in settings, or derive from the first user? (For "From: Williams Construction" on signature page)

2. **Email service:** What are you using for transactional email? (SendGrid, Postmark, SES, etc.)

3. **PDF storage:** Store locally or cloud storage (S3, etc.)?

4. **Expiration length:** 30 days default, or configurable?

5. **Quote editing after send:** Lock the quote when awaiting signature, or allow edits (which would invalidate the pending request)?