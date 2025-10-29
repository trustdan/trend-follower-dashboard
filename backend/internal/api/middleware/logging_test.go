package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestLogging tests the Logging middleware
func TestLogging(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		path                string
		correlationIDHeader string
		statusCode          int
		responseDelay       time.Duration
		expectCorrelationID bool
		expectSlowWarning   bool
	}{
		{
			name:                "GET request with correlation ID",
			method:              http.MethodGet,
			path:                "/api/settings",
			correlationIDHeader: "test-correlation-123",
			statusCode:          http.StatusOK,
			responseDelay:       0,
			expectCorrelationID: true,
			expectSlowWarning:   false,
		},
		{
			name:                "POST request generates correlation ID",
			method:              http.MethodPost,
			path:                "/api/decisions",
			correlationIDHeader: "",
			statusCode:          http.StatusCreated,
			responseDelay:       0,
			expectCorrelationID: true,
			expectSlowWarning:   false,
		},
		{
			name:                "Slow request triggers warning",
			method:              http.MethodGet,
			path:                "/api/positions",
			correlationIDHeader: "",
			statusCode:          http.StatusOK,
			responseDelay:       600 * time.Millisecond,
			expectCorrelationID: true,
			expectSlowWarning:   true,
		},
		{
			name:                "Error response logged correctly",
			method:              http.MethodPost,
			path:                "/api/sizing/calculate",
			correlationIDHeader: "",
			statusCode:          http.StatusBadRequest,
			responseDelay:       0,
			expectCorrelationID: true,
			expectSlowWarning:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			var logBuffer bytes.Buffer
			logger := log.New(&logBuffer, "", 0)

			// Create test handler that returns the specified status code
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.responseDelay > 0 {
					time.Sleep(tt.responseDelay)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("test response"))
			})

			// Wrap with logging middleware
			handler := Logging(logger)(testHandler)

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.correlationIDHeader != "" {
				req.Header.Set("X-Correlation-ID", tt.correlationIDHeader)
			}
			w := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(w, req)

			// Check correlation ID in response header
			responseCorrelationID := w.Header().Get("X-Correlation-ID")
			if tt.expectCorrelationID {
				if responseCorrelationID == "" {
					t.Error("Expected X-Correlation-ID header in response")
				}

				// If we provided a correlation ID, it should match
				if tt.correlationIDHeader != "" && responseCorrelationID != tt.correlationIDHeader {
					t.Errorf("Expected correlation ID '%s', got '%s'", tt.correlationIDHeader, responseCorrelationID)
				}
			}

			// Check log output
			logOutput := logBuffer.String()

			// Should log the request
			if !strings.Contains(logOutput, tt.method) {
				t.Errorf("Expected log to contain method '%s'", tt.method)
			}

			if !strings.Contains(logOutput, tt.path) {
				t.Errorf("Expected log to contain path '%s'", tt.path)
			}

			// Should log the response with status code
			statusStr := fmt.Sprintf("%d", tt.statusCode)
			if !strings.Contains(logOutput, statusStr) {
				t.Errorf("Expected log to contain status code %d", tt.statusCode)
			}

			// Check for slow request warning
			if tt.expectSlowWarning {
				if !strings.Contains(logOutput, "SLOW REQUEST") {
					t.Error("Expected 'SLOW REQUEST' warning in log output")
				}
			} else {
				if strings.Contains(logOutput, "SLOW REQUEST") {
					t.Error("Did not expect 'SLOW REQUEST' warning in log output")
				}
			}

			// Should have correlation ID in log
			if responseCorrelationID != "" && !strings.Contains(logOutput, responseCorrelationID) {
				t.Errorf("Expected correlation ID '%s' in log output", responseCorrelationID)
			}
		})
	}
}

// TestResponseWriter_WriteHeader tests that the responseWriter correctly captures status codes
func TestResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"200 OK", http.StatusOK},
		{"201 Created", http.StatusCreated},
		{"400 Bad Request", http.StatusBadRequest},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			rw.WriteHeader(tt.statusCode)

			if rw.statusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, rw.statusCode)
			}

			if w.Code != tt.statusCode {
				t.Errorf("Expected underlying ResponseWriter code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

// TestLogging_CorrelationIDGeneration tests that correlation IDs are generated when not provided
func TestLogging_CorrelationIDGeneration(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Logging(logger)(testHandler)

	// Make multiple requests without correlation IDs
	correlationIDs := make(map[string]bool)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		correlationID := w.Header().Get("X-Correlation-ID")
		if correlationID == "" {
			t.Error("Expected correlation ID to be generated")
		}

		// Check for duplicates
		if correlationIDs[correlationID] {
			t.Errorf("Duplicate correlation ID generated: %s", correlationID)
		}
		correlationIDs[correlationID] = true
	}

	// Should have 5 unique correlation IDs
	if len(correlationIDs) != 5 {
		t.Errorf("Expected 5 unique correlation IDs, got %d", len(correlationIDs))
	}
}
