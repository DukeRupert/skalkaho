package importer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// RawRow represents a single parsed row from a supplier spreadsheet.
type RawRow struct {
	Name         string
	Unit         string
	Price        float64
	SupplierCode string
	RowIndex     int // 1-based row number in the spreadsheet (for error reporting)
}

// ParseExcel reads an Excel file and returns parsed rows.
// It reads the first sheet, skips the header row, and skips blank rows.
func ParseExcel(path string) ([]RawRow, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening excel file: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in file")
	}

	allRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}

	if len(allRows) < 2 {
		return nil, fmt.Errorf("spreadsheet has no data rows (only %d rows)", len(allRows))
	}

	// Detect column mapping from header row
	header := allRows[0]
	colMap := detectColumns(header)

	if colMap.name < 0 {
		return nil, fmt.Errorf("could not find a 'name' or 'description' column in header: %v", header)
	}

	var rows []RawRow
	for i := 1; i < len(allRows); i++ {
		row := allRows[i]
		raw := parseRow(row, colMap, i+1) // i+1 for 1-based row number
		if raw.Name == "" {
			continue // skip blank rows
		}
		rows = append(rows, raw)
	}

	return rows, nil
}

// columnMap maps header names to column indices.
type columnMap struct {
	name         int
	unit         int
	price        int
	supplierCode int
}

func detectColumns(header []string) columnMap {
	cm := columnMap{name: -1, unit: -1, price: -1, supplierCode: -1}

	for i, h := range header {
		key := strings.ToLower(strings.TrimSpace(h))
		switch {
		case contains(key, "name", "description", "product", "item"):
			if cm.name < 0 {
				cm.name = i
			}
		case contains(key, "unit") && !contains(key, "price", "cost"):
			if cm.unit < 0 {
				cm.unit = i
			}
		case contains(key, "price", "cost", "amount"):
			if cm.price < 0 {
				cm.price = i
			}
		case contains(key, "code", "sku", "part", "item #", "item#", "number"):
			if cm.supplierCode < 0 {
				cm.supplierCode = i
			}
		}
	}

	return cm
}

func contains(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func parseRow(row []string, cm columnMap, rowIndex int) RawRow {
	raw := RawRow{RowIndex: rowIndex}

	if cm.name >= 0 && cm.name < len(row) {
		raw.Name = strings.TrimSpace(row[cm.name])
	}
	if cm.unit >= 0 && cm.unit < len(row) {
		raw.Unit = normalizeUnit(row[cm.unit])
	}
	if cm.price >= 0 && cm.price < len(row) {
		raw.Price = parsePrice(row[cm.price])
	}
	if cm.supplierCode >= 0 && cm.supplierCode < len(row) {
		raw.SupplierCode = strings.TrimSpace(row[cm.supplierCode])
	}

	return raw
}

func parsePrice(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

var unitNormalizer = map[string]string{
	"each":        "ea",
	"linear foot": "lf",
	"linear feet": "lf",
	"lineal foot": "lf",
	"lineal feet": "lf",
	"square foot": "sf",
	"square feet": "sf",
	"board foot":  "bf",
	"board feet":  "bf",
	"cubic yard":  "cy",
	"cubic yards": "cy",
	"hour":        "hr",
	"hours":       "hr",
	"gallon":      "gal",
	"gallons":     "gal",
	"pound":       "lb",
	"pounds":      "lb",
}

func normalizeUnit(raw string) string {
	key := strings.ToLower(strings.TrimSpace(raw))
	if normalized, ok := unitNormalizer[key]; ok {
		return normalized
	}
	return key
}
