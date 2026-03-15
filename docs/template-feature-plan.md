# Project Template Feature — Implementation Plan

## Overview

Allow the owner to define named structural templates (sections → subcategories →
component groups) that can be stamped into a new project at creation time. Once
applied, the project structure is fully independent of the template — editing the
template never affects existing projects. Templates contain no line items; they
are skeleton structure only.

---

## Data Model

Three template tables mirror the existing project structure exactly:

```
templates
  └── template_sections          (template_id)
        └── template_subcategories    (template_section_id)
              └── template_component_groups  (template_subcategory_id)
```

No self-referencing parent_id. Each level is its own table with a foreign key
to its parent, consistent with how sections, subcategories, and component_groups
work on the project side.

---

## Step 1 — Database Migration

**File:** `migrations/XXXX_add_templates.sql`

```sql
CREATE TABLE templates (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE template_sections (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE template_subcategories (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    template_section_id INTEGER NOT NULL REFERENCES template_sections(id) ON DELETE CASCADE,
    name                TEXT NOT NULL,
    sort_order          INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE template_component_groups (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    template_subcategory_id  INTEGER NOT NULL REFERENCES template_subcategories(id) ON DELETE CASCADE,
    name                     TEXT NOT NULL,
    sort_order               INTEGER NOT NULL DEFAULT 0
);
```

**Test:** Run migration, confirm all four tables exist with correct foreign keys.
Manually insert a template with one section → subcategory → component group,
confirm CASCADE delete works (delete the template, all children gone).

---

## Step 2 — sqlc Queries

**File:** `db/query/templates.sql`

### Template CRUD

```sql
-- name: ListTemplates :many
SELECT * FROM templates ORDER BY name ASC;

-- name: GetTemplate :one
SELECT * FROM templates WHERE id = ?;

-- name: CreateTemplate :one
INSERT INTO templates (name, description)
VALUES (?, ?)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE templates
SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM templates WHERE id = ?;
```

### Sections

```sql
-- name: ListTemplateSections :many
SELECT * FROM template_sections
WHERE template_id = ?
ORDER BY sort_order ASC;

-- name: CreateTemplateSection :one
INSERT INTO template_sections (template_id, name, sort_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateTemplateSection :one
UPDATE template_sections
SET name = ?, sort_order = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTemplateSection :exec
DELETE FROM template_sections WHERE id = ?;
```

### Subcategories

```sql
-- name: ListTemplateSubcategories :many
SELECT * FROM template_subcategories
WHERE template_section_id = ?
ORDER BY sort_order ASC;

-- name: CreateTemplateSubcategory :one
INSERT INTO template_subcategories (template_section_id, name, sort_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateTemplateSubcategory :one
UPDATE template_subcategories
SET name = ?, sort_order = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTemplateSubcategory :exec
DELETE FROM template_subcategories WHERE id = ?;
```

### Component Groups

```sql
-- name: ListTemplateComponentGroups :many
SELECT * FROM template_component_groups
WHERE template_subcategory_id = ?
ORDER BY sort_order ASC;

-- name: CreateTemplateComponentGroup :one
INSERT INTO template_component_groups (template_subcategory_id, name, sort_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateTemplateComponentGroup :one
UPDATE template_component_groups
SET name = ?, sort_order = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTemplateComponentGroup :exec
DELETE FROM template_component_groups WHERE id = ?;
```

Run `sqlc generate` and confirm no errors.

**Test:** Confirm generated types. Spot-check that nullable columns (`description`)
come through as `sql.NullString`.

---

## Step 3 — Stamp Function

**File:** `handler/template/stamp.go`

This is the core of the feature. `StampTemplate` is called when a new project is
created with a template selected. It reads the full template tree and inserts
copied rows into the live project tables.

```go
// StampTemplate copies a template's structure into an existing project.
// All inserts run inside a single transaction — if anything fails, nothing is written.
func StampTemplate(ctx context.Context, q *db.Queries, tx *sql.Tx, templateID, projectID int64) error
```

Logic:

1. Load all `template_sections` for the template
2. For each section → insert a real `section` row (project_id, name, sort_order)
3. Load all `template_subcategories` for that section
4. For each subcategory → insert a real `subcategory` row (section_id = new id from step 2)
5. Load all `template_component_groups` for that subcategory
6. For each component group → insert a real `component_group` row (subcategory_id = new id from step 4)

Key point: the new IDs from each insert must be used as the parent FK for the
next level — do not carry over template IDs.

**Test:** Call `StampTemplate` with a template that has 2 sections, 3 subcategories
each, and 2 component groups each. Confirm the project has exactly 2 sections,
6 subcategories, and 12 component groups with correct parent linkage.

---

## Step 4 — Template Management Handler

**File:** `handler/template/template.go`

CRUD handler for managing templates and their structure. Routes:

| Method | Route                                        | Action                        |
|--------|----------------------------------------------|-------------------------------|
| GET    | `/templates`                                 | List all templates            |
| GET    | `/templates/new`                             | New template form             |
| POST   | `/templates`                                 | Create template               |
| GET    | `/templates/{id}`                            | Detail / edit view            |
| PUT    | `/templates/{id}`                            | Update template name/desc     |
| DELETE | `/templates/{id}`                            | Delete template               |
| POST   | `/templates/{id}/sections`                   | Add section                   |
| PUT    | `/templates/{id}/sections/{sid}`             | Update section                |
| DELETE | `/templates/{id}/sections/{sid}`             | Delete section (cascades)     |
| POST   | `/templates/{id}/sections/{sid}/subcategories`          | Add subcategory      |
| PUT    | `/templates/{id}/sections/{sid}/subcategories/{scid}`   | Update subcategory   |
| DELETE | `/templates/{id}/sections/{sid}/subcategories/{scid}`   | Delete subcategory   |
| POST   | `/templates/{id}/sections/{sid}/subcategories/{scid}/groups`        | Add group   |
| PUT    | `/templates/{id}/sections/{sid}/subcategories/{scid}/groups/{gid}`  | Update group|
| DELETE | `/templates/{id}/sections/{sid}/subcategories/{scid}/groups/{gid}`  | Delete group|

The detail view is the primary editing surface — the owner builds out the tree
structure inline using HTMX to add/remove/rename nodes at each level without
full page reloads.

**Test:** Hit each route, confirm correct sqlc method is called, no template
render errors.

---

## Step 5 — Templates UI

**Directory:** `templates/templates/`

### `index.html`

- List of templates: name, description, section count, created date
- "New Template" button
- Delete button per row (HTMX confirm + swap out)

### `form.html` (new template)

Simple form: name (required) + description (optional). On submit, redirects to
the detail/edit view where structure is built out.

### `show.html` (edit view)

The main template editor. Displays the full tree inline:

```
▼ Section: Framing                        [rename] [delete]
    ▼ Subcategory: Wall Framing           [rename] [delete]
        • Component Group: Exterior Walls [rename] [delete]
        • Component Group: Interior Walls [rename] [delete]
        [+ Add Component Group]
    [+ Add Subcategory]
▼ Section: Electrical                     [rename] [delete]
    ...
[+ Add Section]
```

Each add/rename/delete uses HTMX to swap the relevant fragment in place.
Inline renaming can be an Alpine.js toggled input — click the name to edit,
blur or Enter to submit via `hx-put`.

**Test:** Add a full three-level tree via the UI. Rename a node at each level.
Delete a subcategory and confirm its component groups are gone. Reload and
confirm persistence.

---

## Step 6 — Project Creation Integration

**File:** Wherever the new project / quote form lives.

Add an optional template selector to the new project form:

```html
<select name="template_id">
  <option value="">— No template —</option>
  {{ range .Templates }}
  <option value="{{ .ID }}">{{ .Name }}</option>
  {{ end }}
</select>
```

In the project creation handler:

1. Create the project row as normal
2. If `template_id` is present and valid, call `StampTemplate` inside the same
   transaction as the project insert
3. Redirect to the project / quote editor as usual

The user selects a template, lands in the editor, and the structure is already
there. No extra step, no separate "apply template" button.

**Test:** Create a project with a template selected, confirm structure is present
in the editor. Create a project without a template, confirm it starts empty.
Edit the template, create another project from it, confirm the old project is
unchanged.

---

## Step 7 — Router Wiring

**File:** `router/router.go`

Register template routes alongside existing ones. Add a nav link wherever
project/quote management lives — templates are an owner/admin concern, so
consider whether they belong in a settings or admin section vs. the main nav.

**Test:** Full end-to-end — create template, build structure, create project from
template, confirm stamp, edit template, confirm old project unaffected.

---

## Commit Sequence

| # | Commit message                                                        |
|---|-----------------------------------------------------------------------|
| 1 | `db: add templates, template_sections, template_subcategories, template_component_groups migration` |
| 2 | `db: add sqlc queries for template CRUD and tree nodes`               |
| 3 | `handler: add StampTemplate function with transaction support`        |
| 4 | `handler: add template management handler and routes`                 |
| 5 | `templates: add template list, form, and tree editor UI`              |
| 6 | `handler: integrate template stamp into project creation flow`        |
| 7 | `router: wire template routes, add nav link`                          |
