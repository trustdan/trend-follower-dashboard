package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRecovery tests the Recovery middleware
func TestRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handlerPanics  bool
		panicValue     interface{}
		expectedStatus int
		expectLog      bool
	}{
		{
			name:           "Normal request without panic",
			handlerPanics:  false,
			panicValue:     nil,
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name:           "Handler panics with string",
			handlerPanics:  true,
			panicValue:     "something went wrong",
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
		},
		{
			name:           "Handler panics with error",
			handlerPanics:  true,
			panicValue:     http.ErrAbortHandler,
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
		},
		{
			name:           "Handler panics with nil",
			handlerPanics:  true,
			panicValue:     nil,
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture log output
			var logBuffer bytes.Buffer
			logger := log.New(&logBuffer, "", 0)

			// Create test handler that may panic
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.handlerPanics {
					panic(tt.panicValue)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			// Wrap with recovery middleware
			handler := Recovery(logger)(testHandler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			// Execute request (should not panic even if handler panics)
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Recovery middleware did not catch panic: %v", r)
					}
				}()
				handler.ServeHTTP(w, req)
			}()

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check log output
			logOutput := logBuffer.String()
			if tt.expectLog {
				if !strings.Contains(logOutput, "PANIC") {
					t.Error("Expected 'PANIC' in log output")
				}
			} else {
				if strings.Contains(logOutput, "PANIC") {
					t.Error("Did not expect 'PANIC' in log output for successful request")
				}
			}

			// Check response body
			if tt.handlerPanics {
				responseBody := w.Body.String()
				if !strings.Contains(responseBody, "Internal Server Error") {
					t.Errorf("Expected error response, got: %s", responseBody)
				}
			} else {
				if w.Body.String() != "success" {
					t.Errorf("Expected 'success' response, got: %s", w.Body.String())
				}
			}
		})
	}
}

// TestRecovery_StackTrace tests that the recovery middleware logs stack traces
func TestRecovery_StackTrace(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler := Recovery(logger)(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	logOutput := logBuffer.String()

	// Should contain panic message
	if !strings.Contains(logOutput, "test panic") {
		t.Error("Expected panic message in log output")
	}

	// Should contain PANIC indicator
	if !strings.Contains(logOutput, "PANIC") {
		t.Error("Expected 'PANIC' indicator in log output")
	}

	// Should contain stack trace information (goroutine info)
	if !strings.Contains(logOutput, "goroutine") {
		t.Error("Expected stack trace information in log output")
	}
}

// TestRecovery_MultipleRequests tests that recovery middleware works correctly across multiple requests
func TestRecovery_MultipleRequests(t *testing.T) {
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)

	callCount := 0
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 2 {
			panic("second request panics")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	handler := Recovery(logger)(testHandler)

	// First request should succeed
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request: expected status %d, got %d", http.StatusOK, w1.Code)
	}

	// Second request should panic but be recovered
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusInternalServerError {
		t.Errorf("Second request: expected status %d, got %d", http.StatusInternalServerError, w2.Code)
	}

	// Third request should succeed again
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Third request: expected status %d, got %d", http.StatusOK, w3.Code)
	}

	// Should have exactly one panic logged
	logOutput := logBuffer.String()
	panicCount := strings.Count(logOutput, "PANIC")
	if panicCount != 1 {
		t.Errorf("Expected 1 PANIC log entry, got %d", panicCount)
	}
}
