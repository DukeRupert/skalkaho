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
		"add":           add,
		"sub":           sub,
		"mul":           func(a, b float64) float64 { return a * b },
		"eq":            func(a, b interface{}) bool { return a == b },
		"gt":            gt,
		"typeIndicator": typeIndicator,
		"dict":          dict,
	}
}

// add handles both int and int64 types
func add(a, b interface{}) int64 {
	return toInt64(a) + toInt64(b)
}

// sub handles both int and int64 types
func sub(a, b interface{}) int64 {
	return toInt64(a) - toInt64(b)
}

// gt handles both int and int64 types
func gt(a, b interface{}) bool {
	return toInt64(a) > toInt64(b)
}

// toInt64 converts various numeric types to int64
func toInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	default:
		return 0
	}
}

// dict creates a map from key-value pairs for passing to templates.
func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		return nil
	}
	d := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			continue
		}
		d[key] = values[i+1]
	}
	return d
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
