package models

// ReviewRequest represents the incoming review request
type ReviewRequest struct {
	AIModel      string   `json:"ai_model"`
	AIProvider   string   `json:"ai_provider"`
	Language     string   `json:"language"`
	ReviewMode   string   `json:"review_mode"`
	GitDiff      string   `json:"git_diff"`
	GitInfo      *GitInfo `json:"git_info,omitempty"`
}

// GitInfo contains git repository information
type GitInfo struct {
	CommitHash string      `json:"commit_hash"`
	BranchName string      `json:"branch_name"`
	PRNumber   string      `json:"pr_number"`
	RepoURL    string      `json:"repo_url"`
	Author     *GitUser    `json:"author,omitempty"`
	Committer  *GitUser    `json:"committer,omitempty"`
}

// GitUser represents a git user
type GitUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ReviewResponse represents the diagnostic format response
type ReviewResponse struct {
	Source      Source       `json:"source"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Overview    string       `json:"overview,omitempty"`
}

// Source represents the source of diagnostics
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Diagnostic represents a single code review issue
type Diagnostic struct {
	Message  string   `json:"message"`
	Location Location `json:"location"`
	Severity string   `json:"severity"` // ERROR, WARNING, INFO
	Code     Code     `json:"code"`
	Original string   `json:"original,omitempty"` // Original code snippet
	Suggestion string `json:"suggestion,omitempty"` // Suggested fix
}

// Location represents the location of an issue in the code
type Location struct {
	Path  string `json:"path"`
	Range Range  `json:"range"`
}

// Range represents a range in the code
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in a file
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Code represents an issue code
type Code struct {
	Value string `json:"value"`
	URL   string `json:"url"`
}

// AIProviderResponse represents the raw response from AI providers
type AIProviderResponse struct {
	Overview    string
	Diagnostics []Diagnostic
}

