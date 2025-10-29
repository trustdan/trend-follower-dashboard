package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/yourusername/trading-engine/internal/api/responses"
	"github.com/yourusername/trading-engine/internal/storage"
)

// CalendarHandler handles calendar-related API requests
type CalendarHandler struct {
	db     *storage.DB
	logger *log.Logger
}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler(db *storage.DB, logger *log.Logger) *CalendarHandler {
	return &CalendarHandler{
		db:     db,
		logger: logger,
	}
}

// WeekData represents positions for a specific week
type WeekData struct {
	WeekStart string               `json:"week_start"` // Monday of the week in YYYY-MM-DD format
	WeekEnd   string               `json:"week_end"`   // Sunday of the week in YYYY-MM-DD format
	Sectors   map[string][]PositionInfo `json:"sectors"`    // sector -> positions in that sector
}

// PositionInfo represents a position in the calendar view
type PositionInfo struct {
	Ticker      string  `json:"ticker"`
	EntryPrice  float64 `json:"entry_price"`
	RiskDollars float64 `json:"risk_dollars"`
	Status      string  `json:"status"` // OPEN or CLOSED
	DaysHeld    int     `json:"days_held"`
}

// CalendarResponse represents the calendar API response
type CalendarResponse struct {
	Weeks   []WeekData `json:"weeks"`   // 10 weeks of data (2 back + 8 forward)
	Sectors []string   `json:"sectors"` // All unique sectors
}

// GetCalendar handles GET /api/calendar
// Returns a rolling 10-week view (2 weeks back + 8 weeks forward)
func (h *CalendarHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Get all open positions
	positions, err := h.db.GetPositions()
	if err != nil {
		h.logger.Printf("Error getting positions: %v", err)
		responses.InternalError(w, err)
		return
	}

	// Calculate 10-week window (2 back + 8 forward)
	now := time.Now()
	startWeek := getMondayOfWeek(now.AddDate(0, 0, -14)) // 2 weeks back

	weeks := make([]WeekData, 10)
	sectorsMap := make(map[string]bool)

	// Build 10 weeks
	for i := 0; i < 10; i++ {
		weekStart := startWeek.AddDate(0, 0, i*7)
		weekEnd := weekStart.AddDate(0, 0, 6) // Sunday

		weeks[i] = WeekData{
			WeekStart: weekStart.Format("2006-01-02"),
			WeekEnd:   weekEnd.Format("2006-01-02"),
			Sectors:   make(map[string][]PositionInfo),
		}

		// Add positions that overlap with this week
		for _, pos := range positions {
			if pos.Status != "OPEN" {
				continue // Only show open positions
			}

			// Check if position overlaps with this week
			posOpenDate := pos.OpenedAt
			if posOpenDate.After(weekEnd) {
				continue // Position opened after this week
			}

			// Position is active during this week
			daysHeld := int(now.Sub(posOpenDate).Hours() / 24)

			posInfo := PositionInfo{
				Ticker:      pos.Ticker,
				EntryPrice:  pos.EntryPrice,
				RiskDollars: pos.RiskDollars,
				Status:      pos.Status,
				DaysHeld:    daysHeld,
			}

			bucket := pos.Bucket
			if bucket == "" {
				bucket = "Other"
			}

			weeks[i].Sectors[bucket] = append(weeks[i].Sectors[bucket], posInfo)
			sectorsMap[bucket] = true
		}
	}

	// Build unique sectors list
	sectors := make([]string, 0, len(sectorsMap))
	for sector := range sectorsMap {
		sectors = append(sectors, sector)
	}

	// Add common sectors even if no positions (for consistent grid)
	commonSectors := []string{"Tech/Comm", "Energy", "Industrial", "Finance", "Healthcare", "Consumer", "Materials", "Other"}
	for _, sector := range commonSectors {
		if !sectorsMap[sector] {
			sectors = append(sectors, sector)
			sectorsMap[sector] = true
		}
	}

	response := CalendarResponse{
		Weeks:   weeks,
		Sectors: sectors,
	}

	responses.Success(w, response)
}

// getMondayOfWeek returns the Monday (start of week) for a given date
func getMondayOfWeek(t time.Time) time.Time {
	// Get the weekday (0 = Sunday, 1 = Monday, ...)
	weekday := int(t.Weekday())

	// Calculate days to subtract to get to Monday
	daysToMonday := (weekday - 1 + 7) % 7
	if weekday == 0 { // Sunday
		daysToMonday = 6
	}

	monday := t.AddDate(0, 0, -daysToMonday)

	// Normalize to midnight
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}
