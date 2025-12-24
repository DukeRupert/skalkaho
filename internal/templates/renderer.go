package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

//go:embed layouts/*.html pages/*.html partials/*.html
var templateFS embed.FS

// Renderer handles template rendering.
type Renderer struct {
	templates *template.Template
}

// NewRenderer creates a new template renderer.
func NewRenderer() (*Renderer, error) {
	// Parse all templates
	tmpl, err := template.New("").Funcs(templateFuncs()).ParseFS(templateFS, "layouts/*.html", "pages/*.html", "partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return &Renderer{templates: tmpl}, nil
}

// Render renders a full page template.
func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	return nil
}

// RenderPartial renders a partial template without the base layout.
func (r *Renderer) RenderPartial(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.templates.ExecuteTemplate(w, filepath.Base(name), data); err != nil {
		return fmt.Errorf("executing partial %s: %w", name, err)
	}
	return nil
}

// RenderToWriter renders a template to any io.Writer.
func (r *Renderer) RenderToWriter(w io.Writer, name string, data interface{}) error {
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	return nil
}

// templateFuncs returns custom template functions.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatMoney":   formatMoney,
		"formatPercent": formatPercent,
		"add":           func(a, b int) int { return a + b },
		"sub":           func(a, b int) int { return a - b },
		"mul":           func(a, b float64) float64 { return a * b },
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"deref": func(p *float64) float64 {
			if p == nil {
				return 0
			}
			return *p
		},
		"derefStr": func(p *string) string {
			if p == nil {
				return ""
			}
			return *p
		},
		"isNil": func(p interface{}) bool {
			return p == nil
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			m := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				m[key] = values[i+1]
			}
			return m
		},
	}
}

// formatMoney formats a float as currency.
func formatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// formatPercent formats a float as a percentage.
func formatPercent(amount float64) string {
	return fmt.Sprintf("%.1f%%", amount)
}
