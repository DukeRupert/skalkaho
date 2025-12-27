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
- Add line items (materials, labor, and equipment)
- Apply surcharges at job, category, and line-item levels
- Real-time total calculations
- Configurable surcharge modes (stacking vs. override)
- Order List and Site Materials reports

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
    └── provides defaults for ──▶ Job
                                   │
                                   └── has many ──▶ Category
                                                      │
                                                      ├── has many ──▶ Category (nested, max 3 levels)
                                                      │
                                                      └── has many ──▶ LineItem
```

---

### Job

The top-level container for a quote.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `id` | UUID | PK | Unique identifier |
| `name` | string | required | e.g., "Smith Kitchen Remodel" |
| `customer_name` | string | nullable | Reference only, not for invoicing |
| `surcharge_percent` | decimal | default: 0 | Job-level surcharge (e.g., 15.0 for 15%) |
| `surcharge_mode` | enum | "stacking" \| "override" | How surcharges combine |
| `created_at` | timestamp | auto | Creation timestamp |

**Notes:**
- `surcharge_mode` defaults to value from Settings when creating new jobs
- `customer_name` is informal reference; formal customer management is future scope

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
| `type` | enum | "material" \| "labor" \| "equipment" | Critical for material list generation |
| `name` | string | required | e.g., "2x4 Lumber", "Electrician" |
| `description` | string | nullable | e.g., "8ft pressure treated" |
| `quantity` | decimal | required | Supports partial units (e.g., 2.5) |
| `unit` | string | required | Free-form with UI suggestions |
| `unit_price` | decimal | required | Price per unit |
| `surcharge_percent` | decimal | nullable | Null = inherit from category |
| `sort_order` | int | default: 0 | Manual ordering within category |

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

### Inheritance Rules

Surcharges cascade down the hierarchy with explicit values overriding inherited ones:

1. **LineItem** uses its own `surcharge_percent` if set
2. Otherwise, inherits from parent **Category**
3. Category inherits from its parent Category (if nested)
4. Top-level Category inherits from **Job**

```
Job (15%)
└── Category A (null → inherits 15%)
    ├── LineItem 1 (null → inherits 15%)
    └── LineItem 2 (5% → uses 5%)
└── Category B (10% → uses 10%)
    └── Subcategory B1 (null → inherits 10%)
        └── LineItem 3 (null → inherits 10%)
```

### Surcharge Modes

#### Stacking Mode (Additive)

All applicable surcharges add together:

```
Effective Surcharge = Job% + Category% + LineItem%
```

**Example:**
- Job surcharge: 15%
- Category surcharge: 10%
- LineItem surcharge: 5%
- **Total: 30%**

A $100 item becomes $130.

#### Override Mode

Only the most specific (lowest-level) surcharge applies:

```
Effective Surcharge = LineItem% ?? Category% ?? Job%
```

**Example:**
- Job surcharge: 15%
- Category surcharge: 10%
- LineItem surcharge: 5%
- **Total: 5%** (LineItem value used)

A $100 item becomes $105.

### Calculation Pseudocode

```go
func (li *LineItem) EffectiveSurcharge(job *Job, category *Category) decimal {
    if job.SurchargeMode == "override" {
        // Most specific non-null value wins
        if li.SurchargePercent != nil {
            return *li.SurchargePercent
        }
        return category.ResolveSurcharge(job) // walks up tree
    }
    
    // Stacking mode - sum all levels
    total := job.SurchargePercent
    total += category.StackedSurcharge() // sums category + ancestors
    if li.SurchargePercent != nil {
        total += *li.SurchargePercent
    }
    return total
}

func (li *LineItem) FinalPrice(job *Job, category *Category) decimal {
    base := li.Quantity * li.UnitPrice
    surcharge := li.EffectiveSurcharge(job, category)
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

Optionally break down by type:
- **Materials Subtotal**: Sum of LineItems where type = "material"
- **Labor Subtotal**: Sum of LineItems where type = "labor"

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