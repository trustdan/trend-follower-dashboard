package domain

import (
	"fmt"
	"strings"
	"time"
)

// ChecklistRequest contains checklist evaluation input
type ChecklistRequest struct {
	Ticker        string `json:"ticker"`
	FromPreset    bool   `json:"from_preset"`
	TrendPass     bool   `json:"trend_pass"`
	LiquidityPass bool   `json:"liquidity_pass"`
	TVConfirm     bool   `json:"tv_confirm"`
	EarningsOK    bool   `json:"earnings_ok"`
	JournalOK     bool   `json:"journal_ok"`
}

// ChecklistResult contains checklist evaluation output
type ChecklistResult struct {
	Banner              string    `json:"banner"` // "GREEN", "YELLOW", "RED"
	MissingCount        int       `json:"missing_count"`
	MissingItems        []string  `json:"missing_items"`
	EvaluationTimestamp time.Time `json:"evaluation_timestamp,omitempty"`
	AllowSave           bool      `json:"allow_save"`
}

// Banner colors
const (
	BannerGreen  = "GREEN"
	BannerYellow = "YELLOW"
	BannerRed    = "RED"
)

// ChecklistItemNames defines the required checklist items in order
var ChecklistItemNames = []string{
	"FromPreset",
	"TrendPass",
	"LiquidityPass",
	"TVConfirm",
	"EarningsOK",
	"JournalOK",
}

// EvaluateChecklist determines banner color based on missing items
//
// Checklist Rules (from CLAUDE.md):
//   - 0 missing → GREEN (go) - Start impulse timer, allow save
//   - 1 missing → YELLOW (caution) - Do not allow save
//   - 2+ missing → RED (no-go) - Do not allow save
//
// Only GREEN banner starts the impulse timer and allows eventual save.
// This enforces discipline by preventing impulsive trades with incomplete setups.
func EvaluateChecklist(req ChecklistRequest) (*ChecklistResult, error) {
	// Validate ticker
	if strings.TrimSpace(req.Ticker) == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	// Build checklist map
	checklist := map[string]bool{
		"FromPreset":    req.FromPreset,
		"TrendPass":     req.TrendPass,
		"LiquidityPass": req.LiquidityPass,
		"TVConfirm":     req.TVConfirm,
		"EarningsOK":    req.EarningsOK,
		"JournalOK":     req.JournalOK,
	}

	// Count missing items
	missingItems := []string{}
	for _, itemName := range ChecklistItemNames {
		if !checklist[itemName] {
			missingItems = append(missingItems, itemName)
		}
	}

	missingCount := len(missingItems)

	// Determine banner color and permissions
	var banner string
	var allowSave bool
	var evalTimestamp time.Time

	switch missingCount {
	case 0:
		// Perfect setup - all items checked
		banner = BannerGreen
		allowSave = true
		evalTimestamp = time.Now() // Only record timestamp on GREEN
	case 1:
		// Caution - one item missing
		banner = BannerYellow
		allowSave = false
		// No timestamp recorded
	default: // 2 or more
		// No-go - multiple items missing
		banner = BannerRed
		allowSave = false
		// No timestamp recorded
	}

	return &ChecklistResult{
		Banner:              banner,
		MissingCount:        missingCount,
		MissingItems:        missingItems,
		EvaluationTimestamp: evalTimestamp,
		AllowSave:           allowSave,
	}, nil
}
