package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/trading-engine/internal/storage"
)

// TestCalendarHandler_GetCalendar tests the GET /api/calendar endpoint
func TestCalendarHandler_GetCalendar(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewCalendarHandler(db, logger)

	t.Run("Returns 10-week calendar structure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/calendar", nil)
		w := httptest.NewRecorder()

		handler.GetCalendar(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check response structure (API returns {data: T})
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map")
		}

		// Check for weeks array
		weeks, ok := data["weeks"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'weeks' to be an array")
		}

		// Should have 10 weeks (2 back + 8 forward)
		if len(weeks) != 10 {
			t.Errorf("Expected 10 weeks, got %d", len(weeks))
		}

		// Check structure of first week
		if len(weeks) > 0 {
			firstWeek, ok := weeks[0].(map[string]interface{})
			if !ok {
				t.Fatalf("Expected first week to be a map")
			}

			expectedFields := []string{"week_start", "week_end", "sectors"}
			for _, field := range expectedFields {
				if _, exists := firstWeek[field]; !exists {
					t.Errorf("Expected field '%s' in week data", field)
				}
			}
		}

		// Check for sectors array
		sectors, ok := data["sectors"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'sectors' to be an array")
		}

		// Should have at least the common sectors
		if len(sectors) < 8 {
			t.Errorf("Expected at least 8 sectors, got %d", len(sectors))
		}
	})

	t.Run("Method not allowed for non-GET", func(t *testing.T) {
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/api/calendar", nil)
			w := httptest.NewRecorder()

			handler.GetCalendar(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, w.Code)
			}
		}
	})
}

// TestCalendarHandler_GetCalendar_DatabaseError tests error handling
func TestCalendarHandler_GetCalendar_DatabaseError(t *testing.T) {
	// Create handler with closed database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewCalendarHandler(db, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar", nil)
	w := httptest.NewRecorder()

	handler.GetCalendar(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
