package excel

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Row represents a parsed row from an Excel spreadsheet.
type Row struct {
	RowNumber int
	Name      string
	Unit      string
	Price     float64
}

// ParseResult contains the parsed data from an Excel file.
type ParseResult struct {
	Rows     []Row
	Filename string
}

// Parser handles Excel file parsing.
type Parser struct{}

// NewParser creates a new Excel parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads an Excel file and extracts item data.
// It attempts to auto-detect columns containing name, unit, and price data.
func (p *Parser) Parse(r io.Reader, filename string) (*ParseResult, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("opening excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}
	sheetName := sheets[0]

	// Get all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no data found in sheet")
	}

	// Detect column indices
	nameCol, unitCol, priceCol := p.detectColumns(rows)
	if nameCol == -1 {
		return nil, fmt.Errorf("could not detect name column")
	}
	if priceCol == -1 {
		return nil, fmt.Errorf("could not detect price column")
	}

	// Parse data rows (skip header)
	var parsedRows []Row
	startRow := 1 // Skip header row
	if len(rows) > 0 && !p.isHeaderRow(rows[0]) {
		startRow = 0
	}

	for i := startRow; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		name := ""
		if nameCol < len(row) {
			name = strings.TrimSpace(row[nameCol])
		}
		if name == "" {
			continue // Skip rows without names
		}

		unit := ""
		if unitCol >= 0 && unitCol < len(row) {
			unit = strings.TrimSpace(row[unitCol])
		}

		price := 0.0
		if priceCol < len(row) {
			price = p.parsePrice(row[priceCol])
		}

		// Only include rows with valid price
		if price <= 0 {
			continue
		}

		parsedRows = append(parsedRows, Row{
			RowNumber: i + 1, // 1-indexed for user display
			Name:      name,
			Unit:      unit,
			Price:     price,
		})
	}

	if len(parsedRows) == 0 {
		return nil, fmt.Errorf("no valid data rows found (need name and price)")
	}

	return &ParseResult{
		Rows:     parsedRows,
		Filename: filename,
	}, nil
}

// detectColumns attempts to identify which columns contain name, unit, and price.
func (p *Parser) detectColumns(rows [][]string) (nameCol, unitCol, priceCol int) {
	nameCol = -1
	unitCol = -1
	priceCol = -1

	if len(rows) == 0 {
		return
	}

	// Check first row for headers
	header := rows[0]
	for i, cell := range header {
		cellLower := strings.ToLower(strings.TrimSpace(cell))

		// Name column
		if nameCol == -1 && (strings.Contains(cellLower, "name") ||
			strings.Contains(cellLower, "description") ||
			strings.Contains(cellLower, "item") ||
			strings.Contains(cellLower, "product") ||
			strings.Contains(cellLower, "material")) {
			nameCol = i
			continue
		}

		// Unit column
		if unitCol == -1 && (strings.Contains(cellLower, "unit") ||
			strings.Contains(cellLower, "uom") ||
			strings.Contains(cellLower, "measure")) {
			unitCol = i
			continue
		}

		// Price column
		if priceCol == -1 && (strings.Contains(cellLower, "price") ||
			strings.Contains(cellLower, "cost") ||
			strings.Contains(cellLower, "rate") ||
			strings.Contains(cellLower, "amount")) {
			priceCol = i
			continue
		}
	}

	// If we couldn't detect from headers, try to infer from data
	if nameCol == -1 || priceCol == -1 {
		nameCol, unitCol, priceCol = p.inferColumnsFromData(rows)
	}

	return
}

// inferColumnsFromData attempts to detect columns by analyzing the data.
func (p *Parser) inferColumnsFromData(rows [][]string) (nameCol, unitCol, priceCol int) {
	nameCol = -1
	unitCol = -1
	priceCol = -1

	if len(rows) < 2 {
		return
	}

	// Sample a few data rows to determine column types
	maxSamples := 5
	if len(rows) < maxSamples {
		maxSamples = len(rows)
	}

	// Count how many cells in each column look like numbers
	numericCount := make(map[int]int)
	textCount := make(map[int]int)

	for i := 1; i < maxSamples; i++ {
		row := rows[i]
		for j, cell := range row {
			cell = strings.TrimSpace(cell)
			if cell == "" {
				continue
			}
			if p.parsePrice(cell) > 0 {
				numericCount[j]++
			} else if len(cell) > 3 {
				textCount[j]++
			}
		}
	}

	// Find the column with the most text (likely name)
	maxText := 0
	for col, count := range textCount {
		if count > maxText {
			maxText = count
			nameCol = col
		}
	}

	// Find the column with the most numbers (likely price)
	maxNum := 0
	for col, count := range numericCount {
		if count > maxNum && col != nameCol {
			maxNum = count
			priceCol = col
		}
	}

	return
}

// isHeaderRow checks if a row appears to be a header row.
func (p *Parser) isHeaderRow(row []string) bool {
	if len(row) == 0 {
		return false
	}

	headerKeywords := []string{"name", "description", "item", "product", "price", "cost", "unit", "qty", "quantity"}
	for _, cell := range row {
		cellLower := strings.ToLower(strings.TrimSpace(cell))
		for _, keyword := range headerKeywords {
			if strings.Contains(cellLower, keyword) {
				return true
			}
		}
	}
	return false
}

// parsePrice attempts to parse a price value from a string.
func (p *Parser) parsePrice(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// Remove currency symbols and formatting
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
