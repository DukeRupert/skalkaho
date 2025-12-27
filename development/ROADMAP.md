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

## Phase 2: Persistence & Export

| Feature | Priority | Notes |
|---------|----------|-------|
| PDF export | High | Generate professional quote PDFs |
| Quote duplication | Medium | Copy existing quote as template |
| Quote archiving | Medium | Mark quotes as won/lost/archived |
| Data backup/restore | Medium | Export/import database |

---

## Phase 3: Catalog & Efficiency

| Feature | Priority | Notes |
|---------|----------|-------|
| Product catalog | High | Saved materials with pricing |
| Labor rate presets | High | Common labor types with rates |
| Category templates | Medium | Reusable category structures |
| Supplier price imports | Low | CSV/Excel import for pricing |

---

## Phase 4: Business Features

| Feature | Priority | Notes |
|---------|----------|-------|
| Customer management | Medium | Track customers and their quotes |
| Quote versioning | Medium | Track revisions |
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
