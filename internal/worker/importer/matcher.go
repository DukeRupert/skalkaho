package importer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// MatchResult holds Claude's assessment for a single unmatched row.
type MatchResult struct {
	RowIndex      int
	CandidateID   string  // empty if no match
	Confidence    float64
	Reasoning     string
	CanonicalName string // normalized name in "Category, Dimensions, Details" format
}

type matchResponse struct {
	Matches []matchResultJSON `json:"matches"`
}

type matchResultJSON struct {
	RowIndex      int     `json:"row_index"`
	CandidateID   *string `json:"candidate_id"`
	Confidence    float64 `json:"confidence"`
	Reasoning     string  `json:"reasoning"`
	CanonicalName string  `json:"canonical_name"`
}

type existingMaterialJSON struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Unit  string  `json:"unit"`
	Price float64 `json:"price"`
}

type unmatchedRowJSON struct {
	RowIndex int     `json:"row_index"`
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Price    float64 `json:"price"`
}

const namingRules = `
NAMING CONVENTION — apply to every canonical_name you return:
Format: "Category, Dimensions/Weight, Details"
The category (what the item IS) always comes first so searching by category
returns all sizes. Dimensions use trade shorthand. Species, brand, or grade last.

Examples:
  "2x4x8 Douglas Fir"           → "Stud, 2x4x8, Douglas Fir"
  "2x4 Stud 12ft DF"            → "Stud, 2x4x12, Douglas Fir"
  "2x6-16 #2 SPF"               → "Lumber, 2x6x16, #2 SPF"
  "3/4 CDX Plywood 4x8"         → "Plywood, 4x8, 3/4 CDX"
  "OSB 7/16 4x8"                → "OSB, 4x8, 7/16"
  "Simpson LUS210 Joist Hanger"  → "Joist Hanger, 2x10, Simpson LUS210"
  "Quikrete 80lb Concrete Mix"   → "Concrete Mix, 80lb, Quikrete"
  "R-19 Kraft Faced Batt 6.25x23"→ "Insulation Batt, 6-1/4x23, R-19 Kraft Faced"
  "#10 x 3 GRK RSS Screw"       → "Screw, #10x3, GRK RSS"
  "Tapcon 1/4 x 2-3/4"          → "Concrete Anchor, 1/4x2-3/4, Tapcon"
  "16d Framing Nail 3-1/2 Bright"→ "Nail, 16d 3-1/2, Bright Framing"
  "PL Premium Adhesive 28oz"     → "Construction Adhesive, 28oz, PL Premium"

Rules:
- Category is the common trade name for what the item IS (Stud, Plywood, Nail, etc.)
- Dimensions use trade format: lumber as WxDxL, sheet goods as WxL, fasteners as gauge/dia x length
- Convert feet spelled out ("12ft", "12'") to just the number in dimension shorthand
- Spell out species/grade abbreviations: DF → Douglas Fir, SPF stays SPF
- If the item doesn't fit neatly, use your best judgment — category first is the priority`

const systemPrompt = `You are a materials matching assistant for a construction quoting application.
You will be given a list of existing materials and a list of unmatched rows
from a supplier spreadsheet. For each unmatched row:
1. Identify the best matching existing material and return a confidence score (0.0–1.0)
2. Return a canonical_name using the naming convention below

` + namingRules + `

Additional matching rules:
- Use name similarity as the primary signal
- Use unit as a supporting signal (ea, lf, sf, etc.)
- Use price proximity as a weak supporting signal only
- Return null for candidate_id if no reasonable match exists (confidence < 0.50)
- Return only valid JSON, no explanation outside the JSON structure

Response format:
{
  "matches": [
    {
      "row_index": 0,
      "candidate_id": "abc123" | null,
      "confidence": 0.94,
      "reasoning": "Name is nearly identical, units match",
      "canonical_name": "Stud, 2x4x8, Douglas Fir"
    }
  ]
}`

const normalizePrompt = `You are a materials naming assistant for a construction quoting application.
You will be given a list of raw material names from a supplier spreadsheet.
Return a canonical_name for each row using the naming convention below.

` + namingRules + `

Return only valid JSON, no explanation outside the JSON structure.

Response format:
{
  "names": [
    {
      "row_index": 0,
      "canonical_name": "Stud, 2x4x8, Douglas Fir"
    }
  ]
}`

// BatchMatch sends unmatched rows to Claude for fuzzy matching against existing materials.
// Returns a result for each input row, including a canonical name.
func BatchMatch(ctx context.Context, rows []RawRow, existing []repository.Material) ([]MatchResult, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	// If no existing materials, everything is a new material — still normalize names
	if len(existing) == 0 {
		names, err := NormalizeNames(ctx, rows)
		if err != nil {
			// Fall back to raw names
			results := make([]MatchResult, len(rows))
			for i, row := range rows {
				results[i] = MatchResult{
					RowIndex:      row.RowIndex,
					Confidence:    0.0,
					Reasoning:     "No existing materials to match against",
					CanonicalName: row.Name,
				}
			}
			return results, nil
		}
		results := make([]MatchResult, len(rows))
		for i, row := range rows {
			canonical := row.Name
			if n, ok := names[row.RowIndex]; ok {
				canonical = n
			}
			results[i] = MatchResult{
				RowIndex:      row.RowIndex,
				Confidence:    0.0,
				Reasoning:     "No existing materials to match against",
				CanonicalName: canonical,
			}
		}
		return results, nil
	}

	// Build the user message payload
	existingJSON := make([]existingMaterialJSON, len(existing))
	for i, m := range existing {
		existingJSON[i] = existingMaterialJSON{
			ID:    m.ID,
			Name:  m.Name,
			Unit:  m.Unit,
			Price: m.UnitPrice,
		}
	}

	unmatchedJSON := make([]unmatchedRowJSON, len(rows))
	for i, r := range rows {
		unmatchedJSON[i] = unmatchedRowJSON{
			RowIndex: r.RowIndex,
			Name:     r.Name,
			Unit:     r.Unit,
			Price:    r.Price,
		}
	}

	payload := struct {
		ExistingMaterials []existingMaterialJSON `json:"existing_materials"`
		UnmatchedRows     []unmatchedRowJSON     `json:"unmatched_rows"`
	}{
		ExistingMaterials: existingJSON,
		UnmatchedRows:     unmatchedJSON,
	}

	userMsg, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling match request: %w", err)
	}

	// Call Claude API
	client := anthropic.NewClient()
	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(string(userMsg))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("calling Claude API: %w", err)
	}

	// Extract text content from response
	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText = block.Text
			break
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("empty response from Claude")
	}

	// Parse the JSON response
	var resp matchResponse
	if err := json.Unmarshal([]byte(responseText), &resp); err != nil {
		return nil, fmt.Errorf("parsing Claude response: %w (response: %s)", err, responseText)
	}

	results := make([]MatchResult, len(resp.Matches))
	for i, m := range resp.Matches {
		result := MatchResult{
			RowIndex:      m.RowIndex,
			Confidence:    m.Confidence,
			Reasoning:     m.Reasoning,
			CanonicalName: m.CanonicalName,
		}
		if m.CandidateID != nil {
			result.CandidateID = *m.CandidateID
		}
		results[i] = result
	}

	return results, nil
}

type normalizeResponse struct {
	Names []normalizeResultJSON `json:"names"`
}

type normalizeResultJSON struct {
	RowIndex      int    `json:"row_index"`
	CanonicalName string `json:"canonical_name"`
}

// NormalizeNames sends raw material names to Claude for canonical name formatting.
// Returns a map of row_index → canonical_name.
func NormalizeNames(ctx context.Context, rows []RawRow) (map[int]string, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	type nameRow struct {
		RowIndex int    `json:"row_index"`
		Name     string `json:"name"`
		Unit     string `json:"unit,omitempty"`
	}

	payload := struct {
		Rows []nameRow `json:"rows"`
	}{
		Rows: make([]nameRow, len(rows)),
	}
	for i, r := range rows {
		payload.Rows[i] = nameRow{
			RowIndex: r.RowIndex,
			Name:     r.Name,
			Unit:     r.Unit,
		}
	}

	userMsg, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling normalize request: %w", err)
	}

	client := anthropic.NewClient()
	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: normalizePrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(string(userMsg))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("calling Claude API for normalization: %w", err)
	}

	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText = block.Text
			break
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("empty response from Claude")
	}

	var resp normalizeResponse
	if err := json.Unmarshal([]byte(responseText), &resp); err != nil {
		return nil, fmt.Errorf("parsing normalize response: %w (response: %s)", err, responseText)
	}

	result := make(map[int]string, len(resp.Names))
	for _, n := range resp.Names {
		if n.CanonicalName != "" {
			result[n.RowIndex] = n.CanonicalName
		}
	}

	return result, nil
}
