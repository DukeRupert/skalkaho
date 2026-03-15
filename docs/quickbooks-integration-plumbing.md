# QuickBooks Integration — Implementation Plan

## Overview

Connect Skalkaho to QuickBooks Online to push invoices against signed quotes
and reflect incoming payments back in Skalkaho. The owner completes a one-time
OAuth flow to authorize the connection. Customers are manually linked to their
QBO counterparts to avoid duplicates. Line item sync strategy and draw schedule
design are deferred pending client conversations.

---

## Scope (This Plan)

| In scope | Deferred |
|----------|----------|
| OAuth 2.0 connection + token management | Line item / category sync to QBO Items |
| Token refresh background job | Draw schedule (percentage / milestone / ad hoc) |
| Customer ↔ QBO manual linking UI | Full invoice line item detail |
| Invoice push to QBO (summary level) | |
| Webhook handler for payment sync | |
| Payment reflection in Skalkaho | |
| QBO connection status page | |

---

## Key QuickBooks API Facts

- **Auth:** OAuth 2.0 only — no static API keys. One-time owner authorization,
  then fully automatic via token refresh.
- **Access token TTL:** 1 hour — must be refreshed automatically.
- **Refresh token TTL:** 100 days — if not refreshed within this window, the
  owner must re-authorize. App should warn well before expiry.
- **Invoice line items:** Require a QBO `ItemRef` — deferred until line item
  sync is designed.
- **Partial payments:** QBO Payments allows one online payment per invoice
  (partial or full). Additional payments are recorded manually or via the
  Payment entity in the API. Contractor draw schedules should use separate
  invoices per draw, not partial payments on one invoice.
- **Webhooks:** Fire on Payment create/update/delete. Must use CloudEvents
  format (required by QBO after May 15, 2026).
- **SyncToken:** Required on all QBO updates — fetch the current record first,
  include its SyncToken, or the update will be rejected with a stale object error.

---

## Data Model

### Schema additions

```sql
-- OAuth tokens and connection state for the QBO integration
-- Single row table — this app has one QBO connection
CREATE TABLE qb_connection (
    id                INTEGER PRIMARY KEY DEFAULT 1,  -- enforces single row
    realm_id          TEXT NOT NULL,                  -- QBO company ID
    access_token      TEXT NOT NULL,
    refresh_token     TEXT NOT NULL,
    access_token_expiry  TIMESTAMPTZ NOT NULL,
    refresh_token_expiry TIMESTAMPTZ NOT NULL,
    connected_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_refreshed_at TIMESTAMPTZ,
    CHECK (id = 1)                                    -- single row constraint
);

-- QBO customer link — manually set by owner
ALTER TABLE customers
    ADD COLUMN qb_customer_id TEXT;                   -- QBO CustomerRef value

ALTER TABLE customers
    ADD COLUMN qb_customer_name TEXT;                 -- display name from QBO, for confirmation UI

-- Invoice tracking
CREATE TABLE invoices (
    id              TEXT PRIMARY KEY,                 -- nanoid
    quote_id        TEXT NOT NULL REFERENCES quotes(id),
    qb_invoice_id   TEXT,                             -- NULL until pushed to QBO
    qb_invoice_num  TEXT,                             -- QBO DocNumber, for display
    qb_sync_token   TEXT,                             -- required for QBO updates
    amount          REAL NOT NULL,
    status          TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'pushed', 'partially_paid', 'paid', 'voided')),
    due_date        DATE,
    memo            TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payment records synced back from QBO
CREATE TABLE invoice_payments (
    id              TEXT PRIMARY KEY,                 -- nanoid
    invoice_id      TEXT NOT NULL REFERENCES invoices(id),
    qb_payment_id   TEXT NOT NULL UNIQUE,             -- QBO Payment Id
    amount          REAL NOT NULL,
    payment_date    DATE NOT NULL,
    synced_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for webhook handler lookups
CREATE INDEX idx_invoices_qb_invoice_id ON invoices(qb_invoice_id)
    WHERE qb_invoice_id IS NOT NULL;

CREATE INDEX idx_invoice_payments_qb_payment_id ON invoice_payments(qb_payment_id);
```

---

## Step 1 — Database Migration

**File:** `migrations/XXXX_add_quickbooks_integration.sql`

Run the schema additions above. Note the `CHECK (id = 1)` constraint on
`qb_connection` — this is a single-row table enforced at the DB level. Any
attempt to insert a second row will fail.

**Test:** Insert a `qb_connection` row, confirm a second insert is rejected.
Add `qb_customer_id` to an existing customer, confirm nullable. Insert an
invoice and a payment, confirm cascade behavior if quote is deleted (decide:
restrict or cascade — invoices against a deleted quote is a meaningful state
to preserve, so `RESTRICT` is probably correct).

---

## Step 2 — sqlc Queries

**File:** `db/query/quickbooks.sql`

### QBO Connection

```sql
-- name: GetQBConnection :one
SELECT * FROM qb_connection WHERE id = 1;

-- name: UpsertQBConnection :one
INSERT INTO qb_connection (id, realm_id, access_token, refresh_token,
    access_token_expiry, refresh_token_expiry)
VALUES (1, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    realm_id              = EXCLUDED.realm_id,
    access_token          = EXCLUDED.access_token,
    refresh_token         = EXCLUDED.refresh_token,
    access_token_expiry   = EXCLUDED.access_token_expiry,
    refresh_token_expiry  = EXCLUDED.refresh_token_expiry,
    last_refreshed_at     = NOW()
RETURNING *;

-- name: UpdateQBTokens :one
UPDATE qb_connection
SET access_token         = ?,
    refresh_token        = ?,
    access_token_expiry  = ?,
    refresh_token_expiry = ?,
    last_refreshed_at    = NOW()
WHERE id = 1
RETURNING *;

-- name: DeleteQBConnection :exec
DELETE FROM qb_connection WHERE id = 1;
```

### Customer Linking

```sql
-- name: LinkCustomerToQB :one
UPDATE customers
SET qb_customer_id = ?, qb_customer_name = ?
WHERE id = ?
RETURNING *;

-- name: UnlinkCustomerFromQB :one
UPDATE customers
SET qb_customer_id = NULL, qb_customer_name = NULL
WHERE id = ?
RETURNING *;

-- name: ListUnlinkedCustomers :many
SELECT * FROM customers
WHERE qb_customer_id IS NULL
ORDER BY name ASC;

-- name: ListLinkedCustomers :many
SELECT * FROM customers
WHERE qb_customer_id IS NOT NULL
ORDER BY name ASC;
```

### Invoices

```sql
-- name: CreateInvoice :one
INSERT INTO invoices (id, quote_id, amount, due_date, memo)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = ?;

-- name: ListInvoicesByQuote :many
SELECT i.*, COALESCE(SUM(p.amount), 0) AS total_paid
FROM invoices i
LEFT JOIN invoice_payments p ON p.invoice_id = i.id
WHERE i.quote_id = ?
GROUP BY i.id
ORDER BY i.created_at ASC;

-- name: MarkInvoicePushed :one
UPDATE invoices
SET qb_invoice_id  = ?,
    qb_invoice_num = ?,
    qb_sync_token  = ?,
    status         = 'pushed',
    updated_at     = NOW()
WHERE id = ?
RETURNING *;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET status = ?, updated_at = NOW()
WHERE id = ?
RETURNING *;

-- name: UpdateInvoiceSyncToken :one
UPDATE invoices
SET qb_sync_token = ?, updated_at = NOW()
WHERE id = ?
RETURNING *;
```

### Payments

```sql
-- name: CreateInvoicePayment :one
INSERT INTO invoice_payments (id, invoice_id, qb_payment_id, amount, payment_date)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetInvoicePaymentByQBId :one
SELECT * FROM invoice_payments WHERE qb_payment_id = ?;

-- name: ListPaymentsByInvoice :many
SELECT * FROM invoice_payments
WHERE invoice_id = ?
ORDER BY payment_date ASC;
```

Run `sqlc generate`, confirm no errors.

---

## Step 3 — QBO Client

**File:** `internal/qbo/client.go`

A thin wrapper around the QBO REST API. Handles token refresh transparently.
All other parts of the app call this — nothing else touches QBO HTTP directly.

```go
type Client struct {
    db       *db.Queries
    httpClient *http.Client
    baseURL  string  // https://quickbooks.api.intuit.com/v3/company/{realmId}
}

func New(db *db.Queries) *Client

// GetAccessToken returns a valid access token, refreshing if within 5 minutes of expiry.
func (c *Client) GetAccessToken(ctx context.Context) (string, error)

// CreateInvoice pushes a new invoice to QBO and returns the QBO Invoice ID and SyncToken.
func (c *Client) CreateInvoice(ctx context.Context, inv QBInvoice) (id, syncToken string, docNumber string, err error)

// VoidInvoice voids an existing QBO invoice by ID.
func (c *Client) VoidInvoice(ctx context.Context, qbInvoiceID, syncToken string) error

// ListCustomers fetches all customers from QBO for the linking UI.
func (c *Client) ListCustomers(ctx context.Context) ([]QBCustomer, error)
```

### QBO wire types

```go
type QBInvoice struct {
    CustomerRef QBRef
    Line        []QBLineItem
    DueDate     string   // "2026-04-01"
    PrivateNote string   // internal memo
}

type QBLineItem struct {
    Amount      float64
    Description string
    DetailType  string  // "SalesItemLineDetail" — ItemRef deferred
}

type QBCustomer struct {
    ID          string  // QBO value
    DisplayName string
}

type QBRef struct {
    Value string `json:"value"`
    Name  string `json:"name,omitempty"`
}
```

**Note on line items:** Until the line item sync strategy is decided, invoices
are pushed as a single summary line item (total amount, description = quote
name/number). This requires a generic "Construction Services" Item to exist in
QBO — the owner creates this manually once. Replace with proper line items when
that design is complete.

### Token refresh logic

```go
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
    conn, err := c.db.GetQBConnection(ctx)
    // If access token expires within 5 minutes, refresh it
    if time.Until(conn.AccessTokenExpiry) < 5*time.Minute {
        return c.refreshTokens(ctx, conn.RefreshToken)
    }
    return conn.AccessToken, nil
}
```

Token refresh runs inline before each API call. For a single-owner app with
low request volume this is sufficient — no need for a dedicated refresh goroutine.

**Test:** Mock the QBO API, confirm token refresh is called when expiry is
within 5 minutes. Confirm `UpdateQBTokens` is called after a successful refresh.
Confirm a `CreateInvoice` call includes the refreshed token in the Authorization
header.

---

## Step 4 — OAuth Handler

**File:** `handler/qbo/oauth.go`

Handles the one-time authorization flow and connection management.

| Method | Route | Action |
|--------|-------|--------|
| GET | `/settings/quickbooks` | Connection status page |
| GET | `/settings/quickbooks/connect` | Redirect to QBO OAuth consent screen |
| GET | `/settings/quickbooks/callback` | Handle OAuth callback, store tokens |
| POST | `/settings/quickbooks/disconnect` | Delete tokens, unlink connection |

### OAuth flow

1. Owner visits `/settings/quickbooks/connect`
2. Handler generates a random `state` parameter (stored in session), redirects
   to Intuit's authorization URL with `scope=com.intuit.quickbooks.accounting`
3. Owner approves, Intuit redirects to `/settings/quickbooks/callback?code=...&state=...&realmId=...`
4. Handler validates `state`, exchanges `code` for tokens via POST to
   `https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer`
5. Stores tokens via `UpsertQBConnection`, redirects to
   `/settings/quickbooks` with success flash

### Refresh token expiry warning

The connection status page should warn the owner when the refresh token is
within 14 days of expiry (100 day TTL means this is rare, but it will happen
if the app is unused for a long period). Warning: "Your QuickBooks connection
expires on [date]. Visit this page to reconnect."

**Test:** Step through the OAuth flow in QBO sandbox. Confirm tokens are stored.
Confirm disconnect deletes the connection row. Confirm the status page shows
connected state with `realm_id` and `connected_at`.

---

## Step 5 — Customer Linking UI

**File:** `handler/qbo/customers.go`

**Directory:** `templates/qbo/`

Allows the owner to manually link Skalkaho customers to their QBO counterparts.

| Method | Route | Action |
|--------|-------|--------|
| GET | `/settings/quickbooks/customers` | Customer linking page |
| POST | `/settings/quickbooks/customers/{id}/link` | Link customer to QBO customer |
| POST | `/settings/quickbooks/customers/{id}/unlink` | Remove link |

### Linking page

Two columns:

- **Unlinked Skalkaho customers** — list with a QBO customer dropdown next to each
- **Linked customers** — name, QBO display name, unlink button

The QBO customer dropdown is populated by calling `client.ListCustomers()` on
page load. This is a live call to QBO — if the connection is not active, the
page shows a warning and disables linking.

On link submit: store `qb_customer_id` and `qb_customer_name` on the Skalkaho
customer row. HTMX swap moves the customer from the unlinked to the linked list.

**Test:** Link a customer, confirm `qb_customer_id` is stored. Unlink, confirm
it's cleared. Load the page with no QBO connection, confirm graceful degradation
(no dropdown, warning shown).

---

## Step 6 — Invoice Handler

**File:** `handler/invoice/invoice.go`

Manages invoices against a quote and pushes them to QBO.

| Method | Route | Action |
|--------|-------|--------|
| GET | `/quotes/{id}/invoices` | Invoice list for a quote |
| POST | `/quotes/{id}/invoices` | Create draft invoice |
| POST | `/quotes/{id}/invoices/{iid}/push` | Push draft invoice to QBO |
| POST | `/quotes/{id}/invoices/{iid}/void` | Void invoice in QBO |

### Create draft

Accepts `amount` and optional `due_date` and `memo`. Validates that
`amount + existing_invoiced_total <= quote_total`. Stores as `status = 'draft'`
— not yet in QBO.

### Push to QBO

1. Load invoice + quote + customer
2. Confirm customer has `qb_customer_id` — if not, return error with link to
   the customer linking page
3. Call `client.CreateInvoice()`
4. On success → call `MarkInvoicePushed` with returned `qb_invoice_id`,
   `qb_invoice_num`, `qb_sync_token`
5. Return updated invoice partial via HTMX

### Invoice list view

Shows per-quote invoicing summary:

```
Quote total:      $48,500.00
Total invoiced:   $33,950.00  (70%)
Total paid:       $14,550.00  (30%)
Remaining:        $14,550.00

Invoice #1001  $14,550  pushed    paid in full    2026-01-15
Invoice #1002  $19,400  pushed    unpaid          due 2026-03-01
[+ New Invoice]
```

**Test:** Create a draft invoice, confirm it doesn't appear in QBO. Push it,
confirm QBO invoice is created, `qb_invoice_id` is stored. Attempt to create
an invoice that would exceed the quote total, confirm validation error.
Attempt to push an invoice for an unlinked customer, confirm blocked with
helpful error.

---

## Step 7 — Webhook Handler

**File:** `handler/qbo/webhook.go`

Receives payment notifications from QBO and reflects them in Skalkaho.

| Method | Route | Action |
|--------|-------|--------|
| POST | `/webhooks/quickbooks` | Receive QBO webhook events |

### Webhook verification

QBO signs webhook payloads with HMAC-SHA256 using a webhook verifier token
from the Intuit Developer dashboard. Always verify before processing:

```go
func verifyWebhookSignature(payload []byte, signature, verifierToken string) bool {
    mac := hmac.New(sha256.New, []byte(verifierToken))
    mac.Write(payload)
    expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(expected), []byte(signature))
}
```

Signature is in the `intuit-signature` request header. Return 401 if invalid.

### CloudEvents format

Per QBO requirements (mandatory after May 15, 2026), the handler expects
CloudEvents format. Check `Content-Type: application/cloudevents+json`.

### Payment processing

On a `Payment` create or update event:

1. Extract `entityId` (QBO Payment ID) from the event
2. Fetch the full Payment object from QBO API (webhook only contains IDs,
   not the full payload)
3. Extract `TotalAmt`, `TxnDate`, and `Line[].LinkedTxn[].TxnId` (the
   linked QBO invoice ID)
4. Look up the Skalkaho invoice by `qb_invoice_id`
5. If payment already exists (`qb_payment_id` unique) → idempotent, skip
6. Otherwise → insert `invoice_payments` row
7. Recalculate invoice status:
   - `total_paid >= invoice.amount` → `paid`
   - `total_paid > 0` → `partially_paid`
8. Return 200 immediately — do all processing synchronously but keep it fast.
   QBO retries on non-200 responses.

**Test:** Send a test webhook payload (QBO sandbox provides this). Confirm
signature verification rejects tampered payloads. Confirm idempotent — sending
the same payment event twice creates only one `invoice_payments` row. Confirm
invoice status updates correctly after payment.

---

## Step 8 — Token Refresh Background Job

**File:** `worker/jobs/qb_token_refresh.go`

Plugs into the existing worker harness from the material import feature.
Runs on a schedule to proactively refresh the QBO access token before it
expires, rather than waiting for an API call to trigger it inline.

```go
// Registers with the worker to run every 30 minutes.
// Checks if access token expires within 1 hour — if so, refreshes.
// Checks if refresh token expires within 14 days — if so, logs a warning
// (future: send owner notification).
func RefreshQBTokens(ctx context.Context, q *db.Queries) error
```

This is a time-based job rather than a queue-based job — it runs on an interval
regardless of any job row. The worker harness should support a `Schedule`
registration alongside the existing `Register` for queue-based jobs:

```go
w.Schedule("qb_token_refresh", 30*time.Minute, qbjob.RefreshQBTokens)
```

**Test:** Set the access token expiry to 30 minutes from now, confirm the job
refreshes it. Set refresh token expiry to 10 days from now, confirm a warning
is logged.

---

## Step 9 — Settings UI + Router Wiring

**File:** `templates/qbo/`

### `settings.html` (connection status page)

```
QuickBooks Connection

Status:       ● Connected
Company:      {realm_id}
Connected:    January 15, 2026
Token expiry: April 25, 2026  ← warning if < 14 days

[Disconnect]
```

Or if not connected:

```
QuickBooks Connection

Status:       ○ Not connected

[Connect to QuickBooks]
```

### Router wiring

```go
qboHandler := qbo.NewHandler(queries, qboClient)
invHandler := invoice.NewHandler(queries, qboClient)

r.HandleFunc("GET  /settings/quickbooks",                          qboHandler.Settings)
r.HandleFunc("GET  /settings/quickbooks/connect",                  qboHandler.Connect)
r.HandleFunc("GET  /settings/quickbooks/callback",                 qboHandler.Callback)
r.HandleFunc("POST /settings/quickbooks/disconnect",               qboHandler.Disconnect)
r.HandleFunc("GET  /settings/quickbooks/customers",                qboHandler.Customers)
r.HandleFunc("POST /settings/quickbooks/customers/{id}/link",      qboHandler.LinkCustomer)
r.HandleFunc("POST /settings/quickbooks/customers/{id}/unlink",    qboHandler.UnlinkCustomer)
r.HandleFunc("GET  /quotes/{id}/invoices",                         invHandler.Index)
r.HandleFunc("POST /quotes/{id}/invoices",                         invHandler.Create)
r.HandleFunc("POST /quotes/{id}/invoices/{iid}/push",              invHandler.Push)
r.HandleFunc("POST /quotes/{id}/invoices/{iid}/void",              invHandler.Void)
r.HandleFunc("POST /webhooks/quickbooks",                          qboHandler.Webhook)
```

Add a "Settings" section to the nav containing the QuickBooks link. The
webhook route must be **unauthenticated** — QBO cannot log in.

**Test:** Full end-to-end in QBO sandbox — connect, link a customer, create and
push an invoice, trigger a payment webhook, confirm payment appears in Skalkaho
and invoice status updates.

---

## Commit Sequence

| # | Commit message |
|---|----------------|
| 1 | `db: add qb_connection, invoices, invoice_payments migration, qb fields on customers` |
| 2 | `db: add sqlc queries for qbo connection, customer linking, invoices, payments` |
| 3 | `internal/qbo: add QBO API client with token refresh` |
| 4 | `handler/qbo: add OAuth flow (connect, callback, disconnect)` |
| 5 | `handler/qbo: add customer linking UI and routes` |
| 6 | `handler/invoice: add invoice create, push, and void` |
| 7 | `handler/qbo: add webhook handler with signature verification` |
| 8 | `worker: add scheduled job support, add QBO token refresh job` |
| 9 | `templates + router: wire all QBO and invoice routes, add settings nav` |

---

## Deferred — Pending Client Conversations

### Line item sync

When decided, this will require:
- A mapping between Skalkaho material categories and QBO Item IDs
- A new table: `qb_item_mappings (category TEXT, qb_item_id TEXT)`
- Updated `QBInvoice` wire type with per-line-item detail
- A sync UI under `/settings/quickbooks/items`

### Draw schedule

When decided (ad hoc / percentage / milestone), this will shape:
- Whether the `invoices` table needs a `draw_type` or `milestone_id` column
- Whether invoice amounts are calculated or entered manually
- UI for defining and managing the draw schedule per quote

---

## Environment Variables Required

```
QBO_CLIENT_ID        -- from Intuit Developer dashboard
QBO_CLIENT_SECRET    -- from Intuit Developer dashboard
QBO_REDIRECT_URI     -- must match exactly what's registered in Intuit dashboard
QBO_WEBHOOK_TOKEN    -- verifier token for webhook signature validation
QBO_SANDBOX          -- "true" to use sandbox base URL
```
