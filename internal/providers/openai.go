package providers

import (
	"context"
	"fmt"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/prompt"
	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the AIProvider interface for OpenAI
type OpenAIProvider struct {
	client *openai.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client: client,
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// SupportedModels returns the list of supported models
func (p *OpenAIProvider) SupportedModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-4o",
		"gpt-3.5-turbo",
	}
}

// Review performs a code review using OpenAI
func (p *OpenAIProvider) Review(ctx context.Context, request *models.ReviewRequest) (*models.AIProviderResponse, error) {
	// Get model, default to gpt-4o if not specified
	modelName := request.AIModel
	if modelName == "" {
		modelName = "gpt-4o"
	}

	// Generate prompts
	systemPrompt := prompt.GenerateSystemPrompt(request.Language)
	userPrompt := prompt.GenerateUserPrompt(request)

	// Create chat completion request
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.3,
			MaxTokens:   4096,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := resp.Choices[0].Message.Content

	// Parse the response
	return prompt.ParseAIResponse(responseText)
}
