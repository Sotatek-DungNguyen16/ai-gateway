package providers

import (
	"context"
	"fmt"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
)

// AIProvider defines the interface for AI providers
type AIProvider interface {
	Review(ctx context.Context, request *models.ReviewRequest) (*models.AIProviderResponse, error)
	Name() string
	SupportedModels() []string
}

// Registry manages AI providers
type Registry struct {
	providers map[string]AIProvider
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]AIProvider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(name string, provider AIProvider) {
	r.providers[name] = provider
}

// Get retrieves a provider by name
func (r *Registry) Get(name string) (AIProvider, error) {
	provider, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	return provider, nil
}

// List returns all registered provider names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

