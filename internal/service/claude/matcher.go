package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/service/excel"
)

// ExtractedItem represents an item extracted from a spreadsheet by Claude.
type ExtractedItem struct {
	RowNumber int     `json:"row_number"`
	Name      string  `json:"name"`
	Unit      string  `json:"unit,omitempty"`
	Price     float64 `json:"price"`
}

// ExtractAndMatchResponse contains extracted items with their matches.
type ExtractAndMatchResponse struct {
	Items []ExtractedItemWithMatch `json:"items"`
}

// ExtractedItemWithMatch combines an extracted item with its template match.
type ExtractedItemWithMatch struct {
	RowNumber    int     `json:"row_number"`
	Name         string  `json:"name"`
	Unit         string  `json:"unit,omitempty"`
	Price        float64 `json:"price"`
	TemplateID   *int64  `json:"template_id,omitempty"`
	TemplateName string  `json:"template_name,omitempty"`
	Confidence   float64 `json:"confidence"`
	Reason       string  `json:"reason"`
}

// MatchResult represents a single match between a spreadsheet row and an item template.
type MatchResult struct {
	RowNumber    int     `json:"row_number"`
	TemplateID   *int64  `json:"template_id,omitempty"`
	TemplateName string  `json:"template_name,omitempty"`
	Confidence   float64 `json:"confidence"`
	Reason       string  `json:"reason"`
}

// MatchResponse contains all matches from Claude.
type MatchResponse struct {
	Matches []MatchResult `json:"matches"`
}

// Matcher handles matching spreadsheet items to templates using Claude AI.
type Matcher struct {
	client anthropic.Client
}

// NewMatcher creates a new Claude matcher.
func NewMatcher(apiKey string) *Matcher {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Matcher{client: client}
}

// MatchItems sends spreadsheet rows and templates to Claude for matching.
func (m *Matcher) MatchItems(ctx context.Context, rows []excel.Row, templates []repository.ItemTemplate) (*MatchResponse, error) {
	if len(rows) == 0 {
		return &MatchResponse{Matches: []MatchResult{}}, nil
	}

	prompt := m.buildPrompt(rows, templates)

	resp, err := m.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract text content from response
	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	textContent := ""
	for _, block := range resp.Content {
		if block.Type == "text" {
			textContent = block.Text
			break
		}
	}

	if textContent == "" {
		return nil, fmt.Errorf("no text content in Claude response")
	}

	// Parse JSON response
	result, err := m.parseResponse(textContent)
	if err != nil {
		return nil, fmt.Errorf("parsing claude response: %w", err)
	}

	return result, nil
}

// ExtractAndMatchItems extracts items from raw spreadsheet text and matches them against templates.
// This uses a single Claude API call to both parse the spreadsheet and match items.
func (m *Matcher) ExtractAndMatchItems(ctx context.Context, spreadsheet *excel.RawSpreadsheet, templates []repository.ItemTemplate) (*ExtractAndMatchResponse, error) {
	prompt := m.buildExtractAndMatchPrompt(spreadsheet, templates)

	resp, err := m.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 8192,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract text content from response
	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Claude")
	}

	textContent := ""
	for _, block := range resp.Content {
		if block.Type == "text" {
			textContent = block.Text
			break
		}
	}

	if textContent == "" {
		return nil, fmt.Errorf("no text content in Claude response")
	}

	// Parse JSON response
	result, err := m.parseExtractAndMatchResponse(textContent)
	if err != nil {
		return nil, fmt.Errorf("parsing claude response: %w", err)
	}

	return result, nil
}

func (m *Matcher) buildExtractAndMatchPrompt(spreadsheet *excel.RawSpreadsheet, templates []repository.ItemTemplate) string {
	var sb strings.Builder

	sb.WriteString(`You are a construction materials data extraction and matching assistant. Your task is to:
1. Extract product items from a supplier price list spreadsheet
2. Match each extracted item to existing item templates

## Important Instructions for Extraction
- The spreadsheet may have category headers (rows with a category name but no price)
- When you encounter a category header, PREPEND that category to all subsequent item names until a new category is found
- For example, if you see "Sheeting" as a category, then "3/8 CDX" should become "Sheeting 3/8 CDX"
- Only extract rows that have both a name AND a price
- Look for price columns (may be labeled "Price", "Cost", "Rate", or just contain dollar amounts)
- Look for unit columns (may be labeled "Unit", "UOM", "Measure")
- Be smart about identifying the actual product data vs. headers, totals, or notes

## Existing Item Templates
`)

	// Format templates as a list
	for _, t := range templates {
		sb.WriteString(fmt.Sprintf("- ID: %d, Name: %s, Unit: %s, Current Price: $%.2f\n",
			t.ID, t.Name, t.DefaultUnit, t.DefaultPrice))
	}

	sb.WriteString(`
## Raw Spreadsheet Content
`)
	sb.WriteString(spreadsheet.Content)

	sb.WriteString(`

## Instructions for Matching
After extracting items, match each one to the most appropriate template:
1. Compare names considering: abbreviations, common construction terminology, dimensions
2. Return confidence score (0.0-1.0):
   - 0.9-1.0: Exact or near-exact match (same product)
   - 0.7-0.89: Strong match (clearly the same item with different naming)
   - 0.5-0.69: Probable match (likely same item, needs verification)
   - 0.0-0.49: Weak or no match (different items or too uncertain)
3. Provide brief reason for match or non-match

## Response Format (JSON only, no other text)
{
  "items": [
    {
      "row_number": 5,
      "name": "Sheeting 3/8 CDX",
      "unit": "sheet",
      "price": 25.99,
      "template_id": 42,
      "template_name": "Sheeting 3/8 CDX Plywood",
      "confidence": 0.95,
      "reason": "Near-exact name match"
    },
    {
      "row_number": 6,
      "name": "Sheeting 1/2 CDX",
      "unit": "sheet",
      "price": 32.50,
      "template_id": null,
      "template_name": "",
      "confidence": 0.0,
      "reason": "No matching template found"
    }
  ]
}

Return ONLY valid JSON with no additional text or explanation.`)

	return sb.String()
}

func (m *Matcher) parseExtractAndMatchResponse(text string) (*ExtractAndMatchResponse, error) {
	// Try to extract JSON from the response
	text = strings.TrimSpace(text)

	// Handle markdown code blocks
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	var result ExtractAndMatchResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w (response was: %s)", err, text[:min(200, len(text))])
	}

	return &result, nil
}

func (m *Matcher) buildPrompt(rows []excel.Row, templates []repository.ItemTemplate) string {
	var sb strings.Builder

	sb.WriteString(`You are a construction materials matching assistant. Match items from a supplier price list to existing item templates.

## Existing Item Templates
`)

	// Format templates as a list
	for _, t := range templates {
		sb.WriteString(fmt.Sprintf("- ID: %d, Name: %s, Unit: %s, Current Price: $%.2f\n",
			t.ID, t.Name, t.DefaultUnit, t.DefaultPrice))
	}

	sb.WriteString(`
## Supplier Price List Items
`)

	// Format spreadsheet rows
	for _, r := range rows {
		if r.Unit != "" {
			sb.WriteString(fmt.Sprintf("- Row %d: %s, Unit: %s, Price: $%.2f\n",
				r.RowNumber, r.Name, r.Unit, r.Price))
		} else {
			sb.WriteString(fmt.Sprintf("- Row %d: %s, Price: $%.2f\n",
				r.RowNumber, r.Name, r.Price))
		}
	}

	sb.WriteString(`
## Instructions
1. For each supplier item, find the best matching template by comparing names
2. Consider: name similarity, abbreviations, common construction terminology
3. Return confidence score (0.0-1.0):
   - 0.9-1.0: Exact or near-exact match (same product, minor spelling/format differences)
   - 0.7-0.89: Strong match (clearly the same item with different naming convention)
   - 0.5-0.69: Probable match (likely the same item but needs human verification)
   - 0.0-0.49: Weak or no match (different items or too uncertain)
4. Provide brief reason for match or non-match

## Response Format (JSON only, no other text)
{
  "matches": [
    {
      "row_number": 1,
      "template_id": 42,
      "template_name": "2x4 Lumber 8ft",
      "confidence": 0.95,
      "reason": "Exact name match"
    },
    {
      "row_number": 2,
      "template_id": null,
      "template_name": "",
      "confidence": 0.0,
      "reason": "No matching template found"
    }
  ]
}

Return ONLY valid JSON with no additional text or explanation.`)

	return sb.String()
}

func (m *Matcher) parseResponse(text string) (*MatchResponse, error) {
	// Try to extract JSON from the response
	text = strings.TrimSpace(text)

	// Handle markdown code blocks
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}

	var result MatchResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w (response was: %s)", err, text[:min(200, len(text))])
	}

	return &result, nil
}
