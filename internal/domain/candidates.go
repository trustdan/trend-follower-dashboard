package domain

import (
	"fmt"
	"strings"
	"time"
)

// Candidate represents a daily screening candidate
type Candidate struct {
	ID        int
	Date      string
	Ticker    string
	PresetID  *int
	Sector    string
	Bucket    string
	CreatedAt time.Time
}

// ImportCandidatesRequest represents a request to import candidates
type ImportCandidatesRequest struct {
	Tickers string
	Preset  string
	Sector  string
	Bucket  string
	Date    string // Optional: defaults to today
}

// ImportCandidatesResult represents the result of importing candidates
type ImportCandidatesResult struct {
	Count      int      `json:"count"`
	Date       string   `json:"date"`
	Tickers    []string `json:"tickers"`
	Preset     string   `json:"preset,omitempty"`
	Sector     string   `json:"sector,omitempty"`
	Bucket     string   `json:"bucket,omitempty"`
	Normalized bool     `json:"normalized"`
}

// NormalizeTickers normalizes a comma-separated list of tickers
//
// Rules:
//   - Trim whitespace from each ticker
//   - Convert to uppercase
//   - Remove empty entries
//   - Return as slice
func NormalizeTickers(tickers string) ([]string, error) {
	if strings.TrimSpace(tickers) == "" {
		return nil, fmt.Errorf("at least one ticker required")
	}

	parts := strings.Split(tickers, ",")
	normalized := make([]string, 0, len(parts))

	for _, ticker := range parts {
		ticker = strings.TrimSpace(ticker)
		ticker = strings.ToUpper(ticker)
		if ticker != "" {
			normalized = append(normalized, ticker)
		}
	}

	if len(normalized) == 0 {
		return nil, fmt.Errorf("at least one ticker required")
	}

	return normalized, nil
}

// ValidateImportRequest validates the import candidates request
func ValidateImportRequest(req ImportCandidatesRequest) error {
	// Validate tickers
	if _, err := NormalizeTickers(req.Tickers); err != nil {
		return err
	}

	// Date is optional but must be valid if provided
	if req.Date != "" {
		_, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
		}
	}

	return nil
}

// GetImportDate returns the date for import (today if not specified)
func GetImportDate(dateStr string) string {
	if dateStr != "" {
		return dateStr
	}
	return time.Now().Format("2006-01-02")
}
