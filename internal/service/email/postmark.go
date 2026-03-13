package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client sends transactional emails via Postmark.
type Client struct {
	apiKey string
	from   string
}

// NewClient creates a new Postmark email client.
// If apiKey is empty, all sends will be no-ops (copy-link only mode).
func NewClient(apiKey, fromAddress string) *Client {
	return &Client{
		apiKey: apiKey,
		from:   fromAddress,
	}
}

// Enabled returns whether email sending is configured.
func (c *Client) Enabled() bool {
	return c.apiKey != ""
}

type postmarkMessage struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Subject  string `json:"Subject"`
	HtmlBody string `json:"HtmlBody"`
	TextBody string `json:"TextBody"`
}

type postmarkResponse struct {
	MessageID string `json:"MessageID"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// SendQuoteEmail sends a quote link to the recipient.
func (c *Client) SendQuoteEmail(to, quoteURL, projectName, contractorName string) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("email not configured")
	}

	subject := fmt.Sprintf("Quote from %s — %s", contractorName, projectName)

	htmlBody := fmt.Sprintf(`<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
	<h2 style="color: #1f2937;">Quote from %s</h2>
	<p style="color: #4b5563;">You have received a quote for <strong>%s</strong>.</p>
	<p style="color: #4b5563;">Please review the details and sign the quote using the link below:</p>
	<p style="margin: 24px 0;">
		<a href="%s" style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: 600;">
			View &amp; Sign Quote
		</a>
	</p>
	<p style="color: #9ca3af; font-size: 14px;">Or copy this link: %s</p>
	<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 24px 0;" />
	<p style="color: #9ca3af; font-size: 12px;">This email was sent by %s via Skalkaho.</p>
</div>`, contractorName, projectName, quoteURL, quoteURL, contractorName)

	textBody := fmt.Sprintf("Quote from %s\n\nYou have received a quote for %s.\n\nView and sign: %s\n\n— %s",
		contractorName, projectName, quoteURL, contractorName)

	msg := postmarkMessage{
		From:     c.from,
		To:       to,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshaling email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending email: %w", err)
	}
	defer resp.Body.Close()

	var result postmarkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if result.ErrorCode != 0 {
		return "", fmt.Errorf("postmark error %d: %s", result.ErrorCode, result.Message)
	}

	return result.MessageID, nil
}
