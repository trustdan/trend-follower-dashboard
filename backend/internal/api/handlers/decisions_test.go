package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/trading-engine/internal/storage"
)

// TestDecisionsHandler_SaveDecision tests the POST /api/decisions endpoint
func TestDecisionsHandler_SaveDecision(t *testing.T) {
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
	handler := NewDecisionHandler(db, logger)

	tests := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		expectedStatus int
		expectData     bool
	}{
		{
			name:   "Valid GO decision with stock",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":            "AAPL",
				"decision":          "GO",
				"entry":             180.0,
				"atr":               1.5,
				"shares":            250,
				"method":            "stock",
				"sector":            "Tech/Comm",
				"banner_status":     "GREEN",
				"risk_dollars":      750.0,
				"notes":             "Strong breakout",
				"banner_green":      true,
				"timer_complete":    true,
				"not_on_cooldown":   true,
				"heat_passed":       true,
				"sizing_complete":   true,
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:   "Valid NO-GO decision",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":            "TSLA",
				"decision":          "NO-GO",
				"notes":             "Heat cap exceeded",
				"banner_green":      false,
				"timer_complete":    false,
				"not_on_cooldown":   true,
				"heat_passed":       false,
				"sizing_complete":   false,
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:   "Missing ticker",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"decision": "GO",
				"notes":    "Missing ticker",
			},
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
		{
			name:   "Invalid decision",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":   "AAPL",
				"decision": "MAYBE", // Invalid
				"notes":    "Invalid decision",
			},
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
		{
			name:   "GO without entry price",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":            "AAPL",
				"decision":          "GO",
				"atr":               1.5,
				"notes":             "Missing entry",
				"banner_green":      true,
				"timer_complete":    true,
				"not_on_cooldown":   true,
				"heat_passed":       true,
				"sizing_complete":   true,
			},
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
		{
			name:           "Method not allowed (GET)",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectData:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.requestBody != nil {
				bodyBytes, _ := json.Marshal(tt.requestBody)
				reqBody = bytes.NewBuffer(bodyBytes)
			} else {
				reqBody = bytes.NewBuffer([]byte{})
			}

			req := httptest.NewRequest(tt.method, "/api/decisions", reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.SaveDecision(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectData {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Check response structure (API returns {data: T})
				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected data to be a map")
				}

				// Check for decision ID
				if _, exists := data["id"]; !exists {
					t.Errorf("Expected 'id' in response")
				}
			}
		})
	}
}

