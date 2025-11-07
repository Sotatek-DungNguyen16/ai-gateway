package providers

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/prompt"
	"google.golang.org/api/option"
)

// GeminiProvider implements the AIProvider interface for Google Gemini
type GeminiProvider struct {
	client *genai.Client
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string) (*GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiProvider{
		client: client,
	}, nil
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "google"
}

// SupportedModels returns the list of supported models
func (p *GeminiProvider) SupportedModels() []string {
	return []string{
		"gemini-2.0-flash",
		"gemini-1.5-pro",
		"gemini-1.5-flash",
		"gemini-pro",
	}
}

// Review performs a code review using Gemini
func (p *GeminiProvider) Review(ctx context.Context, request *models.ReviewRequest) (*models.AIProviderResponse, error) {
	// Get model, default to gemini-2.0-flash if not specified
	modelName := request.AIModel
	if modelName == "" {
		modelName = "gemini-2.0-flash"
	}

	model := p.client.GenerativeModel(modelName)

	// Configure model for structured output
	model.SetTemperature(0.3)
	model.SetTopP(0.95)
	model.SetTopK(40)
	model.SetMaxOutputTokens(8192)

	// Generate prompt
	systemPrompt := prompt.GenerateSystemPrompt(request.Language)
	userPrompt := prompt.GenerateUserPrompt(request)

	// Create the prompt parts
	fullPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, userPrompt)

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates from Gemini")
	}

	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Parse the response
	return prompt.ParseAIResponse(responseText)
}
