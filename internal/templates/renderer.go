package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed layouts/*.html partials/*.html pages/*.html
var templateFS embed.FS

// Renderer parses and renders HTML templates.
// Pages are composed with layouts and partials using template inheritance.
type Renderer struct {
	templates map[string]*template.Template
}

// NewRenderer creates a new Renderer by parsing all page templates.
// Each page is parsed together with layouts and partials so that
// {{template "base" .}} and {{template "sidebar" .}} work.
func NewRenderer() (*Renderer, error) {
	r := &Renderer{
		templates: make(map[string]*template.Template),
	}

	pages, err := templateFS.ReadDir("pages")
	if err != nil {
		return nil, fmt.Errorf("reading pages directory: %w", err)
	}

	for _, page := range pages {
		if page.IsDir() {
			continue
		}
		name := page.Name()

		// Parse the page together with layouts and partials
		tmpl, err := template.ParseFS(templateFS,
			"pages/"+name,
			"layouts/*.html",
			"partials/*.html",
		)
		if err != nil {
			return nil, fmt.Errorf("parsing template %s: %w", name, err)
		}
		r.templates[name] = tmpl
	}

	return r, nil
}

// Render executes a named template and writes it to the response.
// For pages that use the base layout, the page template should call
// {{template "base" .}} which renders the full layout with sidebar.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}

// RenderPartial executes a named block within a page template.
// Used for HTMX partial responses that return just a fragment.
func (r *Renderer) RenderPartial(w http.ResponseWriter, page, block string, data any) error {
	tmpl, ok := r.templates[page]
	if !ok {
		return fmt.Errorf("template %q not found", page)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, block, data)
}
