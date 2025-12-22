# Go Coding Style Guide

A comprehensive style guide for Go web applications covering full-stack apps (Go + PostgreSQL + htmx + Alpine.js) and Hugo static sites.

---

## Table of Contents

1. [Go Code Style](#section-1-go-code-style)
2. [Database Conventions](#section-2-database-conventions)
3. [HTTP Layer](#section-3-http-layer)
4. [Frontend Patterns](#section-4-frontend-patterns)
5. [Static Site Development](#section-5-static-site-development)
6. [Testing](#section-6-testing)
7. [Deployment & Infrastructure](#section-7-deployment--infrastructure)
8. [Project Documentation](#section-8-project-documentation)

---

## Section 1: Go Code Style

### 1.1 Naming Conventions

#### Variables and Functions

```go
// camelCase for variables and unexported functions
userID := "123"
func parseConfig() {}

// PascalCase for exported functions and types
func ParseConfig() {}
type UserService struct {}
```

#### Packages

```go
// Lowercase, single word preferred
package repository
package handler
package middleware

// Avoid generic names
package utils    // ❌ Too vague
package helpers  // ❌ Too vague
```

#### Interfaces

```go
// Single-method interfaces: verb + "er"
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Multi-method interfaces: descriptive noun
type Storage interface {
    Get(key string) ([]byte, error)
    Put(key string, value []byte) error
    Delete(key string) error
}
```

#### Acronyms

```go
// Keep acronyms uppercase
userID := "123"      // Not userId
httpClient := &http.Client{}
urlPath := "/api"    // Not urlpath
```

### 1.2 Error Handling

#### Wrap Errors with Context

```go
if err != nil {
    return nil, fmt.Errorf("GetProduct(%s): %w", id, err)
}
```

#### Use Early Returns

```go
// ✅ Good: Early returns
func ProcessOrder(order *Order) error {
    if order == nil {
        return errors.New("order is nil")
    }
    
    if order.Total <= 0 {
        return errors.New("order total must be positive")
    }
    
    // Happy path continues
    return saveOrder(order)
}

// ❌ Bad: Deep nesting
func ProcessOrder(order *Order) error {
    if order != nil {
        if order.Total > 0 {
            return saveOrder(order)
        } else {
            return errors.New("order total must be positive")
        }
    } else {
        return errors.New("order is nil")
    }
}
```

#### Domain Errors

Define typed errors for consistent handling:

```go
// internal/domain/error.go
package domain

// Error codes for mapping to HTTP status
const (
    EINVALID      = "invalid"       // 400
    EUNAUTHORIZED = "unauthorized"  // 401
    EFORBIDDEN    = "forbidden"     // 403
    ENOTFOUND     = "not_found"     // 404
    ECONFLICT     = "conflict"      // 409
    EINTERNAL     = "internal"      // 500
)

type Error struct {
    Code    string
    Op      string
    Message string
    Err     error
}

func (e *Error) Error() string {
    return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *Error) Unwrap() error {
    return e.Err
}

func Errorf(code, op, message string, args ...interface{}) *Error {
    return &Error{
        Code:    code,
        Op:      op,
        Message: fmt.Sprintf(message, args...),
    }
}

func ErrorCode(err error) string {
    if err == nil {
        return ""
    }
    var e *Error
    if errors.As(err, &e) {
        return e.Code
    }
    return EINTERNAL
}

func ErrorMessage(err error) string {
    if err == nil {
        return ""
    }
    var e *Error
    if errors.As(err, &e) {
        return e.Message
    }
    return "An unexpected error occurred"
}
```

### 1.3 Package Organization

#### Project Structure (Full-Stack)

```
project-root/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, wires dependencies
├── internal/
│   ├── config/                  # Environment configuration
│   ├── database/                # Connection and migrations
│   ├── domain/                  # Business logic, validation, errors
│   ├── handler/
│   │   ├── admin/               # Admin handlers
│   │   ├── storefront/          # Customer-facing handlers
│   │   └── webhook/             # External integrations
│   ├── middleware/              # HTTP middleware
│   ├── repository/              # sqlc generated code
│   ├── router/                  # Route definitions
│   └── templates/               # html/template files
├── migrations/                  # Goose SQL migrations
├── static/                      # CSS, JS, images
├── Makefile
├── sqlc.yaml
└── CLAUDE.md
```

#### Organize by User Context, Not Entity

```go
// ✅ Good: Handlers grouped by who uses them
internal/handler/
├── admin/           # Tenant admin operations
│   ├── products.go
│   ├── orders.go
│   └── customers.go
├── storefront/      # Customer-facing pages
│   ├── catalog.go
│   ├── cart.go
│   └── checkout.go
└── webhook/         # External service callbacks
    ├── stripe.go
    └── shipping.go

// ❌ Bad: Handlers grouped by entity
internal/handler/
├── products.go      # Mixes admin and storefront logic
├── orders.go
└── customers.go
```

### 1.4 Interface Design

#### Accept Interfaces, Return Structs

```go
// ✅ Good: Accept interface
func NewProductService(store ProductStore) *ProductService {
    return &ProductService{store: store}
}

// ❌ Bad: Accept concrete type
func NewProductService(store *PostgresProductStore) *ProductService {
    return &ProductService{store: store}
}
```

#### Define Interfaces Where Used

```go
// internal/handler/admin/products.go

// Define interface in the consumer package
type ProductStore interface {
    GetProduct(ctx context.Context, tenantID, productID uuid.UUID) (*domain.Product, error)
    ListProducts(ctx context.Context, tenantID uuid.UUID) ([]domain.Product, error)
}

type ProductHandler struct {
    store ProductStore
}
```

### 1.5 Struct Initialization

#### Use Named Fields

```go
// ✅ Good: Named fields
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
}

// ❌ Bad: Positional (fragile if struct changes)
server := &http.Server{":8080", mux, 15 * time.Second, 15 * time.Second}
```

### 1.6 Comments

```go
// Package repository provides database access using sqlc-generated queries.
package repository

// ProductService handles product-related business operations.
// It coordinates between the repository layer and external services.
type ProductService struct {
    // ...
}

// GetProduct retrieves a product by ID within a tenant's scope.
// Returns ENOTFOUND if the product doesn't exist.
func (s *ProductService) GetProduct(ctx context.Context, tenantID, productID uuid.UUID) (*Product, error) {
    // ...
}
```

**When to comment:**
- Package documentation (required)
- Exported types and functions
- Non-obvious "why" decisions
- Complex algorithms

**When not to comment:**
- Obvious code (`i++ // increment i`)
- Self-documenting names

---

## Section 2: Database Conventions

### 2.1 Table Naming

```sql
-- Plural, snake_case
CREATE TABLE products (...);
CREATE TABLE order_items (...);
CREATE TABLE price_list_entries (...);

-- Junction tables: alphabetical order
CREATE TABLE product_categories (...);  -- Not categories_products
```

### 2.2 Column Naming

```sql
CREATE TABLE products (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Foreign keys: singular_table_id
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    category_id UUID REFERENCES categories(id),
    
    -- Data columns: snake_case
    name TEXT NOT NULL,
    base_price_cents INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- Timestamps (always include)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 2.3 Multi-Tenant Pattern

Every tenant-scoped table includes `tenant_id`:

```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    -- ...
    
    -- Unique constraints include tenant_id
    UNIQUE(tenant_id, slug)
);

-- Composite index for tenant queries
CREATE INDEX idx_products_tenant ON products(tenant_id);
CREATE INDEX idx_products_tenant_status ON products(tenant_id, status);
```

All queries include tenant_id filter:

```sql
-- name: GetProduct :one
SELECT * FROM products
WHERE tenant_id = $1 AND id = $2;

-- name: ListProducts :many
SELECT * FROM products
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
```

### 2.4 Migration Structure

Use Goose with numbered migrations:

```
migrations/
├── 00001_create_tenants.sql
├── 00002_create_products.sql
├── 00003_create_customers.sql
├── 00004_add_product_metadata.sql
└── 00005_create_orders.sql
```

#### Migration Template

```sql
-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    base_price_cents INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    UNIQUE(tenant_id, slug)
);

CREATE INDEX idx_products_tenant ON products(tenant_id);
CREATE INDEX idx_products_tenant_status ON products(tenant_id, status);

-- +goose Down
DROP TABLE IF EXISTS products;
```

### 2.5 sqlc Configuration

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sqlc/queries/"
    schema: "migrations/"
    gen:
      go:
        package: "repository"
        out: "internal/repository"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        overrides:
          - db_type: "uuid"
            go_type:
              import: "github.com/google/uuid"
              type: "UUID"
          - db_type: "timestamptz"
            go_type:
              import: "time"
              type: "Time"
```

### 2.6 sqlc Query Patterns

```sql
-- sqlc/queries/products.sql

-- name: GetProduct :one
SELECT * FROM products
WHERE tenant_id = $1 AND id = $2;

-- name: GetProductBySlug :one
SELECT * FROM products
WHERE tenant_id = $1 AND slug = $2;

-- name: ListProducts :many
SELECT * FROM products
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveProducts :many
SELECT * FROM products
WHERE tenant_id = $1 AND status = 'active'
ORDER BY name ASC;

-- name: CreateProduct :one
INSERT INTO products (
    tenant_id, name, slug, description, base_price_cents, status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateProduct :one
UPDATE products SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    base_price_cents = COALESCE(sqlc.narg('base_price_cents'), base_price_cents),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = now()
WHERE tenant_id = $1 AND id = $2
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE tenant_id = $1 AND id = $2;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE tenant_id = $1;
```

### 2.7 Money Handling

Store money as integers (cents):

```sql
base_price_cents INTEGER NOT NULL DEFAULT 0,
discount_cents INTEGER NOT NULL DEFAULT 0,
total_cents INTEGER NOT NULL,
```

Format in Go:

```go
func FormatMoney(cents int, currency string) string {
    dollars := float64(cents) / 100
    switch currency {
    case "USD":
        return fmt.Sprintf("$%.2f", dollars)
    case "EUR":
        return fmt.Sprintf("€%.2f", dollars)
    default:
        return fmt.Sprintf("%.2f %s", dollars, currency)
    }
}
```

### 2.8 Status Enums

Use CHECK constraints instead of enum types:

```sql
status TEXT NOT NULL DEFAULT 'draft' 
    CHECK (status IN ('draft', 'active', 'archived')),
```

**Why:** Easier to add new values without migrations.

### 2.9 JSONB Metadata

For flexible, schemaless data:

```sql
metadata JSONB NOT NULL DEFAULT '{}',
settings JSONB NOT NULL DEFAULT '{}',
```

Query patterns:

```sql
-- name: GetProductsWithTag :many
SELECT * FROM products
WHERE tenant_id = $1 
  AND metadata->>'featured' = 'true';

-- name: UpdateProductMetadata :one
UPDATE products 
SET metadata = metadata || $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2
RETURNING *;
```

### 2.10 Makefile Targets

```makefile
# Database
DB_URL ?= postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable

.PHONY: db-up
db-up:
	docker compose up -d postgres

.PHONY: db-migrate
db-migrate:
	goose -dir migrations postgres "$(DB_URL)" up

.PHONY: db-rollback
db-rollback:
	goose -dir migrations postgres "$(DB_URL)" down

.PHONY: db-reset
db-reset:
	goose -dir migrations postgres "$(DB_URL)" reset
	goose -dir migrations postgres "$(DB_URL)" up

.PHONY: db-status
db-status:
	goose -dir migrations postgres "$(DB_URL)" status

.PHONY: db-new
db-new:
	@read -p "Migration name: " name; \
	goose -dir migrations create $$name sql

.PHONY: sqlc
sqlc:
	sqlc generate
```

---

## Section 3: HTTP Layer

### 3.1 Handler Structure

#### Full-Stack Handler Pattern

```go
// internal/handler/admin/products.go
package admin

type ProductHandler struct {
    queries  *repository.Queries
    renderer *templates.Renderer
    storage  storage.Client
    logger   *slog.Logger
}

func NewProductHandler(
    queries *repository.Queries,
    renderer *templates.Renderer,
    storage storage.Client,
    logger *slog.Logger,
) *ProductHandler {
    return &ProductHandler{
        queries:  queries,
        renderer: renderer,
        storage:  storage,
        logger:   logger,
    }
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    tenantID := domain.TenantIDFromContext(ctx)
    logger := middleware.LoggerFromContext(ctx)
    
    products, err := h.queries.ListProducts(ctx, repository.ListProductsParams{
        TenantID: tenantID,
        Limit:    50,
        Offset:   0,
    })
    if err != nil {
        logger.Error("failed to list products", "error", err)
        h.renderError(w, r, err)
        return
    }
    
    h.render(w, r, "admin/products", map[string]any{
        "Products": products,
    })
}
```

#### Minimal API Pattern (Static Sites)

```go
// api/main.go
package main

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    allowedOrigins := parseOrigins(os.Getenv("ALLOWED_ORIGINS"))

    mux := http.NewServeMux()
    mux.HandleFunc("POST /api/contact", handleContact)
    mux.HandleFunc("GET /api/health", handleHealth)

    handler := corsMiddleware(mux, allowedOrigins)

    log.Printf("Starting API server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
```

### 3.2 Routing (Go 1.22+)

Use method-prefixed patterns:

```go
mux := http.NewServeMux()

// Method + path
mux.HandleFunc("GET /products", handler.ListProducts)
mux.HandleFunc("POST /products", handler.CreateProduct)
mux.HandleFunc("GET /products/{id}", handler.GetProduct)
mux.HandleFunc("PUT /products/{id}", handler.UpdateProduct)
mux.HandleFunc("DELETE /products/{id}", handler.DeleteProduct)

// Access path parameters
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    // ...
}
```

#### Optional Router Wrapper

For middleware support similar to chi:

```go
// internal/router/router.go
package router

type Router struct {
    mux        *http.ServeMux
    middleware []func(http.Handler) http.Handler
}

func New() *Router {
    return &Router{mux: http.NewServeMux()}
}

func (r *Router) Use(mw func(http.Handler) http.Handler) {
    r.middleware = append(r.middleware, mw)
}

func (r *Router) Get(pattern string, handler http.HandlerFunc) {
    r.mux.HandleFunc("GET "+pattern, handler)
}

func (r *Router) Post(pattern string, handler http.HandlerFunc) {
    r.mux.HandleFunc("POST "+pattern, handler)
}

func (r *Router) Put(pattern string, handler http.HandlerFunc) {
    r.mux.HandleFunc("PUT "+pattern, handler)
}

func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
    r.mux.HandleFunc("DELETE "+pattern, handler)
}

func (r *Router) Group(fn func(*Router)) {
    fn(r)
}

func (r *Router) Handler() http.Handler {
    var handler http.Handler = r.mux
    for i := len(r.middleware) - 1; i >= 0; i-- {
        handler = r.middleware[i](handler)
    }
    return handler
}
```

### 3.3 Middleware

#### Standard Signature

```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before
        next.ServeHTTP(w, r)
        // After
    })
}
```

#### Middleware Order (outermost to innermost)

```go
handler := router.Handler()
handler = middleware.Recover(handler)       // 1. Catch panics
handler = middleware.RequestID(handler)     // 2. Assign request ID
handler = middleware.Logger(handler)        // 3. Log requests
handler = middleware.SecureHeaders(handler) // 4. Security headers
handler = middleware.RateLimit(handler)     // 5. Rate limiting
handler = middleware.ResolveTenant(handler) // 6. Multi-tenant
// Authentication applied per-route, not globally
```

#### Context Helpers

```go
// internal/middleware/context.go
package middleware

type contextKey string

const (
    requestIDKey contextKey = "requestID"
    loggerKey    contextKey = "logger"
    tenantIDKey  contextKey = "tenantID"
    userKey      contextKey = "user"
)

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}
```

### 3.4 CORS Middleware

```go
func corsMiddleware(next http.Handler, allowedOrigins map[string]bool) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Vary", "Origin")
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Parse comma-separated origins
func parseOrigins(s string) map[string]bool {
    origins := make(map[string]bool)
    for _, origin := range strings.Split(s, ",") {
        origin = strings.TrimSpace(origin)
        if origin != "" {
            origins[origin] = true
        }
    }
    return origins
}
```

### 3.5 Response Formatting

#### Error to HTTP Status Mapping

```go
func httpStatusFromError(err error) int {
    switch domain.ErrorCode(err) {
    case domain.EINVALID:
        return http.StatusBadRequest
    case domain.EUNAUTHORIZED:
        return http.StatusUnauthorized
    case domain.EFORBIDDEN:
        return http.StatusForbidden
    case domain.ENOTFOUND:
        return http.StatusNotFound
    case domain.ECONFLICT:
        return http.StatusConflict
    default:
        return http.StatusInternalServerError
    }
}
```

#### JSON vs HTML Response

```go
func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, err error) {
    status := httpStatusFromError(err)
    message := domain.ErrorMessage(err)
    
    if r.Header.Get("Accept") == "application/json" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(map[string]string{
            "error": message,
        })
        return
    }
    
    // HTML response
    w.WriteHeader(status)
    h.renderer.Render(w, "error", map[string]any{
        "Status":  status,
        "Message": message,
    })
}
```

### 3.6 Form Handling

#### JSON Input

```go
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var input CreateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        h.badRequest(w, "Invalid JSON")
        return
    }
    
    if errors := input.Validate(); len(errors) > 0 {
        h.validationError(w, errors)
        return
    }
    
    // Process...
}
```

#### Form Input

```go
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        h.badRequest(w, "Invalid form data")
        return
    }
    
    input := CreateProductInput{
        Name:        r.FormValue("name"),
        Slug:        r.FormValue("slug"),
        Description: r.FormValue("description"),
    }
    
    if errors := input.Validate(); len(errors) > 0 {
        h.validationError(w, errors)
        return
    }
    
    // Process...
}
```

#### Validation Pattern

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type CreateProductInput struct {
    Name           string `json:"name"`
    Slug           string `json:"slug"`
    Description    string `json:"description"`
    BasePriceCents int    `json:"base_price_cents"`
}

func (i *CreateProductInput) Validate() []ValidationError {
    var errors []ValidationError
    
    if strings.TrimSpace(i.Name) == "" {
        errors = append(errors, ValidationError{
            Field:   "name",
            Message: "Name is required",
        })
    } else if len(i.Name) > 255 {
        errors = append(errors, ValidationError{
            Field:   "name",
            Message: "Name must be less than 255 characters",
        })
    }
    
    if !slugRegex.MatchString(i.Slug) {
        errors = append(errors, ValidationError{
            Field:   "slug",
            Message: "Slug must contain only lowercase letters, numbers, and hyphens",
        })
    }
    
    if i.BasePriceCents < 0 {
        errors = append(errors, ValidationError{
            Field:   "base_price_cents",
            Message: "Price cannot be negative",
        })
    }
    
    return errors
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
```

### 3.7 Spam Protection

#### Honeypot Field

```html
<!-- Hidden field that bots fill but humans don't see -->
<div class="hidden" aria-hidden="true">
    <label for="website">Website</label>
    <input type="text" name="website" id="website" tabindex="-1" autocomplete="off">
</div>
```

```go
func handleContact(w http.ResponseWriter, r *http.Request) {
    var form ContactForm
    json.NewDecoder(r.Body).Decode(&form)
    
    // Check honeypot - if filled, it's a bot
    if strings.TrimSpace(form.Website) != "" {
        // Return fake success to not alert the bot
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]bool{"success": true})
        return
    }
    
    // Process legitimate submission...
}
```

#### Cloudflare Turnstile

```go
type TurnstileResponse struct {
    Success    bool     `json:"success"`
    ErrorCodes []string `json:"error-codes,omitempty"`
}

func verifyTurnstile(token, secretKey, remoteIP string) (bool, error) {
    if secretKey == "" {
        return true, nil // Skip in development
    }

    payload := map[string]string{
        "secret":   secretKey,
        "response": token,
        "remoteip": remoteIP,
    }
    
    body, _ := json.Marshal(payload)
    resp, err := http.Post(
        "https://challenges.cloudflare.com/turnstile/v0/siteverify",
        "application/json",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    var result TurnstileResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result.Success, nil
}
```

**Development test keys:**
- Site Key: `1x00000000000000000000AA`
- Secret: `1x0000000000000000000000000000000AA`

### 3.8 Health Check

```go
// Simple health check
func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}

// Full health check with database
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    if err := h.db.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "unhealthy",
            "error":  "database unavailable",
        })
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
}
```

---

## Section 4: Frontend Patterns

### 4.1 Alpine.js Conventions

#### Loading

```html
<!-- Load at end of body, collapse plugin first -->
<script defer src="https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
```

#### Component Pattern

```javascript
function contactForm() {
    return {
        // State
        formData: {
            firstName: '',
            lastName: '',
            email: '',
            message: ''
        },
        isSubmitting: false,
        formError: '',
        
        // Lifecycle
        init() {
            // Runs when component initializes
            const urlParams = new URLSearchParams(window.location.search);
            this.formData.service = urlParams.get('service') || '';
        },
        
        // Methods
        async submitForm() {
            this.isSubmitting = true;
            this.formError = '';
            
            try {
                const response = await fetch('/api/contact', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(this.formData)
                });
                
                const result = await response.json();
                
                if (response.ok && result.success) {
                    window.location.href = '/success';
                } else {
                    this.formError = result.error || 'Something went wrong';
                }
            } catch (error) {
                this.formError = 'Unable to submit form';
            } finally {
                this.isSubmitting = false;
            }
        }
    };
}
```

```html
<form x-data="contactForm()" @submit.prevent="submitForm">
    <div x-show="formError" class="error-message" x-text="formError"></div>
    
    <input type="text" x-model="formData.firstName" required>
    <input type="email" x-model="formData.email" required>
    
    <button type="submit" :disabled="isSubmitting">
        <span x-text="isSubmitting ? 'Sending...' : 'Send Message'"></span>
    </button>
</form>
```

#### Common Patterns

**Mobile Menu:**

```html
<div x-data="{ mobileMenuOpen: false }">
    <button @click="mobileMenuOpen = true">Menu</button>
    
    <div x-show="mobileMenuOpen" 
         x-transition:enter="transition ease-out duration-200"
         x-transition:leave="transition ease-in duration-150"
         @click.outside="mobileMenuOpen = false"
         @keydown.escape.window="mobileMenuOpen = false">
        <!-- Menu content -->
    </div>
</div>
```

**Dismissible Banner:**

```html
<div x-data="{ showBanner: true }" x-show="showBanner">
    <p>Announcement message</p>
    <button @click="showBanner = false">Dismiss</button>
</div>
```

### 4.2 htmx Conventions

#### Core Attributes

```html
<!-- GET request, replace inner HTML -->
<button hx-get="/products" hx-target="#product-list">Load Products</button>

<!-- POST request, replace entire element -->
<form hx-post="/products" hx-swap="outerHTML">
    <input name="name" required>
    <button type="submit">Create</button>
</form>

<!-- DELETE with confirmation -->
<button hx-delete="/products/123" 
        hx-confirm="Delete this product?" 
        hx-target="closest tr" 
        hx-swap="outerHTML swap:1s">
    Delete
</button>
```

#### Swap Strategies

| Strategy | Behavior |
|----------|----------|
| `innerHTML` | Replace children (default) |
| `outerHTML` | Replace entire element |
| `beforebegin` | Insert before element |
| `afterbegin` | Insert as first child |
| `beforeend` | Insert as last child |
| `afterend` | Insert after element |
| `delete` | Remove element |
| `none` | No swap (for side effects) |

#### Go Handler Pattern

```go
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    // ... fetch product ...
    
    // Check if htmx request
    if r.Header.Get("HX-Request") == "true" {
        // Return partial HTML
        h.renderer.RenderPartial(w, "products/_detail", product)
        return
    }
    
    // Return full page
    h.renderer.Render(w, "products/show", map[string]any{
        "Product": product,
    })
}
```

#### Response Headers

```go
// Redirect (client-side)
w.Header().Set("HX-Redirect", "/products")

// Refresh entire page
w.Header().Set("HX-Refresh", "true")

// Change target
w.Header().Set("HX-Retarget", "#notification-area")

// Trigger client-side event
w.Header().Set("HX-Trigger", `{"showToast": {"message": "Product saved"}}`)
```

### 4.3 Alpine.js + htmx Integration

**Optimistic UI with confirmation:**

```html
<form hx-post="/products" 
      hx-target="#product-list" 
      hx-swap="beforeend"
      x-data="{ isSubmitting: false }"
      @submit="isSubmitting = true"
      @htmx:after-request="isSubmitting = false">
    
    <input name="name" required :disabled="isSubmitting">
    <button type="submit" :disabled="isSubmitting">
        <span x-show="!isSubmitting">Add Product</span>
        <span x-show="isSubmitting">Adding...</span>
    </button>
</form>
```

**Toast Notifications:**

```html
<div x-data="toastManager()" @show-toast.window="addToast($event.detail)">
    <template x-for="toast in toasts" :key="toast.id">
        <div x-show="toast.visible"
             x-transition
             class="toast"
             :class="toast.type">
            <span x-text="toast.message"></span>
        </div>
    </template>
</div>

<script>
function toastManager() {
    return {
        toasts: [],
        addToast(detail) {
            const id = Date.now();
            this.toasts.push({ id, ...detail, visible: true });
            setTimeout(() => this.removeToast(id), 3000);
        },
        removeToast(id) {
            const toast = this.toasts.find(t => t.id === id);
            if (toast) toast.visible = false;
            setTimeout(() => {
                this.toasts = this.toasts.filter(t => t.id !== id);
            }, 300);
        }
    };
}
</script>
```

### 4.4 Tailwind CSS Conventions

#### Class Ordering

Follow consistent order: Layout → Position → Sizing → Spacing → Typography → Background → Border → Effects → Transitions → States

```html
<button class="flex items-center justify-center 
               w-full max-w-md 
               px-4 py-2 
               text-sm font-medium text-white 
               bg-blue-600 
               rounded-lg border border-blue-700 
               shadow-sm 
               transition-colors 
               hover:bg-blue-500 focus:ring-2 focus:ring-blue-500">
    Submit
</button>
```

#### Responsive Design

Mobile-first approach:

```html
<!-- Base (mobile) → sm → md → lg → xl → 2xl -->
<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
    <!-- Cards -->
</div>

<h1 class="text-2xl sm:text-3xl lg:text-4xl">Responsive Heading</h1>
```

#### Component Patterns

```html
<!-- Primary Button -->
<button class="rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
    Primary Action
</button>

<!-- Secondary Button -->
<button class="rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
    Secondary Action
</button>

<!-- Card -->
<div class="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
    <h3 class="text-lg font-semibold text-gray-900">Card Title</h3>
    <p class="mt-2 text-sm text-gray-600">Card content</p>
</div>

<!-- Form Input -->
<input class="block w-full rounded-md border-0 py-2 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm">
```

### 4.5 Custom Typography System

Define CSS variables for consistent styling:

```css
/* assets/css/main.css */
:root {
  /* Font Families */
  --font-primary: 'Manrope', sans-serif;
  --font-display: 'Outfit', sans-serif;

  /* Type Scale */
  --text-kicker: 0.8125rem;    /* 13px */
  --text-caption: 0.875rem;    /* 14px */
  --text-body: 1rem;           /* 16px */
  --text-body-lg: 1.125rem;    /* 18px */
  --text-subhead: 1.25rem;     /* 20px */
  --text-headline-sm: 1.5rem;  /* 24px */
  --text-headline: 2rem;       /* 32px */
  --text-headline-lg: 2.75rem; /* 44px */
  --text-display-sm: 3.5rem;   /* 56px */
  --text-display: 4.5rem;      /* 72px */

  /* Line Heights */
  --leading-tight: 1.15;
  --leading-snug: 1.25;
  --leading-normal: 1.45;
  --leading-relaxed: 1.55;

  /* Brand Colors */
  --primary-600: #3a61d6;
  --secondary-600: #417848;
}

/* Utility Classes */
.text-headline {
  font-size: var(--text-headline);
  line-height: var(--leading-tight);
}

.font-display-semibold {
  font-family: var(--font-display);
  font-weight: 600;
}

.text-primary-600 {
  color: var(--primary-600);
}
```

### 4.6 Native Browser APIs First

Prefer platform features before JavaScript:

```html
<!-- Modal with <dialog> -->
<dialog id="my-modal">
    <h2>Modal Title</h2>
    <p>Modal content</p>
    <button onclick="this.closest('dialog').close()">Close</button>
</dialog>
<button onclick="document.getElementById('my-modal').showModal()">Open Modal</button>

<!-- Accordion with <details> -->
<details>
    <summary class="cursor-pointer font-semibold">Question?</summary>
    <p class="mt-2">Answer content here.</p>
</details>

<!-- Native Form Validation -->
<input type="email" required pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$">
```

**Decision flow:**
1. HTML (semantic elements)
2. CSS (`:hover`, `:focus`, transitions)
3. Native APIs (`<dialog>`, `popover`, `<details>`)
4. Alpine.js (client-side state)
5. htmx (server communication)
6. Vanilla JS (complex interactions)

---

## Section 5: Static Site Development

### 5.1 Hugo Project Structure

```
project-root/
├── assets/
│   └── css/main.css         # Tailwind + custom CSS
├── content/
│   ├── _index.md            # Homepage
│   ├── about.md
│   ├── services.md
│   ├── contact.md
│   └── privacy.md
├── data/
│   ├── navigation.yaml
│   ├── services.yaml
│   ├── faq.yaml
│   └── team.yaml
├── layouts/
│   ├── _default/
│   │   ├── baseof.html      # Base template
│   │   └── single.html      # Default single page
│   ├── partials/
│   │   ├── head.html
│   │   ├── header.html
│   │   ├── footer.html
│   │   ├── hero.html
│   │   └── ...
│   └── index.html           # Homepage layout
├── static/
│   ├── images/
│   ├── favicon.ico
│   └── ...
├── hugo.toml
├── package.json
└── tailwind.config.js
```

### 5.2 Configuration (hugo.toml)

```toml
baseURL = 'https://www.example.com/'
languageCode = 'en-us'
title = 'Company Name'

[params]
  description = 'Company description for SEO'
  email = 'contact@example.com'
  phone = '(555) 123-4567'
  founder = 'Jane Doe'
  founderTitle = 'Founder & CEO'

  [params.address]
    locality = 'City'
    region = 'ST'
    country = 'US'

  [params.ogImage]
    url = 'https://example.com/og-image.jpg'
    width = 1200
    height = 630

  # Cloudflare Turnstile
  turnstileSiteKey = '0x4AAAAAAA...'

[markup]
  [markup.goldmark]
    [markup.goldmark.renderer]
      unsafe = true
```

### 5.3 Base Template

```html
<!-- layouts/_default/baseof.html -->
<!DOCTYPE html>
<html lang="en">
<head>
  {{ partial "head.html" . }}
</head>
<body class="bg-white">
  {{ partial "header.html" . }}

  <main>
    {{ block "main" . }}{{ end }}
  </main>

  {{ partial "footer.html" . }}

  <!-- Alpine.js -->
  <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
  
  <!-- Analytics -->
  <script defer data-domain="example.com" src="https://plausible.io/js/script.js"></script>
</body>
</html>
```

### 5.4 Head Partial

```html
<!-- layouts/partials/head.html -->
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">

<title>{{ if .IsHome }}{{ .Site.Title }}{{ else }}{{ .Title }} | {{ .Site.Title }}{{ end }}</title>
<meta name="description" content="{{ with .Description }}{{ . }}{{ else }}{{ .Site.Params.description }}{{ end }}">

<!-- Canonical URL -->
<link rel="canonical" href="{{ .Permalink }}">

<!-- Open Graph -->
<meta property="og:type" content="website">
<meta property="og:url" content="{{ .Permalink }}">
<meta property="og:title" content="{{ .Title }}">
<meta property="og:description" content="{{ .Description | default .Site.Params.description }}">
<meta property="og:image" content="{{ .Site.Params.ogImage.url }}">

<!-- Twitter -->
<meta property="twitter:card" content="summary_large_image">

<!-- Favicon -->
<link rel="icon" type="image/x-icon" href="/favicon.ico">

<!-- Non-blocking Font Loading -->
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="preload" as="style" href="https://fonts.googleapis.com/css2?family=...">
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=..." media="print" onload="this.media='all'">

<!-- CSS with Fingerprinting -->
{{ $css := resources.Get "css/main.css" }}
{{ if $css }}
  {{ $css = $css | css.PostCSS | minify | fingerprint }}
  <link rel="stylesheet" href="{{ $css.RelPermalink }}" integrity="{{ $css.Data.Integrity }}">
{{ end }}
```

### 5.5 Data Files

```yaml
# data/navigation.yaml
main:
  - name: Services
    href: /services
  - name: About
    href: /about
  - name: Contact
    href: /contact

footer:
  - name: Privacy
    href: /privacy
```

```yaml
# data/services.yaml
packages:
  - id: essentials
    name: Essentials Package
    price: "$750"
    period: "/month"
    features:
      - "Feature one"
      - "Feature two"
    ctaText: "Get Started"
    ctaHref: "/contact?service=essentials"
```

### 5.6 Template Patterns

```html
<!-- Range over data -->
{{ range site.Data.services.packages }}
  <div class="package">
    <h3>{{ .name }}</h3>
    <p>{{ .price }}{{ .period }}</p>
    <ul>
      {{ range .features }}
        <li>{{ . }}</li>
      {{ end }}
    </ul>
    <a href="{{ .ctaHref }}">{{ .ctaText }}</a>
  </div>
{{ end }}

<!-- Conditionals -->
{{ if .popular }}
  <span class="badge">Most Popular</span>
{{ end }}

<!-- With default -->
{{ with .Description }}
  <p>{{ . }}</p>
{{ else }}
  <p>{{ .Site.Params.description }}</p>
{{ end }}

<!-- Reusable partial with parameters -->
{{ partial "section-header.html" (dict 
  "title" "Our Services" 
  "description" "What we offer"
  "kicker" "Services"
) }}
```

### 5.7 Structured Data (JSON-LD)

```html
<!-- Homepage -->
{{ if .IsHome }}
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Organization",
  "name": "{{ .Site.Title }}",
  "description": "{{ .Site.Params.description }}",
  "url": "{{ .Site.BaseURL }}",
  "founder": {
    "@type": "Person",
    "name": "{{ .Site.Params.founder }}",
    "jobTitle": "{{ .Site.Params.founderTitle }}"
  },
  "contactPoint": {
    "@type": "ContactPoint",
    "telephone": "{{ .Site.Params.phone }}",
    "email": "{{ .Site.Params.email }}"
  }
}
</script>
{{ end }}
```

### 5.8 Image Optimization

**Development:** Use images as-is for fast iteration.

**Pre-deployment:** Run optimization pass:

```makefile
.PHONY: optimize-images
optimize-images:
	@echo "Optimizing images..."
	@find static/images -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | while read img; do \
		npx sharp-cli --input "$$img" --output "$${img%.*}.webp" --format webp --quality 80; \
	done
```

**Template pattern:**

```html
<picture>
  <source srcset="/images/hero.webp" type="image/webp">
  <img src="/images/hero.jpg" alt="Description" loading="lazy" width="1600" height="900">
</picture>
```

**Guidelines:**
- Max 200KB per image
- WebP format preferred
- Max 1600px width for hero images
- Always include `loading="lazy"` (except above-fold)
- Always include `width` and `height` attributes

---

## Section 6: Testing

### 6.1 Testing Philosophy

**Focus on domain logic:**
- Validation rules
- Business calculations
- Data transformations
- Error handling logic
- Authentication/authorization rules

**Skip:**
- Repository layer (trust sqlc + database constraints)
- HTTP handler plumbing (trust stdlib)
- External service integrations
- Thin middleware wrappers

**Why domain-only:**
- No database or network required
- Fast execution (milliseconds)
- No mocking complexity
- Tests the code that matters most

### 6.2 Project Structure

Keep tests alongside the code they test:

```
internal/
└── domain/
    ├── auth.go
    ├── auth_test.go
    ├── error.go
    ├── error_test.go
    ├── money.go
    ├── money_test.go
    ├── validation.go
    └── validation_test.go
```

### 6.3 Validation Tests

```go
// internal/domain/validation_test.go
package domain_test

import (
	"testing"

	"myapp/internal/domain"
)

func TestCreateProductInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.CreateProductInput
		wantErr bool
		field   string
	}{
		{
			name: "valid input",
			input: domain.CreateProductInput{
				Name:           "Test Product",
				Slug:           "test-product",
				BasePriceCents: 1999,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.CreateProductInput{
				Name: "",
				Slug: "test-product",
			},
			wantErr: true,
			field:   "name",
		},
		{
			name: "invalid slug format",
			input: domain.CreateProductInput{
				Name: "Test Product",
				Slug: "Invalid Slug!",
			},
			wantErr: true,
			field:   "slug",
		},
		{
			name: "negative price",
			input: domain.CreateProductInput{
				Name:           "Test Product",
				Slug:           "test-product",
				BasePriceCents: -100,
			},
			wantErr: true,
			field:   "base_price_cents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.Validate()

			if tt.wantErr && len(errors) == 0 {
				t.Error("expected validation error, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("expected no errors, got %v", errors)
			}

			if tt.field != "" {
				found := false
				for _, err := range errors {
					if err.Field == tt.field {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error on field %q, got %v", tt.field, errors)
				}
			}
		})
	}
}
```

### 6.4 Authentication Logic Tests

Extract auth logic to domain, keep middleware thin:

```go
// internal/domain/auth.go
package domain

import (
	"time"
)

type TokenClaims struct {
	UserID    string
	TenantID  string
	Role      string
	ExpiresAt time.Time
}

type TokenValidationResult struct {
	Valid  bool
	Reason string
	Claims *TokenClaims
}

func ValidateTokenClaims(claims *TokenClaims, now time.Time) TokenValidationResult {
	if claims == nil {
		return TokenValidationResult{Valid: false, Reason: "missing claims"}
	}

	if claims.ExpiresAt.Before(now) {
		return TokenValidationResult{Valid: false, Reason: "token expired"}
	}

	if claims.UserID == "" {
		return TokenValidationResult{Valid: false, Reason: "missing user ID"}
	}

	if claims.TenantID == "" {
		return TokenValidationResult{Valid: false, Reason: "missing tenant ID"}
	}

	return TokenValidationResult{Valid: true, Claims: claims}
}

func CanAccessResource(userRole string, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"viewer": 1,
		"editor": 2,
		"admin":  3,
		"owner":  4,
	}

	userLevel := roleHierarchy[userRole]
	requiredLevel := roleHierarchy[requiredRole]

	return userLevel >= requiredLevel
}
```

Test the domain logic:

```go
// internal/domain/auth_test.go
package domain_test

import (
	"testing"
	"time"

	"myapp/internal/domain"
)

func TestValidateTokenClaims(t *testing.T) {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		claims    *domain.TokenClaims
		wantValid bool
		reason    string
	}{
		{
			name:      "nil claims",
			claims:    nil,
			wantValid: false,
			reason:    "missing claims",
		},
		{
			name: "valid claims",
			claims: &domain.TokenClaims{
				UserID:    "user-123",
				TenantID:  "tenant-456",
				Role:      "admin",
				ExpiresAt: now.Add(time.Hour),
			},
			wantValid: true,
		},
		{
			name: "expired token",
			claims: &domain.TokenClaims{
				UserID:    "user-123",
				TenantID:  "tenant-456",
				Role:      "admin",
				ExpiresAt: now.Add(-time.Hour),
			},
			wantValid: false,
			reason:    "token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.ValidateTokenClaims(tt.claims, now)

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if tt.reason != "" && result.Reason != tt.reason {
				t.Errorf("Reason = %q, want %q", result.Reason, tt.reason)
			}
		})
	}
}

func TestCanAccessResource(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		requiredRole string
		want         bool
	}{
		{"owner can access admin", "owner", "admin", true},
		{"admin can access admin", "admin", "admin", true},
		{"editor cannot access admin", "editor", "admin", false},
		{"admin can access viewer", "admin", "viewer", true},
		{"unknown role denied", "unknown", "viewer", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.CanAccessResource(tt.userRole, tt.requiredRole)
			if got != tt.want {
				t.Errorf("CanAccessResource(%q, %q) = %v, want %v",
					tt.userRole, tt.requiredRole, got, tt.want)
			}
		})
	}
}
```

#### What to Test vs Skip

| Layer | Test? | Why |
|-------|-------|-----|
| Token validation logic | ✅ Yes | Security-critical, pure functions |
| Role/permission checks | ✅ Yes | Business rules, easy to test |
| Middleware HTTP glue | ❌ No | Thin wrapper, trust stdlib |
| Token parsing (JWT lib) | ❌ No | Trust the library |

### 6.5 Business Logic Tests

```go
// internal/domain/money_test.go
package domain_test

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		name     string
		cents    int
		currency string
		want     string
	}{
		{"zero", 0, "USD", "$0.00"},
		{"positive", 1999, "USD", "$19.99"},
		{"large amount", 123456, "USD", "$1,234.56"},
		{"negative", -500, "USD", "-$5.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.FormatMoney(tt.cents, tt.currency)
			if got != tt.want {
				t.Errorf("FormatMoney(%d, %q) = %q, want %q", tt.cents, tt.currency, got, tt.want)
			}
		})
	}
}
```

### 6.6 Test Naming Conventions

| Pattern | Example |
|---------|---------|
| `Test<Function>` | `TestFormatMoney` |
| `Test<Type>_<Method>` | `TestContactForm_Validate` |
| `Test<Function>/<scenario>` | `TestFormatMoney/negative` |

Subtest names: lowercase with spaces, readable in output:

```go
t.Run("empty name returns error", func(t *testing.T) { ... })
t.Run("discount exceeds total", func(t *testing.T) { ... })
```

### 6.7 Running Tests

#### Makefile

```makefile
.PHONY: test
test:
	go test ./internal/domain/... -v

.PHONY: test-short
test-short:
	go test ./internal/domain/...

.PHONY: test-coverage
test-coverage:
	go test ./internal/domain/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
```

#### GitHub Actions

```yaml
name: Test and Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      
      - name: Run tests
        run: go test ./internal/domain/... -v

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      # ... deployment steps
```

### 6.8 When to Add More Tests

Expand beyond domain tests only when:

- A bug reaches production (add regression test)
- Complex handler logic that's hard to verify manually
- Critical payment/auth flows

Start minimal, add tests in response to actual problems.

---

## Section 7: Deployment & Infrastructure

### 7.1 Architecture

**Static Sites:**
```
Internet → Outer Caddy (HTTPS) → Container (Inner Caddy + Go API)
```

**Full-Stack Apps:**
```
Internet → Outer Caddy (HTTPS) → Container (Go App) → PostgreSQL
```

### 7.2 Docker Multi-Stage Build

#### Static Site

```dockerfile
# Stage 1: Build Hugo site
FROM hugomods/hugo:exts AS hugo-builder
WORKDIR /src
COPY package*.json ./
RUN npm install
COPY . .
RUN hugo --gc --minify

# Stage 2: Build Go API
FROM golang:1.21-alpine AS go-builder
WORKDIR /app
COPY api/go.mod ./
RUN go mod download
COPY api/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o contact-api .

# Stage 3: Final image with Caddy
FROM caddy:2-alpine

COPY --from=hugo-builder /src/public /srv
COPY --from=go-builder /app/contact-api /usr/local/bin/contact-api
COPY Caddyfile /etc/caddy/Caddyfile
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
```

#### Full-Stack App

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

### 7.3 Entrypoint Script

```bash
#!/bin/sh
set -e

# Start Go API in background (if exists)
if [ -f /usr/local/bin/contact-api ]; then
    /usr/local/bin/contact-api &
    sleep 1
fi

# Start Caddy in foreground
exec caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
```

### 7.4 Caddyfile (Inner Caddy)

```caddyfile
{
    admin off
    auto_https off
}

:{$PORT:80} {
    # Proxy API requests to Go backend
    handle /api/* {
        reverse_proxy localhost:8080
    }

    # Serve static files
    handle {
        root * /srv
        try_files {path} {path}/ /index.html
        file_server
    }

    header {
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        Referrer-Policy strict-origin-when-cross-origin
    }

    encode gzip zstd

    log {
        output stdout
        format console
    }
}
```

### 7.5 Docker Compose

#### Static Site

```yaml
services:
  web:
    image: ${DOCKER_IMAGE:-username/project:latest}
    ports:
      - "${LISTEN_PORT:-8082}:80"
    environment:
      - PORT=80
      - TURNSTILE_SECRET_KEY=${TURNSTILE_SECRET_KEY}
      - POSTMARK_TOKEN=${POSTMARK_TOKEN}
      - POSTMARK_TO=${POSTMARK_TO}
      - POSTMARK_FROM=${POSTMARK_FROM}
      - ALLOWED_ORIGINS=${ALLOWED_ORIGINS}
    restart: unless-stopped
```

#### Full-Stack App

```yaml
services:
  app:
    image: ${DOCKER_IMAGE:-username/project:latest}
    ports:
      - "${LISTEN_PORT:-8080}:8080"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - SESSION_SECRET=${SESSION_SECRET}
      - ENVIRONMENT=production
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres_data:
```

### 7.6 Environment Files

```bash
# .env.example (committed)
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
SESSION_SECRET=change-me-to-32-bytes
POSTMARK_TOKEN=your-token
ALLOWED_ORIGINS=https://example.com

# .env (on VPS only, never committed)
DATABASE_URL=postgres://prod_user:real_password@db:5432/prod_db?sslmode=disable
SESSION_SECRET=real-32-byte-secret-here
```

### 7.7 GitHub Actions CI/CD

```yaml
name: Test and Deploy

on:
  push:
    branches: [main]

env:
  IMAGE_NAME: project-name

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test ./internal/domain/... -v

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: docker/setup-buildx-action@v3

      - uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ${{ secrets.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:latest
            ${{ secrets.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Deploy to VPS
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            cd /opt/${{ env.IMAGE_NAME }}
            docker compose pull
            docker compose up -d
            docker image prune -f
```

### 7.8 VPS Setup

#### Outer Caddy Configuration

```caddyfile
# /etc/caddy/Caddyfile
www.example.com {
    reverse_proxy localhost:8082
}

example.com {
    redir https://www.example.com{uri} permanent
}
```

#### Deploy User Setup

```bash
# Create deploy user
sudo useradd -r -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# Setup SSH key
sudo mkdir -p /home/deploy/.ssh
sudo chmod 700 /home/deploy/.ssh
echo "ssh-ed25519 AAAA... deploy-key" | sudo tee /home/deploy/.ssh/authorized_keys
sudo chmod 600 /home/deploy/.ssh/authorized_keys
sudo chown -R deploy:deploy /home/deploy/.ssh

# Project directory
sudo mkdir -p /opt/project-name
sudo chown deploy:deploy /opt/project-name
```

#### Required Files on VPS

```
/opt/project-name/
├── docker-compose.yml
└── .env
```

### 7.9 Embedded Migrations

Run migrations on application startup:

```go
//go:embed migrations/*.sql
var migrations embed.FS

func main() {
    db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    
    // Run migrations before serving
    if err := runMigrations(db); err != nil {
        log.Fatalf("migrations failed: %v", err)
    }
    
    // Start server...
}

func runMigrations(db *sql.DB) error {
    goose.SetBaseFS(migrations)
    goose.SetDialect("postgres")
    return goose.Up(db, "migrations")
}
```

### 7.10 Database Backups

```bash
#!/bin/bash
# /usr/local/bin/backup-postgres.sh
set -e

BACKUP_DIR="/var/backups/postgresql"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=7
DB_NAME="mydb"

mkdir -p "$BACKUP_DIR"

# Create backup
pg_dump -Fc -h localhost -U postgres "$DB_NAME" > "$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.dump"

# Compress
zstd "$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.dump"
rm "$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.dump"

# Remove old backups
find "$BACKUP_DIR" -name "*.dump.zst" -mtime +$KEEP_DAYS -delete
```

**Cron schedule:**

```bash
# /etc/cron.d/postgres-backup
0 3 * * * postgres /usr/local/bin/backup-postgres.sh >> /var/log/pg-backup.log 2>&1
```

**Restoration:**

```bash
zstd -dc backup.dump.zst | pg_restore -h localhost -U postgres -d mydb --clean --if-exists
```

### 7.11 Rollback Procedures

**Manual rollback:**

```bash
cd /opt/project-name
docker compose down
DOCKER_IMAGE=username/project:previous-sha docker compose up -d
```

**Automatic rollback in CI:**

```yaml
- name: Deploy with rollback
  run: |
    # Tag current as previous
    docker tag $IMAGE:latest $IMAGE:previous || true
    
    # Deploy new version
    ssh deploy@$VPS "cd /opt/project && docker compose pull && docker compose up -d"
    
    # Health check
    for i in {1..12}; do
      if curl -sf https://example.com/api/health; then
        echo "Deployment successful"
        exit 0
      fi
      sleep 5
    done
    
    # Rollback on failure
    ssh deploy@$VPS "cd /opt/project && DOCKER_IMAGE=$IMAGE:previous docker compose up -d"
    exit 1
```

### 7.12 Health Endpoints

```go
// Basic health
mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
})

// Detailed health (for monitoring)
mux.HandleFunc("GET /api/health/ready", func(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "db": "down"})
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
})
```

### 7.13 Monitoring with Uptime Kuma

Simple, self-hosted monitoring:

```yaml
# docker-compose.monitoring.yml
services:
  uptime-kuma:
    image: louislam/uptime-kuma:1
    volumes:
      - uptime-kuma-data:/app/data
    ports:
      - "3001:3001"
    restart: unless-stopped

volumes:
  uptime-kuma-data:
```

Configure monitors for:
- HTTPS endpoint availability
- /api/health response
- SSL certificate expiry
- Response time thresholds

### 7.14 Deployment Checklist

**Pre-deployment:**
- [ ] All tests passing
- [ ] Environment variables configured
- [ ] DNS records set
- [ ] Images optimized (static sites)

**Post-deployment:**
- [ ] HTTPS loads correctly
- [ ] Forms submit successfully
- [ ] Emails received
- [ ] No console errors
- [ ] SSL certificate valid

**Troubleshooting:**
- Container won't start → Check logs: `docker compose logs`
- 502 Bad Gateway → Inner service not running, check ports
- CORS errors → Check ALLOWED_ORIGINS matches domain exactly
- Turnstile fails → Use test keys for development

---

## Section 8: Project Documentation

### 8.1 Purpose of CLAUDE.md

CLAUDE.md is the entry point for AI assistants working with your codebase. It provides:

- Quick orientation to the project
- Commands to run without searching
- Architecture decisions that aren't obvious from code
- Patterns to follow for consistency
- Gotchas and non-obvious behaviors

**Design principle:** Write what you'd tell a new developer on day one.

### 8.2 Template Structure

```markdown
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

[One paragraph: what this project is, who it's for, core technology choices]

## Commands

```bash
# Development
[command]    # [what it does]

# Build
[command]    # [what it does]

# Database (if applicable)
[command]    # [what it does]

# Testing
[command]    # [what it does]
```

## Architecture

[Brief explanation of how the system is structured]

## Project Structure

```
[directory tree with annotations]
```

## Key Patterns

[Important conventions that aren't obvious from code]

## Environment Variables

[Required variables and what they configure]

## Deployment

[How the project gets to production]
```

### 8.3 Full-Stack Application Example

```markdown
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Multi-tenant e-commerce platform for custom apparel businesses. Built with Go, PostgreSQL, htmx, and Alpine.js. Each tenant gets isolated data and customizable storefronts.

## Commands

```bash
# Development
make dev              # Start server with hot reload (air)
make dev-css          # Watch and rebuild Tailwind CSS

# Database
make db-up            # Start PostgreSQL container
make db-migrate       # Run pending migrations
make db-new NAME=x    # Create new migration
make sqlc             # Regenerate query code

# Testing
make test             # Run domain tests

# Build
docker compose up -d  # Run production stack
```

## Architecture

```
Internet → Caddy (HTTPS) → Go App → PostgreSQL
```

**Multi-tenancy:** Tenant resolved from subdomain via middleware. All queries include tenant_id filter.

**Rendering:** Server-side HTML with htmx for dynamic updates. Alpine.js for client-side interactivity.

## Key Patterns

### Handler Structure

Handlers grouped by user context (admin/, storefront/), not by entity.

### Multi-Tenant Queries

Every query includes tenant_id:

```sql
SELECT * FROM products WHERE tenant_id = $1 AND id = $2;
```

### htmx Responses

Check HX-Request header for partial vs full page:

```go
if r.Header.Get("HX-Request") == "true" {
    h.renderer.RenderPartial(w, "products/_row", product)
} else {
    h.renderer.Render(w, "products/show", data)
}
```

## Environment Variables

```bash
DATABASE_URL=postgres://...
SESSION_SECRET=32-byte-random-string
```

## Deployment

Push to main → GitHub Actions builds Docker image → SSH deploy to VPS → Migrations run on startup.
```

### 8.4 Static Site Example

```markdown
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Brochure site for a consulting company. Built with Hugo static site generator and a Go API backend for contact form handling.

## Commands

```bash
# Development
npm install             # Install dependencies (first time)
hugo server             # Start dev server on localhost:1313

# Build
hugo --gc --minify      # Build for production

# Docker
docker compose up -d    # Run production container
```

## Architecture

```
Internet → Caddy (HTTPS) → Container
                              ├── Inner Caddy (static files)
                              └── Go API (/api/*)
```

## Key Patterns

### Data-Driven Content

Content defined in `data/*.yaml`, accessed via `site.Data.services.packages`.

### Contact Form

Alpine.js form → POST /api/contact → Validation + Turnstile → Postmark email.

### Spam Protection

Honeypot field + Cloudflare Turnstile.

## Environment Variables

```bash
POSTMARK_TOKEN=xxx
TURNSTILE_SECRET_KEY=xxx
ALLOWED_ORIGINS=https://example.com
```

## Deployment

GitHub Actions builds Docker image on push to main, deploys via SSH to VPS.
```

### 8.5 Writing Guidelines

**Be specific about commands:**

```markdown
# ❌ Vague
Run the development server to start working.

# ✅ Specific
```bash
make dev    # Start server with hot reload on localhost:8080
```
```

**Document non-obvious patterns:**

```markdown
# ❌ Obvious (skip it)
Use `go build` to build the project.

# ✅ Non-obvious (document it)
Handlers are organized by user context (admin/, storefront/), not by entity.
```

**Include the "why":**

```markdown
# ❌ What only
All queries include tenant_id.

# ✅ What + Why
All queries include tenant_id filter. This enforces data isolation in multi-tenant architecture.
```

### 8.6 What to Include vs Exclude

**Include:**
- Project purpose (one paragraph)
- Commands to run/build/test
- Architecture diagram (ASCII)
- Directory structure with annotations
- Non-obvious patterns and conventions
- Required environment variables
- Deployment process summary

**Exclude:**
- API documentation (use OpenAPI)
- Detailed database schema (migrations are source of truth)
- Step-by-step tutorials
- Duplicate of README content

### 8.7 Keeping CLAUDE.md Updated

Update when:
- Adding new commands
- Changing project structure
- Introducing new patterns
- Modifying environment requirements
- Changing deployment process

**Tip:** Add a PR checklist item: "Update CLAUDE.md if needed"

---

## Quick Reference

### File Locations

| Type | Full-Stack | Static Site |
|------|-----------|-------------|
| Entry point | `cmd/server/main.go` | `hugo.toml` |
| Handlers | `internal/handler/` | `api/handlers.go` |
| Templates | `internal/templates/` | `layouts/` |
| Static files | `static/` | `static/` |
| Database | `migrations/` | N/A |
| Content | N/A | `content/`, `data/` |

### Common Commands

```bash
# Development
make dev              # Start with hot reload
hugo server           # Hugo dev server

# Database
make db-migrate       # Run migrations
make sqlc             # Generate queries

# Testing
make test             # Run tests

# Deployment
docker compose up -d  # Start containers
docker compose logs   # View logs
```

### Port Allocation

| Port | Use |
|------|-----|
| 80/443 | Outer Caddy (host) |
| 8080 | Go app (container) |
| 80 | Inner Caddy (container) |
| 5432 | PostgreSQL |

---

*Last updated: December 2024*