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

// TestPositionsHandler_GetPositions tests the GET /api/positions endpoint
func TestPositionsHandler_GetPositions(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Initialize database
	if err := db.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewPositionsHandler(db, logger)

	t.Run("Empty positions returns empty array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/positions", nil)
		w := httptest.NewRecorder()

		handler.GetPositions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check response structure (API returns {data: T})
		data, ok := response["data"].([]interface{})
		if !ok {
			t.Fatalf("Expected data to be an array")
		}

		if len(data) != 0 {
			t.Errorf("Expected empty array, got %d positions", len(data))
		}
	})

	t.Run("Returns positions after adding one", func(t *testing.T) {
		// Open a position (this will be in "open" status)
		_, err := db.OpenPosition("AAPL")
		if err != nil {
			t.Fatalf("Failed to open test position: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/positions", nil)
		w := httptest.NewRecorder()

		handler.GetPositions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		data, ok := response["data"].([]interface{})
		if !ok {
			t.Fatalf("Expected data to be an array")
		}

		if len(data) != 1 {
			t.Errorf("Expected 1 position, got %d", len(data))
		}

		// Verify we got a position with basic fields
		pos, ok := data[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected position to be a map")
		}

		// Check for essential fields
		if _, exists := pos["ticker"]; !exists {
			t.Error("Expected 'ticker' field in position")
		}
		if _, exists := pos["status"]; !exists {
			t.Error("Expected 'status' field in position")
		}
	})

	t.Run("Method not allowed for non-GET", func(t *testing.T) {
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/api/positions", nil)
			w := httptest.NewRecorder()

			handler.GetPositions(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, w.Code)
			}
		}
	})
}

// TestPositionsHandler_GetPositions_DatabaseError tests error handling
func TestPositionsHandler_GetPositions_DatabaseError(t *testing.T) {
	// Create handler with closed database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewPositionsHandler(db, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/positions", nil)
	w := httptest.NewRecorder()

	handler.GetPositions(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
