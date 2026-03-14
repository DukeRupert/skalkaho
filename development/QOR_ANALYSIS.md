# QOR Framework Analysis

> Analyzed 2026-03-03 from https://github.com/qor/qor
> Purpose: Identify patterns and edge cases relevant to Skalkaho/Hiri development

## What is QOR?

QOR is a collection of ~40 Go libraries (5,363 stars, MIT licensed) for building business/admin applications. Created by The Plant (Japan), it wraps GORM models as "Resources" and auto-generates admin UIs, APIs, and CRUD flows. There is a newer QOR5 rewrite that prefers Go-rendered HTML over templates.

## Patterns Worth Adopting

### 1. State Machine (`transition` package)

Clean lifecycle management with entry/exit hooks and audit logging:

```go
OrderState := transition.New(&Order{})
OrderState.Initial("draft")
OrderState.State("pending")
OrderState.State("processing").Enter(func(order interface{}, tx *gorm.DB) error {
    // Business logic when entering this state
    return nil
})
OrderState.Event("checkout").To("pending").From("draft")
    .Before(func(...) error { /* validate */ })
    .After(func(...) error { /* side effects */ })
```

Key features:
- States have Enter/Exit hooks
- Events have Before/After hooks
- Events can have multiple From states mapping to different To states
- Automatic state change logging in DB
- Transitions run within a DB transaction

**Relevance**: Quote/job lifecycle: `draft → sent → reviewed → approved/rejected → invoiced`

### 2. Processor Pipeline

When saving a resource, QOR runs: Initialize → Validate → Decode → Process → Save. Any step can short-circuit with `ErrProcessorSkipLeft`. Useful for recalculating totals when a line item changes.

### 3. Error Accumulation (not short-circuit)

`qor.Errors` collects *all* validation errors rather than returning on the first. Forms show every issue at once instead of one-at-a-time.

### 4. Composable Model Behaviors via Embedding

Add capabilities by embedding behavior structs:

```go
type Product struct {
    gorm.Model
    audited.AuditedModel  // Tracks CreatedBy/UpdatedBy
    transition.Transition // Adds state machine
    sorting.SortingDESC   // Adds position-based ordering
}
```

Each embedded type adds DB columns, registers callbacks, and auto-configures admin UI.

### 5. Role-Based Field Permissions

Permissions at both resource-level AND field-level:

```go
product.Meta(&admin.Meta{
    Name: "Price",
    Permission: roles.Allow(roles.Read, "admin"),
})
```

The `roles` package uses deny-first, then allow: Deny map checked first → Allow map empty means anyone allowed → Check Allow map.

### 6. Scoped Queries

Named, predefined query filters ("Active Only", "This Month") reusable across handlers.

### 7. Resource-Oriented Architecture

Every model becomes a Resource with standardized CRUD handlers, validators, and processors. A single `Admin.AddResource(&Product{})` call generates index, show, new, edit, delete pages plus RESTful JSON API.

### 8. Draft/Publish System

Shadow table pattern: `products` (live) + `products_draft` (admin edits). "Publish" copies draft to production, "Discard" reverts.

### 9. Exchange (Import/Export)

Clean abstraction for data import/export with field mapping, validation, processing hooks, and progress callbacks.

## Edge Cases & Gotchas

1. **No built-in multi-tenancy** — despite being a business framework. Tenant scoping is DIY. Our `org_id` pattern is the right approach.
2. **CSRF protection is weak** — only checks Referer header host match, no token-based CSRF. Always use proper CSRF tokens.
3. **Create vs Update transaction inconsistency** — `Update` wraps in a DB transaction but `Create` does not. Always use transactions for both.
4. **Reflection-heavy code is fragile** — QOR's auto-mapping via reflection is its weakest part. Code generation (sqlc) is cleaner and more debuggable.
5. **Composite primary keys are painful** — QOR uses a `^|^` separator hack. Use surrogate keys instead.
6. **Template discovery is brittle** — hardcoded search paths from pre-modules era. Design template loading explicitly.

## What NOT to Copy

- **GORM v1 dependency** — outdated
- **Reflection-based field mapping** — sqlc gives type safety for free
- **Full admin auto-generation** — for specialized UIs (construction quoting), hand-built HTMX gives more control. Auto-gen better for back-office screens only.
- **QOR's Context threading** — uses a custom context struct instead of Go's `context.Context`. Prefer standard library patterns.

## Key Takeaway

The most transferable patterns are the **state machine**, **processor pipeline**, **error accumulation**, and **field-level permissions**. QOR proves that a thin core + composable plugins is more maintainable than a monolith, but its reflection-heavy internals show why code generation (sqlc) is the better path.
