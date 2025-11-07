package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/config"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/handlers"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/middleware"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/providers"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize AI providers
	providerRegistry := providers.NewRegistry()

	// Register Google Gemini provider if API key is available
	if cfg.GoogleAPIKey != "" {
		geminiProvider, err := providers.NewGeminiProvider(cfg.GoogleAPIKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize Gemini provider: %v", err)
		} else {
			providerRegistry.Register("google", geminiProvider)
			log.Println("✓ Gemini provider registered")
		}
	}

	// Register OpenAI provider if API key is available
	if cfg.OpenAIAPIKey != "" {
		openaiProvider := providers.NewOpenAIProvider(cfg.OpenAIAPIKey)
		providerRegistry.Register("openai", openaiProvider)
		log.Println("✓ OpenAI provider registered")
	}

	// Register Anthropic Claude provider if API key is available
	if cfg.AnthropicAPIKey != "" {
		claudeProvider := providers.NewClaudeProvider(cfg.AnthropicAPIKey)
		providerRegistry.Register("anthropic", claudeProvider)
		log.Println("✓ Claude provider registered")
	}

	if len(providerRegistry.List()) == 0 {
		log.Fatal("No AI providers configured. Please set at least one API key.")
	}

	// Create handler
	handler := handlers.NewReviewHandler(providerRegistry, cfg)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/review", handler.HandleReview)

	// Apply middleware
	httpHandler := middleware.Logging(
		middleware.CORS(
			middleware.APIKeyAuth(mux, cfg.APIKeys),
		),
	)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 AI Gateway server starting on %s", addr)
	log.Printf("📋 Available providers: %v", providerRegistry.List())

	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"ai-gateway"}`)
}
