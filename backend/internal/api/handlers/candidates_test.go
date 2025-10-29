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

// TestCandidatesHandler_GetCandidates tests the GET /api/candidates endpoint
func TestCandidatesHandler_GetCandidates(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewCandidatesHandler(db, logger)

	t.Run("Returns empty array when no candidates", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/candidates", nil)
		w := httptest.NewRecorder()

		handler.GetCandidates(w, req)

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
			t.Errorf("Expected empty array, got %d candidates", len(data))
		}
	})

	t.Run("Method not allowed for non-GET", func(t *testing.T) {
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/api/candidates", nil)
			w := httptest.NewRecorder()

			handler.GetCandidates(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, w.Code)
			}
		}
	})
}

// TestCandidatesHandler_ImportCandidates tests the POST /api/candidates/import endpoint
func TestCandidatesHandler_ImportCandidates(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewCandidatesHandler(db, logger)

	tests := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Valid import with tickers",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"tickers": []string{"AAPL", "NVDA", "TSLA"},
				"date":    "2025-10-29",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Import with single ticker",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"tickers": []string{"AAPL"},
				"date":    "2025-10-29",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Missing tickers",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"date": "2025-10-29",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty tickers array",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"tickers": []string{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Method not allowed (GET)",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
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

			req := httptest.NewRequest(tt.method, "/api/candidates/import", reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ImportCandidates(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}
