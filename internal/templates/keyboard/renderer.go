package keyboard

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

//go:embed layouts/*.html pages/*.html partials/*.html
var templateFS embed.FS

// Renderer handles keyboard template rendering.
type Renderer struct {
	templates *template.Template
}

// NewRenderer creates a new keyboard template renderer.
func NewRenderer() (*Renderer, error) {
	tmpl, err := template.New("").Funcs(templateFuncs()).ParseFS(templateFS, "layouts/*.html", "pages/*.html", "partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing keyboard templates: %w", err)
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

// RenderPartial renders a partial template.
func (r *Renderer) RenderPartial(w io.Writer, name string, data interface{}) error {
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("executing partial %s: %w", name, err)
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
		"eq":            func(a, b interface{}) bool { return a == b },
		"typeIndicator": typeIndicator,
	}
}

func formatMoney(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func formatPercent(amount float64) string {
	return fmt.Sprintf("%.1f%%", amount)
}

func typeIndicator(itemType string) string {
	switch itemType {
	case "material":
		return "●" // filled circle
	case "labor":
		return "◐" // half circle
	case "equipment":
		return "○" // empty circle
	default:
		return "•"
	}
}
