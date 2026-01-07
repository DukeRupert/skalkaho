# Skalkaho - MVP Technical Guide

> Construction quoting tool for small-medium contractors

---

## Project Overview

Skalkaho is a construction quoting SaaS application designed to help small-medium contractors build professional quotes with real-time pricing calculations. The initial client is a Hamilton, MT-based contractor currently managing quotes manually.

### Product Name Origin

**Skalkaho** - Named after Skalkaho Pass, a scenic mountain route east of Hamilton, Montana in the Bitterroot Valley. The name fits regional naming conventions (similar to Lolo, Kootenai, Flathead) and provides authentic local identity.

**Domain:** skalkaho.com

---

## MVP Scope

### Included in MVP

- Build quotes with hierarchical categories
- Add line items (materials, labor, equipment, and custom types)
- Custom item types per job (e.g., "Subcontractor", "Permits")
- Apply surcharges at job, type, category, and line-item levels
- Per-type surcharge rates (e.g., 20% materials, 10% labor)
- Real-time total calculations with type breakdowns
- Configurable surcharge modes (stacking vs. override)
- Order List and Site Materials reports
- Excel/CSV price import with AI-powered matching
- Line item tags for visual grouping within categories

### Excluded from MVP (Future Phases)

| Feature | Phase |
|---------|-------|
| User accounts / authentication | 2 |
| Saving/loading quotes | 2 |
| PDF export | 2 |
| ~~Material list generation~~ | ~~3~~ Done in MVP |
| Supplier catalog imports | 3 |
| Labor forecasting | 3 |
| QuickBooks integration | 4 |
| Multi-user / team features | 4 |

---

## Data Model

### Entity Relationship

```
Settings (singleton)
    │
    └── provides defaults for ──▶ Job ◀── belongs to ── Client
                                   │
                                   ├── has many ──▶ JobItemType (custom types)
                                   │
                                   └── has many ──▶ Category
                                                      │
                                                      ├── has many ──▶ Category (nested, max 3 levels)
                                                      │
                                                      └── has many ──▶ LineItem
```

---

### Client

Represents a customer who receives quotes.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `name` | string | required, unique | Primary display name |
| `company` | string | nullable | Business/company name |
| `email` | string | nullable | Contact email |
| `phone` | string | nullable | Contact phone |
| `address` | string | nullable | Street address |
| `city` | string | nullable | City |
| `state` | string | nullable | State/Province |
| `zip` | string | nullable | Postal code |
| `tax_id` | string | nullable | Tax ID for invoicing |
| `notes` | string | nullable | Internal notes |
| `created_at` | timestamp | auto | Creation timestamp |

---

### Job

The top-level container for a quote.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `name` | string | required | e.g., "Smith Kitchen Remodel" |
| `customer_name` | string | nullable | Legacy field, prefer client_id |
| `surcharge_percent` | decimal | default: 0 | Default surcharge for unset types |
| `material_surcharge_percent` | decimal | nullable | Material-specific surcharge |
| `labor_surcharge_percent` | decimal | nullable | Labor-specific surcharge |
| `equipment_surcharge_percent` | decimal | nullable | Equipment-specific surcharge |
| `surcharge_mode` | enum | "stacking" \| "override" | How surcharges combine |
| `status` | enum | "draft" \| "sent" \| "accepted" \| "rejected" | Quote lifecycle status |
| `expires_at` | timestamp | nullable | Quote expiration date |
| `client_id` | UUID | FK → Client, nullable | Associated client |
| `created_at` | timestamp | auto | Creation timestamp |

**Notes:**
- `surcharge_mode` defaults to value from Settings when creating new jobs
- `client_id` links to client record; `customer_name` kept for backwards compatibility
- Client can only be changed when status is "draft"
- Type-specific surcharges (material, labor, equipment) override the default when set
- Custom item types have their own surcharge in the JobItemType table

---

### JobItemType

Custom line item types defined per job (beyond the standard material/labor/equipment).

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `job_id` | UUID | FK → Job | Parent job |
| `name` | string | required | Display name (e.g., "Subcontractor") |
| `slug` | string | required, unique per job | URL-safe identifier (e.g., "subcontractor") |
| `color` | string | required | Tailwind color prefix for UI (e.g., "amber") |
| `sort_order` | int | default: 0 | Display order |
| `surcharge_percent` | decimal | nullable | Type-specific surcharge |
| `created_at` | timestamp | auto | Creation timestamp |

**Notes:**
- Custom types appear alongside standard types (material, labor, equipment)
- If `surcharge_percent` is null, falls back to job's default `surcharge_percent`
- Slug must be unique within a job, lowercase alphanumeric with hyphens

---

### Category

Organizational groupings within a job. Supports nesting up to 3 levels deep.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `job_id` | UUID | FK → Job | Parent job |
| `parent_id` | UUID | FK → Category, nullable | Null = top-level category |
| `name` | string | required | e.g., "Framing", "Electrical", "Rough-In" |
| `surcharge_percent` | decimal | nullable | Null = inherit from parent/job |
| `sort_order` | int | default: 0 | Manual ordering within siblings |

**Nesting Examples:**
```
Electrical (top-level)
├── Rough-In (level 2)
│   ├── Main Panel (level 3)
│   └── Circuits (level 3)
└── Finish (level 2)
    ├── Outlets (level 3)
    └── Fixtures (level 3)
```

**Constraints:**
- Maximum nesting depth: 3 levels for MVP
- Enforce via application logic, not database constraint

---

### LineItem

Individual materials or labor entries within a category.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `category_id` | UUID | FK → Category | Parent category |
| `type` | string | required | "material", "labor", "equipment", or custom slug |
| `name` | string | required | e.g., "2x4 Lumber", "Electrician" |
| `description` | string | nullable | e.g., "8ft pressure treated" |
| `quantity` | decimal | required | Supports partial units (e.g., 2.5) |
| `unit` | string | required | Free-form with UI suggestions |
| `unit_price` | decimal | required | Price per unit |
| `surcharge_percent` | decimal | nullable | Null = inherit from type/category |
| `tag` | string | nullable | Visual grouping label within category |
| `sort_order` | int | default: 0 | Manual ordering within category |

**Standard Types:**
- `material` - Physical goods
- `labor` - Work hours/services
- `equipment` - Tools and machinery

**Custom Types:**
- Any slug from job's JobItemType records (e.g., "subcontractor", "permits")

**Tags:**
- Optional label for visual grouping within a category
- Items with the same tag are displayed together with indentation
- Scoped within item type (material tags separate from labor tags)

**Common Units (UI suggestions):**
- Materials: `ea`, `sqft`, `lnft`, `bundle`, `box`, `bag`, `gal`, `sheet`
- Labor: `hr`, `day`, `job`, `sqft`

**Calculated Fields (not stored):**
- `base_price` = quantity × unit_price
- `effective_surcharge` = resolved surcharge based on inheritance
- `final_price` = base_price × (1 + effective_surcharge / 100)

---

### Settings

Application-wide defaults. Single row table.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Single row identifier |
| `default_surcharge_mode` | enum | "stacking" \| "override" | Default for new jobs |
| `default_surcharge_percent` | decimal | default: 0 | Default surcharge for new jobs |

---

## Surcharge Calculation Logic

### Type-Based Surcharges

Surcharges are now applied per item type, allowing different markup rates for materials, labor, and equipment:

```
Job
├── material_surcharge_percent: 20%
├── labor_surcharge_percent: 10%
├── equipment_surcharge_percent: 15%
└── surcharge_percent: 12% (default for unset types and custom types)
```

**Resolution Order for Type Surcharge:**
1. Standard types (material/labor/equipment) use their specific job field if set
2. Custom types use their `surcharge_percent` from JobItemType if set
3. Falls back to job's default `surcharge_percent`

### Inheritance Rules

Surcharges cascade down the hierarchy with explicit values overriding inherited ones:

1. **LineItem** uses its own `surcharge_percent` if set
2. Otherwise, inherits from parent **Category** chain
3. Otherwise, uses **Type-specific surcharge** from Job or JobItemType
4. Falls back to Job's default `surcharge_percent`

```
Job (default: 15%, material: 20%, labor: 10%)
└── Category A (null → inherits type surcharge)
    ├── Material Item 1 (null → uses 20%)
    ├── Labor Item 2 (null → uses 10%)
    └── Material Item 3 (5% → uses 5%)
└── Category B (8% → overrides type surcharge)
    └── Subcategory B1 (null → inherits 8%)
        └── Material Item 4 (null → uses 8% in override mode)
```

### Surcharge Modes

#### Stacking Mode (Additive)

All applicable surcharges add together:

```
Effective Surcharge = TypeSurcharge% + Category% + LineItem%
```

**Example:**
- Material type surcharge: 20%
- Category surcharge: 5%
- LineItem surcharge: 3%
- **Total: 28%**

A $100 material item becomes $128.

#### Override Mode

Only the most specific (lowest-level) surcharge applies:

```
Effective Surcharge = LineItem% ?? Category% ?? TypeSurcharge%
```

**Example:**
- Material type surcharge: 20%
- Category surcharge: 10%
- LineItem surcharge: null
- **Total: 10%** (Category value used)

A $100 item becomes $110.

### Calculation Pseudocode

```go
// GetTypeSurcharge returns the surcharge for a specific item type
func GetTypeSurcharge(job *Job, itemType LineItemType, customTypes []*JobItemType) decimal {
    switch itemType {
    case "material":
        if job.MaterialSurchargePercent != nil {
            return *job.MaterialSurchargePercent
        }
    case "labor":
        if job.LaborSurchargePercent != nil {
            return *job.LaborSurchargePercent
        }
    case "equipment":
        if job.EquipmentSurchargePercent != nil {
            return *job.EquipmentSurchargePercent
        }
    default:
        // Custom type - look up in customTypes
        for _, ct := range customTypes {
            if ct.Slug == itemType && ct.SurchargePercent != nil {
                return *ct.SurchargePercent
            }
        }
    }
    return job.SurchargePercent // Fall back to default
}

func (li *LineItem) EffectiveSurcharge(job *Job, categoryChain []*Category, customTypes []*JobItemType) decimal {
    if job.SurchargeMode == "override" {
        // Most specific non-null value wins
        if li.SurchargePercent != nil {
            return *li.SurchargePercent
        }
        // Walk category chain from deepest to shallowest
        for i := len(categoryChain) - 1; i >= 0; i-- {
            if categoryChain[i].SurchargePercent != nil {
                return *categoryChain[i].SurchargePercent
            }
        }
        return GetTypeSurcharge(job, li.Type, customTypes)
    }

    // Stacking mode - sum all levels
    total := GetTypeSurcharge(job, li.Type, customTypes)
    for _, cat := range categoryChain {
        if cat.SurchargePercent != nil {
            total += *cat.SurchargePercent
        }
    }
    if li.SurchargePercent != nil {
        total += *li.SurchargePercent
    }
    return total
}

func (li *LineItem) FinalPrice(job *Job, categoryChain []*Category, customTypes []*JobItemType) decimal {
    base := li.Quantity * li.UnitPrice
    surcharge := li.EffectiveSurcharge(job, categoryChain, customTypes)
    return base * (1 + surcharge / 100)
}
```

---

## Quote Totals Calculation

### Category Total

```
CategoryTotal = Σ(LineItem.FinalPrice) + Σ(ChildCategory.Total)
```

### Job Total

```
JobTotal = Σ(TopLevelCategory.Total)
```

### Summary Display

The UI should display:
- **Subtotal**: Sum of all base prices (quantity × unit_price)
- **Total Surcharges**: JobTotal - Subtotal
- **Grand Total**: JobTotal

**Type Breakdown** (displayed on job and category pages):
- **Materials Total**: Sum of final prices for type = "material"
- **Labor Total**: Sum of final prices for type = "labor"
- **Equipment Total**: Sum of final prices for type = "equipment"
- **Custom Type Totals**: Sum of final prices for each custom type

Line items are grouped by type in category view, with each type displayed in its own card/section.

---

## Technical Stack (Recommended)

| Layer | Technology |
|-------|------------|
| Backend | Go |
| Frontend | HTMX + Alpine.js + Tailwind CSS |
| Database | SQLite (MVP), PostgreSQL (production) |
| Templating | Go html/template |

### Why This Stack

- **Go**: Fast, simple deployment, great stdlib
- **HTMX**: Real-time updates without SPA complexity
- **Alpine.js**: Lightweight interactivity for UI components
- **Tailwind**: Rapid styling without custom CSS
- **SQLite**: Zero-config for MVP, easy local development

---

## Database Schema (SQLite)

```sql
-- Settings (singleton)
CREATE TABLE settings (
    id TEXT PRIMARY KEY DEFAULT 'default',
    default_surcharge_mode TEXT NOT NULL DEFAULT 'stacking' 
        CHECK (default_surcharge_mode IN ('stacking', 'override')),
    default_surcharge_percent REAL NOT NULL DEFAULT 0
);

INSERT INTO settings (id) VALUES ('default');

-- Jobs
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    customer_name TEXT,
    surcharge_percent REAL NOT NULL DEFAULT 0,
    surcharge_mode TEXT NOT NULL DEFAULT 'stacking'
        CHECK (surcharge_mode IN ('stacking', 'override')),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Categories
CREATE TABLE categories (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    parent_id TEXT REFERENCES categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_categories_job ON categories(job_id);
CREATE INDEX idx_categories_parent ON categories(parent_id);

-- Line Items
CREATE TABLE line_items (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('material', 'labor')),
    name TEXT NOT NULL,
    description TEXT,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    unit_price REAL NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_line_items_category ON line_items(category_id);
```

---

## Go Structs

```go
package models

import (
    "time"
)

type SurchargeMode string

const (
    SurchargeModeStacking SurchargeMode = "stacking"
    SurchargeModeOverride SurchargeMode = "override"
)

type LineItemType string

const (
    LineItemTypeMaterial LineItemType = "material"
    LineItemTypeLabor    LineItemType = "labor"
)

type Settings struct {
    ID                      string        `db:"id"`
    DefaultSurchargeMode    SurchargeMode `db:"default_surcharge_mode"`
    DefaultSurchargePercent float64       `db:"default_surcharge_percent"`
}

type Job struct {
    ID               string        `db:"id"`
    Name             string        `db:"name"`
    CustomerName     *string       `db:"customer_name"`
    SurchargePercent float64       `db:"surcharge_percent"`
    SurchargeMode    SurchargeMode `db:"surcharge_mode"`
    CreatedAt        time.Time     `db:"created_at"`
}

type Category struct {
    ID               string   `db:"id"`
    JobID            string   `db:"job_id"`
    ParentID         *string  `db:"parent_id"`
    Name             string   `db:"name"`
    SurchargePercent *float64 `db:"surcharge_percent"`
    SortOrder        int      `db:"sort_order"`
}

type LineItem struct {
    ID               string       `db:"id"`
    CategoryID       string       `db:"category_id"`
    Type             LineItemType `db:"type"`
    Name             string       `db:"name"`
    Description      *string      `db:"description"`
    Quantity         float64      `db:"quantity"`
    Unit             string       `db:"unit"`
    UnitPrice        float64      `db:"unit_price"`
    SurchargePercent *float64     `db:"surcharge_percent"`
    SortOrder        int          `db:"sort_order"`
}
```

---

## UI Behavior Notes

### Real-Time Updates

- Totals recalculate immediately on any change
- Use HTMX `hx-trigger="change"` on inputs
- Debounce quantity/price inputs (300ms)

### Surcharge Display

- Show effective surcharge on each line item
- In stacking mode, consider tooltip showing breakdown:
  `"Job: 15% + Category: 10% + Item: 5% = 30%"`

### Category Nesting

- Indent nested categories visually
- Collapse/expand for deep hierarchies
- Drag-and-drop reordering (future enhancement)

### Line Item Type

- Visual indicator (icon or color) for material vs. labor
- Filter/group by type in summary views

---

## Open Questions for Future

1. **Rounding**: Round final prices to nearest cent? Per-item or on totals only?
2. **Negative surcharges**: Allow discounts as negative percentages?
3. **Tax handling**: Separate from surcharges? MVP excludes tax.
4. **Templates**: Save category structures as reusable templates?
5. **Versioning**: Track quote revisions?

---

## Revision History

| Date | Version | Notes |
|------|---------|-------|
| 2024-12-22 | 0.1 | Initial MVP specification |
| 2025-12-27 | 0.2 | MVP complete - added equipment type, reports |
| 2025-12-28 | 0.3 | Added client management, quote status, item templates |
| 2026-01-07 | 0.4 | Added custom item types, per-type markup, line item tags, Excel import |