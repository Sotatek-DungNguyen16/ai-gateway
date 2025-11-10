package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/config"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/providers"
)

// ReviewHandler handles code review requests
type ReviewHandler struct {
	registry *providers.Registry
	config   *config.Config
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(registry *providers.Registry, cfg *config.Config) *ReviewHandler {
	return &ReviewHandler{
		registry: registry,
		config:   cfg,
	}
}

// HandleReview handles the /review endpoint
func (h *ReviewHandler) HandleReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Log request details for debugging
	contentType := r.Header.Get("Content-Type")
	log.Printf("Received review request - Content-Type: %s, Content-Length: %d", contentType, r.ContentLength)

	var request models.ReviewRequest

	// Handle both JSON and multipart/form-data
	if strings.Contains(contentType, "application/json") {
		// Handle JSON request (from GitHub Actions)
		log.Printf("Processing as JSON request")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &request); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Invalid JSON: %v"}`, err), http.StatusBadRequest)
			return
		}
	} else {
		// Handle multipart/form-data request (from local/curl)
		log.Printf("Processing as multipart/form-data request")
		if err := r.ParseMultipartForm(h.config.MaxDiffSize); err != nil {
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Failed to parse form: %v"}`, err), http.StatusBadRequest)
			return
		}

		// Get metadata
		metadataStr := r.FormValue("metadata")
		if metadataStr == "" {
			http.Error(w, `{"error":"Missing metadata field"}`, http.StatusBadRequest)
			return
		}

		// Parse metadata JSON
		if err := json.Unmarshal([]byte(metadataStr), &request); err != nil {
			log.Printf("Error parsing metadata JSON: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Invalid metadata JSON: %v"}`, err), http.StatusBadRequest)
			return
		}

		// Get git_diff file
		file, _, err := r.FormFile("git_diff")
		if err != nil {
			log.Printf("Error reading git_diff file: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"Missing or invalid git_diff file: %v"}`, err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read diff content
		diffBytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading diff content: %v", err)
			http.Error(w, `{"error":"Failed to read diff content"}`, http.StatusInternalServerError)
			return
		}
		request.GitDiff = string(diffBytes)
	}

	// Validate request
	if request.GitDiff == "" {
		http.Error(w, `{"error":"Empty git diff"}`, http.StatusBadRequest)
		return
	}

	// Set defaults
	if request.AIProvider == "" {
		request.AIProvider = h.config.DefaultProvider
	}
	if request.AIModel == "" {
		request.AIModel = h.config.DefaultModel
	}
	if request.Language == "" {
		request.Language = "unknown"
	}

	log.Printf("Review request: provider=%s, model=%s, language=%s, diff_size=%d bytes",
		request.AIProvider, request.AIModel, request.Language, len(request.GitDiff))

	// Get provider
	provider, err := h.registry.Get(request.AIProvider)
	if err != nil {
		log.Printf("Provider error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Provider not available: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Call AI provider with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	aiResponse, err := provider.Review(ctx, &request)
	if err != nil {
		log.Printf("AI review error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"AI review failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Build response in reviewdog diagnostic format
	response := models.ReviewResponse{
		Source: models.Source{
			Name: "ai-review",
			URL:  "",
		},
		Diagnostics: aiResponse.Diagnostics,
		Overview:    aiResponse.Overview,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

	log.Printf("Review completed: %d diagnostics found", len(response.Diagnostics))
}

