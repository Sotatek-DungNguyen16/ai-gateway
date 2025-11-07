package prompt

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sotatek-DungNguyen16/ai-review-gateway/internal/models"
)

// GenerateSystemPrompt creates the system prompt for the AI
func GenerateSystemPrompt(language string) string {
	return fmt.Sprintf(`You are an expert code reviewer specializing in %s. Review ALL code changes and provide comprehensive feedback on these specific categories:

## Review Categories (Check ALL for every request):

1. **Possible Bug** - Logic errors, null pointer risks, off-by-one errors, race conditions, edge cases not handled
2. **Best Practice** - Coding standards violations, naming conventions, code organization, design patterns misuse
3. **Performance** - Inefficient algorithms, unnecessary loops, memory leaks, N+1 queries, blocking operations
4. **Maintainability** - Code complexity, lack of documentation, unclear variable names, hard-coded values, tight coupling
5. **Possible Issue** - Code smells, anti-patterns, deprecated API usage, potential future problems
6. **Enhancement** - Optimization opportunities, better approaches, missing features, code improvements

## Output Format
You must respond ONLY with valid JSON in this exact format:

{
  "overview": "Brief summary covering findings across all 6 categories (2-4 sentences)",
  "issues": [
    {
      "file": "path/to/file.ext",
      "line": 42,
      "column": 10,
      "severity": "ERROR|WARNING|INFO",
      "category": "possible-bug|best-practice|performance|maintainability|possible-issue|enhancement",
      "message": "Clear description with category context",
      "suggestion": "Specific actionable fix"
    }
  ]
}

## Severity Guidelines:
- **ERROR**: Definite bugs, security vulnerabilities, critical performance issues
- **WARNING**: Maintainability concerns, performance bottlenecks, likely bugs, anti-patterns
- **INFO**: Best practice suggestions, enhancements, minor optimizations

## Important Rules:
- Review EVERY changed line against ALL 6 categories
- Provide specific line numbers and actionable suggestions
- Include code examples in suggestions when helpful
- If no issues found, still acknowledge what was reviewed well
- Focus on changed code (marked with + or -)
- Be thorough but constructive
- Prioritize issues by severity and impact
- Consider %s-specific best practices and idioms`, language, language)
}

// GenerateUserPrompt creates the user prompt with the git diff
func GenerateUserPrompt(request *models.ReviewRequest) string {
	var builder strings.Builder

	builder.WriteString("Please review the following code changes:\n\n")

	// Add metadata if available
	if request.GitInfo != nil {
		builder.WriteString("**Context:**\n")
		if request.GitInfo.RepoURL != "" {
			builder.WriteString(fmt.Sprintf("- Repository: %s\n", request.GitInfo.RepoURL))
		}
		if request.GitInfo.BranchName != "" {
			builder.WriteString(fmt.Sprintf("- Branch: %s\n", request.GitInfo.BranchName))
		}
		if request.GitInfo.PRNumber != "" {
			builder.WriteString(fmt.Sprintf("- PR Number: #%s\n", request.GitInfo.PRNumber))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("**Git Diff:**\n```diff\n")
	builder.WriteString(request.GitDiff)
	builder.WriteString("\n```\n\n")

	builder.WriteString("**Review Instructions:**\n")
	builder.WriteString("1. Check EVERY changed line against ALL 6 categories:\n")
	builder.WriteString("   - Possible Bug\n")
	builder.WriteString("   - Best Practice\n")
	builder.WriteString("   - Performance\n")
	builder.WriteString("   - Maintainability\n")
	builder.WriteString("   - Possible Issue\n")
	builder.WriteString("   - Enhancement\n\n")
	builder.WriteString("2. Provide specific line numbers and actionable suggestions\n")
	builder.WriteString("3. Respond ONLY with valid JSON in the format specified\n")

	return builder.String()
}

// ParseAIResponse parses the AI response into structured diagnostics
func ParseAIResponse(responseText string) (*models.AIProviderResponse, error) {
	// Try to extract JSON from code blocks if present
	jsonStr := extractJSON(responseText)

	// Parse JSON
	var rawResponse struct {
		Overview string `json:"overview"`
		Issues   []struct {
			File       string `json:"file"`
			Line       int    `json:"line"`
			Column     int    `json:"column,omitempty"`
			Severity   string `json:"severity"`
			Category   string `json:"category"`
			Message    string `json:"message"`
			Suggestion string `json:"suggestion,omitempty"`
		} `json:"issues"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &rawResponse); err != nil {
		// If JSON parsing fails, try to extract issues from text
		return parseUnstructuredResponse(responseText)
	}

	// Convert to diagnostics
	diagnostics := make([]models.Diagnostic, 0, len(rawResponse.Issues))

	for _, issue := range rawResponse.Issues {
		// Normalize severity
		severity := normalizeSeverity(issue.Severity)

		// Set default column if not provided
		column := issue.Column
		if column == 0 {
			column = 1
		}

		diagnostic := models.Diagnostic{
			Message: issue.Message,
			Location: models.Location{
				Path: issue.File,
				Range: models.Range{
					Start: models.Position{
						Line:   issue.Line,
						Column: column,
					},
					End: models.Position{
						Line:   issue.Line,
						Column: column + 1,
					},
				},
			},
			Severity: severity,
			Code: models.Code{
				Value: issue.Category,
				URL:   "",
			},
			Suggestion: issue.Suggestion,
		}

		diagnostics = append(diagnostics, diagnostic)
	}

	return &models.AIProviderResponse{
		Overview:    rawResponse.Overview,
		Diagnostics: diagnostics,
	}, nil
}

// extractJSON tries to extract JSON from markdown code blocks or raw text
func extractJSON(text string) string {
	// Try to find JSON in code blocks
	codeBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*({.*?})\\s*```")
	matches := codeBlockRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try to find raw JSON
	jsonRegex := regexp.MustCompile(`(?s)({[\s\S]*"overview"[\s\S]*"issues"[\s\S]*})`)
	matches = jsonRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}

	// Return as is and hope it's valid JSON
	return text
}

// parseUnstructuredResponse tries to parse a response that isn't proper JSON
func parseUnstructuredResponse(text string) (*models.AIProviderResponse, error) {
	// Look for common patterns in unstructured text
	diagnostics := make([]models.Diagnostic, 0)

	// Pattern: File: path/to/file.js Line: 42 - Message
	fileLineRegex := regexp.MustCompile(`(?i)(?:file|path):\s*([^\s]+).*?(?:line|L):\s*(\d+).*?[-:]\s*(.+?)(?:\n|$)`)
	matches := fileLineRegex.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			line, _ := strconv.Atoi(match[2])
			diagnostic := models.Diagnostic{
				Message: strings.TrimSpace(match[3]),
				Location: models.Location{
					Path: match[1],
					Range: models.Range{
						Start: models.Position{Line: line, Column: 1},
						End:   models.Position{Line: line, Column: 2},
					},
				},
				Severity: "INFO",
				Code: models.Code{
					Value: "ai-review",
					URL:   "",
				},
			}
			diagnostics = append(diagnostics, diagnostic)
		}
	}

	overview := "AI code review completed"
	if len(diagnostics) == 0 {
		// Extract first paragraph as overview
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				overview = line
				break
			}
		}
	}

	return &models.AIProviderResponse{
		Overview:    overview,
		Diagnostics: diagnostics,
	}, nil
}

// normalizeSeverity normalizes severity values to standard levels
func normalizeSeverity(severity string) string {
	severity = strings.ToUpper(strings.TrimSpace(severity))

	switch severity {
	case "ERROR", "CRITICAL", "HIGH":
		return "ERROR"
	case "WARNING", "WARN", "MEDIUM":
		return "WARNING"
	case "INFO", "INFORMATION", "LOW", "NOTE":
		return "INFO"
	default:
		return "INFO"
	}
}
