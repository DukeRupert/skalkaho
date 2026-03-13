package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

//go:embed pages/*.html
var templateFS embed.FS

// Renderer parses and renders HTML templates.
type Renderer struct {
	templates map[string]*template.Template
}

// NewRenderer creates a new Renderer by parsing all page templates.
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
		tmpl, err := template.ParseFS(templateFS, "pages/"+name)
		if err != nil {
			return nil, fmt.Errorf("parsing template %s: %w", name, err)
		}
		r.templates[name] = tmpl
	}

	return r, nil
}

// Render executes a named template and writes it to the response.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}

// RenderToWriter executes a named template and writes to an io.Writer.
func (r *Renderer) RenderToWriter(w io.Writer, name string, data any) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	return tmpl.Execute(w, data)
}
