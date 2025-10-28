package logx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCorrelationID(t *testing.T) {
	corrID1 := GenerateCorrelationID()
	corrID2 := GenerateCorrelationID()

	assert.NotEmpty(t, corrID1, "Correlation ID should not be empty")
	assert.NotEmpty(t, corrID2, "Correlation ID should not be empty")
	assert.NotEqual(t, corrID1, corrID2, "Correlation IDs should be unique")
	assert.Len(t, corrID1, 36, "UUID should be 36 characters")
}

func TestInitialize(t *testing.T) {
	logPath := "test_logger.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err, "Initialize should not return error")

	// Verify log file was created
	_, err = os.Stat(logPath)
	require.NoError(t, err, "Log file should exist")

	// Verify logger is set
	assert.NotNil(t, logger, "Logger should be initialized")
}

func TestInitializeWithDirectory(t *testing.T) {
	logPath := "test_logs/nested/test.log"
	defer os.RemoveAll("test_logs")

	err := Initialize(logPath)
	require.NoError(t, err, "Initialize should create directories")

	// Verify log file was created
	_, err = os.Stat(logPath)
	require.NoError(t, err, "Log file should exist")
}

func TestWithCorrelationID(t *testing.T) {
	logPath := "test_corr_id.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err)

	corrID := "test-corr-123"
	entry := WithCorrelationID(corrID)

	require.NotNil(t, entry, "Entry should not be nil")
	assert.Equal(t, corrID, entry.Data["corr_id"], "Correlation ID should match")
}

func TestWithCorrelationIDFallback(t *testing.T) {
	// Reset logger to test fallback
	logger = nil

	corrID := "fallback-test-456"
	entry := WithCorrelationID(corrID)

	require.NotNil(t, entry, "Entry should not be nil even without Initialize")
	assert.Equal(t, corrID, entry.Data["corr_id"], "Correlation ID should match")
}

func TestInfoLogging(t *testing.T) {
	logPath := "test_info.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err)

	// Should not panic
	Info("Test info message", map[string]interface{}{
		"test_field": "test_value",
	})

	// Verify log file has content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test info message")
	assert.Contains(t, string(content), "test_field")
}

func TestErrorLogging(t *testing.T) {
	logPath := "test_error.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err)

	testErr := assert.AnError

	// Should not panic
	Error("Test error message", testErr, map[string]interface{}{
		"error_field": "error_value",
	})

	// Verify log file has content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test error message")
	assert.Contains(t, string(content), "error_field")
	assert.Contains(t, string(content), "error")
}

func TestDebugLogging(t *testing.T) {
	logPath := "test_debug.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err)

	// Set DEBUG level
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	// Reinitialize with debug level
	err = Initialize(logPath)
	require.NoError(t, err)

	// Should not panic
	Debug("Test debug message", map[string]interface{}{
		"debug_field": "debug_value",
	})

	// Verify log file has content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test debug message")
}

func TestWarnLogging(t *testing.T) {
	logPath := "test_warn.log"
	defer os.Remove(logPath)

	err := Initialize(logPath)
	require.NoError(t, err)

	// Should not panic
	Warn("Test warn message", map[string]interface{}{
		"warn_field": "warn_value",
	})

	// Verify log file has content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test warn message")
	assert.Contains(t, string(content), "warn_field")
}
