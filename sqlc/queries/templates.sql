-- Templates

-- name: ListTemplates :many
SELECT * FROM templates ORDER BY name ASC;

-- name: GetTemplate :one
SELECT * FROM templates WHERE id = $1;

-- name: CreateTemplate :one
INSERT INTO templates (id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE templates
SET name = $2, description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM templates WHERE id = $1;

-- Template Sections

-- name: ListTemplateSections :many
SELECT * FROM template_sections
WHERE template_id = $1
ORDER BY sort_order ASC;

-- name: CreateTemplateSection :one
INSERT INTO template_sections (id, template_id, name, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTemplateSection :one
UPDATE template_sections
SET name = $2, sort_order = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTemplateSection :exec
DELETE FROM template_sections WHERE id = $1;

-- name: CountTemplateSections :one
SELECT COUNT(*) FROM template_sections WHERE template_id = $1;

-- Template Subcategories

-- name: ListTemplateSubcategories :many
SELECT * FROM template_subcategories
WHERE template_section_id = $1
ORDER BY sort_order ASC;

-- name: CreateTemplateSubcategory :one
INSERT INTO template_subcategories (id, template_section_id, name, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTemplateSubcategory :one
UPDATE template_subcategories
SET name = $2, sort_order = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTemplateSubcategory :exec
DELETE FROM template_subcategories WHERE id = $1;

-- Template Component Groups

-- name: ListTemplateComponentGroups :many
SELECT * FROM template_component_groups
WHERE template_subcategory_id = $1
ORDER BY sort_order ASC;

-- name: CreateTemplateComponentGroup :one
INSERT INTO template_component_groups (id, template_subcategory_id, name, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTemplateComponentGroup :one
UPDATE template_component_groups
SET name = $2, sort_order = $3
WHERE id = $1
RETURNING *;

-- name: DeleteTemplateComponentGroup :exec
DELETE FROM template_component_groups WHERE id = $1;
