package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/prompt"
)

// ClaudeProvider implements the AIProvider interface for Anthropic Claude
type ClaudeProvider struct {
	apiKey     string
	httpClient *http.Client
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "anthropic"
}

// SupportedModels returns the list of supported models
func (p *ClaudeProvider) SupportedModels() []string {
	return []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}
}

// ClaudeRequest represents the request structure for Claude API
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	System      string          `json:"system,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// ClaudeMessage represents a message in Claude API
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the response from Claude API
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Review performs a code review using Claude
func (p *ClaudeProvider) Review(ctx context.Context, request *models.ReviewRequest) (*models.AIProviderResponse, error) {
	// Get model, default to claude-3-5-sonnet if not specified
	modelName := request.AIModel
	if modelName == "" {
		modelName = "claude-3-5-sonnet-20241022"
	}

	// Generate prompts
	systemPrompt := prompt.GenerateSystemPrompt(request.Language)
	userPrompt := prompt.GenerateUserPrompt(request)

	// Create request
	reqBody := ClaudeRequest{
		Model:       modelName,
		MaxTokens:   4096,
		Temperature: 0.3,
		System:      systemPrompt,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return nil, fmt.Errorf("Claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	responseText := claudeResp.Content[0].Text

	// Parse the response
	return prompt.ParseAIResponse(responseText)
}
