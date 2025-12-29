# Skalkaho Roadmap

## MVP Status: Feature Complete

The MVP is functionally complete and deployed. Core quoting functionality is working.

---

## Completed Features

### Core Functionality
- [x] Quote management (create, rename, delete)
- [x] Hierarchical categories (up to 3 levels deep)
- [x] Line items: materials, labor, equipment
- [x] Surcharge system with inheritance (job → category → item)
- [x] Surcharge modes: stacking and override
- [x] Real-time total calculations
- [x] Settings page for defaults
- [x] Quote status tracking (draft, sent, accepted, rejected)
- [x] Quote pagination with filtering by status
- [x] Quote sorting (newest, oldest, name A-Z/Z-A)

### Client Management
- [x] Client CRUD with pagination and search
- [x] Client fields: name, company, email, phone, address, tax ID, notes
- [x] Associate quotes with clients
- [x] Client card display on quote detail page
- [x] Restrict client changes to draft quotes only
- [x] Auto-migrate legacy customer_name to clients

### Item Templates
- [x] Save commonly used items as templates
- [x] Search templates when adding line items
- [x] Template management page (/items)

### Reports
- [x] Order List - aggregated materials/equipment for entire quote
- [x] Site Materials - materials/equipment by category

### UI/UX
- [x] Keyboard-driven interface (vim-style navigation)
- [x] Inline editing for line items
- [x] Rename functionality (r key) for quotes and categories
- [x] Markup editing (% key) for quotes and categories
- [x] Delete functionality (d key)
- [x] Help overlay (? key)
- [x] Autocomplete search for line items
- [x] Print-friendly report layouts
- [x] Touch-friendly action menus (44px tap targets)
- [x] Mobile-responsive new quote button

### Infrastructure
- [x] Docker deployment with multi-stage build
- [x] Caddy reverse proxy configuration
- [x] GitHub Actions CI/CD pipeline
- [x] SQLite database with migrations

### Design
- [x] Typography system (Barlow + JetBrains Mono)
- [x] Color palette (slate, forest, copper)
- [x] Brand identity and logo

---

## Phase 2: Quote Workflow & Export

| Feature | Priority | Notes |
|---------|----------|-------|
| Send quote to client | High | Email quote to client, update status |
| PDF export | High | Generate professional quote PDFs |
| Quote duplication | Medium | Copy existing quote as template |
| Data backup/restore | Medium | Export/import database |

---

## Phase 3: Catalog & Efficiency

| Feature | Priority | Notes |
|---------|----------|-------|
| Labor rate presets | High | Common labor types with rates |
| Category templates | Medium | Reusable category structures |
| Supplier price imports | Low | CSV/Excel import for pricing |

---

## Phase 4: Business Features

| Feature | Priority | Notes |
|---------|----------|-------|
| Quote versioning | Medium | Track revisions |
| Client reports | Medium | Quote history by client |
| QuickBooks integration | Low | Export to accounting |
| User accounts | Low | Multi-user support |

---

## Known Issues / Polish

- [ ] Mobile responsiveness improvements
- [ ] Category drag-and-drop reordering
- [ ] Bulk line item operations
- [ ] Undo/redo support

---

## Technical Debt

- [ ] Add comprehensive test coverage
- [ ] PostgreSQL migration for production scale
- [ ] API documentation
- [ ] Performance profiling for large quotes

---

## Revision History

| Date | Notes |
|------|-------|
| 2025-12-27 | Initial roadmap - MVP feature complete |
| 2025-12-28 | Added client management, quote status, item templates, touch UI |
