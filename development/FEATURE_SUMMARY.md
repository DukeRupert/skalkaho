# Skalkaho Feature Summary

> A construction quoting tool built for how contractors actually work

---

## What Skalkaho Does

Skalkaho replaces spreadsheets and manual quote-building for small-to-medium contractors. A contractor opens the app, builds a quote by organizing work into categories and line items, applies markup, assigns a client, and sends a professional estimate for digital signature. The whole flow — from blank quote to signed contract — happens inside a single keyboard-driven interface.

The application was designed around one real contractor's workflow in Hamilton, Montana. Every feature exists because the user needed it.

---

## The Quoting Workflow

This is the core of Skalkaho. Everything else supports this flow.

### 1. Start a Quote

From the dashboard, the user hits **`n`** or taps "New Quote." A job is created instantly with the default name "New Quote" and inherits the organization's default markup settings. No forms, no modals — the user lands directly on the quote's spreadsheet view ready to work.

The dashboard shows all quotes with status badges (Draft, Sent, Accepted, Rejected, Expired), client names, totals, and estimate/signature indicators. Quotes can be filtered by status and sorted by date or name.

### 2. Organize Work into Categories

Quotes are structured as a hierarchy:

```
Job (Quote)
└── Category (e.g., "Electrical")
    └── Subcategory (e.g., "Rough-In")
        └── Sub-subcategory (e.g., "Main Panel")
            └── Line Items
```

Categories nest up to 3 levels deep. This matches how contractors think about a job — major phases broken into specific work areas. The user presses **`c`** to add a category, and subcategories can be added from any category's action menu.

On the job page, categories display as expandable/collapsible rows in the spreadsheet view. Press **`Space`** to toggle a category open or closed. The category tree also appears as a navigation sidebar on desktop for quick jumping between sections.

### 3. Add Line Items

This is where the user spends most of their time. Every item has a type, name, quantity, unit, and unit price.

**Three standard types**, each color-coded throughout the UI:
- **Material** (forest green) — physical goods: lumber, wire, fixtures
- **Labor** (copper) — work hours: electrician, plumber, general labor
- **Equipment** (slate) — tools and machinery: excavator rental, scaffolding

**Custom types** can be created per job for anything that doesn't fit — "Subcontractor", "Permits", "Disposal". Each gets a name, color (chosen from a curated palette), and optional type-specific markup.

**Adding items happens multiple ways, depending on speed needed:**

**Keyboard shortcuts** — The fastest path. From a category view, press **`m`** for material, **`l`** for labor, **`e`** for equipment. An inline form appears with the type pre-selected and sensible default units (ea for materials, hr for labor, day for equipment).

**Template autocomplete** — As the user types an item name, the form searches the item template library and suggests matches. Selecting a template pre-fills the name, unit, and price. This is how the user avoids retyping "12/2 Romex NM-B Wire" and "$0.49/lnft" on every quote.

**Batch add** — Press **`b`** to open batch mode. This pulls all templates from a named category (e.g., all items tagged "Electrical > Rough-In") and creates them as line items in one action. Useful for repeating standard scope across quotes.

**Spreadsheet inline editing** — On the job page's spreadsheet view, press **`a`** or **`+`** on a selected category row to add an item directly inline, without leaving the page. Items can also be edited inline from the spreadsheet by pressing **`Enter`** on a selected item row.

### 4. Tags for Visual Grouping

Within a category, line items can be tagged with a free-form label. Items sharing the same tag are grouped together visually under a colored sub-header. Tags are scoped per item type — material tags are separate from labor tags.

This handles the case where a category like "Rough-In" has multiple clusters of related items (e.g., "Main Panel" materials vs. "Branch Circuits" materials) without needing another level of nesting.

### 5. Apply Markup

Markup (surcharge) is the profit margin applied on top of base costs. Skalkaho's markup system is built for flexibility because different contractors mark up different types of work differently.

**Markup can be set at four levels:**

| Level | Example | How to set |
|-------|---------|-----------|
| **Job default** | 15% on everything | Press **`%`** on job page |
| **Per type** | 20% materials, 10% labor | Press **`%`**, set per-type rates |
| **Per category** | 25% on "Custom Millwork" | Press **`%`** on category page |
| **Per line item** | 5% on a specific fixture | Edit the item |

**Two modes control how these levels interact:**

- **Override mode** (default): The most specific markup wins. If an item has its own 5%, that's used regardless of category or type rates. Simple and predictable.
- **Stacking mode**: All applicable markups add together. Job 15% + Category 10% + Item 5% = 30% total. Useful when markup represents separate concerns (overhead, profit, contingency).

The mode is set at the organization level in Settings and inherited by new jobs.

When a category or item is using inherited markup rather than its own, the UI displays "inherit" to make the source clear. The user always knows where a number is coming from.

### 6. Review the Numbers

The job page shows a running total breakdown:

- **Per-type subtotals** — Materials: $X, Labor: $Y, Equipment: $Z
- **Grand total** — All types combined, with markup applied

Two additional reports are available:

**Order List** (`o` key) — Aggregates all materials and equipment across the entire quote by name and unit, sorted alphabetically. This is the list the contractor hands to a supplier: "I need 200 lnft of 12/2 Romex, 50 ea outlet boxes, 30 ea switches." Duplicate items from different categories are combined.

**Site Materials** (`s` key) — Shows materials and equipment broken down by category with full path labels (e.g., "Electrical > Rough-In > Main Panel"). This is the list for on-site logistics: what materials go where on the jobsite.

Both reports have print-friendly layouts.

### 7. Assign a Client

Before creating an estimate, a client must be assigned to the quote. From the job page, the user clicks the client card to select from existing clients or create a new one.

Client records store: name, company, email, phone, full address (street, city, state, zip), tax ID, and internal notes. Clients can only be changed while the quote is in draft status — once an estimate is sent, the client is locked.

The Clients section of the app provides full CRUD with paginated list and search by name. A client's detail page shows all associated quotes with their status.

### 8. Create an Estimate

An estimate is a snapshot of the quote at a point in time. When the user clicks "Create Estimate," the system freezes the current state of the quote — categories, items, totals — into a versioned record. Estimates auto-version: v1, v2, v3.

The estimate flattens the 3-level category hierarchy into 2 levels for presentation. Tier 3 items roll up into their parent tier 2 category. This keeps the client-facing document clean while the working quote can be as detailed as needed.

**Estimate categories have editable descriptions.** This is where the contractor writes scope-of-work language: "Supply and install 200A main panel with 40-space load center. Includes all breakers, grounding, and bonding per NEC 2023." Descriptions are edited inline with a click-to-edit interface — no separate form or modal.

### 9. Preview and Send

The estimate preview is a print-friendly, client-facing document showing:

- Company header
- Client name, company, and full address
- Category table with descriptions and totals
- Grand total
- Notes section
- Professional footer

The user reviews the preview, then clicks "Send for Signature."

### 10. E-Signature

The signature flow is a complete digital signing system:

1. **Send**: The user enters the client's name and email, optionally adds a personal message. The system generates a secure signing URL with a cryptographic token, creates an immutable snapshot of the estimate, and computes a SHA-256 document hash.

2. **Client views**: The client opens the public link (no login required). They see the company name, estimate details, scope of work by category, and any message from the contractor. The page handles edge cases — expired links, cancelled requests, already-signed documents.

3. **Client signs**: The client checks a consent box, enters their legal name, and submits. The system records the signature with: legal name, IP address, user agent, document hash, and timestamp. This creates an immutable record.

4. **Automatic status update**: The estimate is marked "Accepted" and the quote status updates. The client sees a confirmation page.

Signature requests expire after 30 days. The contractor can cancel a pending request and re-send.

---

## Supporting Features

### Item Template Library

Frequently used items are saved as templates with a type, category label, name, default unit, and default price. Templates serve two purposes:

1. **Autocomplete during item creation** — Search-as-you-type finds matching templates and pre-fills form fields
2. **Batch creation** — All templates in a category can be added to a quote category in one action

The template library is managed from the Items section of the app, with search and filtering by type and category.

### AI-Powered Price Import

Vendor price sheets arrive as Excel files in inconsistent formats. The import system handles this:

1. **Upload** an Excel file (up to 10MB)
2. **AI parsing** — Claude reads the spreadsheet as tab-separated text, identifies items, and matches them against existing templates with confidence scores
3. **Review** — The user sees matches organized by confidence level:
   - 90-100%: Exact match (auto-approved by default)
   - 70-89%: Strong match
   - 50-69%: Probable match
   - Below 50%: Weak match
4. **Actions** — Approve/reject individual matches, bulk approve above a threshold, or create new templates from unmatched rows
5. **Apply** — Approved matches update existing template prices

The parser handles hierarchical spreadsheets where category headers are mixed in with product rows. It uses heuristics (dimensions, unit patterns, fractions) to distinguish products from headers and prepends category context to item names.

### Organization Settings

- **Default markup percentage** — Applied to new quotes
- **Surcharge mode** — Stacking vs. override, with explanatory text and examples
- Saved via HTMX with toast notification feedback

### Multi-Tenancy and Authentication

- **Registration**: Organization name, subdomain, user details. First user gets admin role.
- **Login**: Email/password with Argon2id hashing, session tokens stored as SHA-256 hashes
- **Sessions**: 30-day default lifetime, secure cookie-based
- **Tenant isolation**: Every database query scoped by organization ID extracted from session
- **Roles**: Owner, admin, member

---

## UI/UX Design Principles

### Keyboard-First, Mouse-Friendly

The interface was designed for speed. A contractor pricing a job wants to move fast — navigating categories, adding items, adjusting quantities. Keyboard shortcuts handle all common actions:

| Key | Action |
|-----|--------|
| `n` | New quote |
| `c` | New category |
| `m` / `l` / `e` | Add material / labor / equipment |
| `b` | Batch add items |
| `j` / `k` | Navigate down / up (vim-style) |
| `Enter` | Open or edit selected row |
| `Space` | Expand/collapse category |
| `r` | Rename quote or category |
| `%` | Edit markup |
| `d` | Delete selected item |
| `o` / `s` | Open order list / site materials |
| `Escape` | Go back or close form |
| `?` | Show all keyboard shortcuts |

Every action also works with mouse/touch. Action menus (three-dot icons) provide the same operations. Mobile gets 44px minimum tap targets, a slide-in navigation sidebar, and stacked layouts for item rows.

### Spreadsheet View

The job page is a spreadsheet-style interface — not a literal spreadsheet, but designed to feel like one. Columns for Name, Qty, Unit, Price, Total. Categories as bold header rows. Items indented under their category. Inline forms for adding and editing without page navigation.

This was a deliberate choice: contractors live in spreadsheets. The UI should feel familiar, not like learning new software.

### Color as Information

The three standard item types each have a consistent color throughout the app:

- **Forest green** = materials (filled circle indicator)
- **Copper** = labor (half-circle indicator)
- **Slate** = equipment (empty circle indicator)

These colors appear on: item row hover highlights, type section headers, tag group backgrounds, dot indicators, and the brand accent (copper). Custom types get their own assigned color from a curated palette.

### Inline Everything

Almost nothing in Skalkaho opens a separate page or modal for editing. Rename a quote? Inline form replaces the title. Edit markup? Inline form replaces the percentage. Edit an item? The row transforms into an editable form. Add an item? A form slides in below the category header.

This keeps the user in context. They see the quote totals update in real time as they work, without navigating away and losing their place.

### HTMX-Driven Interactions

The frontend uses HTMX for all dynamic behavior. Form submissions, inline edits, filter changes, and delete operations happen without full page reloads. Alpine.js handles client-side state: expand/collapse, menu open/close, form toggle, consent checkbox state. There is no JavaScript framework, no build step, no bundle.

---

## Visual Identity

**Name origin**: Skalkaho Pass, a mountain route east of Hamilton, Montana in the Bitterroot Valley. Regional identity, not tech branding.

**Typography**: Barlow (UI text) + JetBrains Mono (numbers, monospace). Clean and readable.

**Color palette**: Drawn from the Montana landscape — granite slate, evergreen forest, copper mining heritage. Grounded, natural, professional.

---

## What We Learned Building This

The features in Skalkaho exist because a real contractor told us what they needed, often by showing us what was broken or slow in their spreadsheet workflow. Key lessons:

1. **Markup is not simple.** Different contractors mark up materials, labor, and equipment at different rates. Some stack overhead and profit separately, others use a single flat rate. The system had to support both without being confusing.

2. **Speed matters more than polish.** A contractor pricing a $200k remodel doesn't want to click through forms. Keyboard shortcuts and inline editing aren't nice-to-haves — they're the reason someone would use this instead of Excel.

3. **The quote is not the estimate.** The working quote can be messy — deep nesting, dozens of line items per category, notes and tags for internal tracking. The estimate the client sees should be clean and professional. Snapshotting and flattening the hierarchy bridges this gap.

4. **Vendor price sheets are chaos.** No two suppliers format their Excel sheets the same way. AI-powered parsing was the only realistic answer to importing prices without manual cleanup.

5. **Categories mirror how the work gets done.** Three levels of nesting handle the real structure of construction work (Phase > System > Component) without becoming unwieldy. Tags handle sub-grouping without adding a fourth level.

6. **Reports serve different audiences.** The Order List is for the supplier. Site Materials is for the crew. The Estimate is for the client. Same data, three views.

7. **E-signature closes the loop.** A quote sitting in someone's inbox waiting for a "looks good" reply is not a signed contract. Digital signature with immutable snapshots and audit trails turns a quote into a commitment.
