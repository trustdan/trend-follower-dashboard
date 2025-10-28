package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSetting_ValidEquity(t *testing.T) {
	err := ValidateSetting("Equity_E", "10000")
	assert.NoError(t, err)

	err = ValidateSetting("Equity_E", "50000.50")
	assert.NoError(t, err)
}

func TestValidateSetting_InvalidEquity(t *testing.T) {
	err := ValidateSetting("Equity_E", "0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Equity_E must be positive")

	err = ValidateSetting("Equity_E", "-1000")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Equity_E must be positive")
}

func TestValidateSetting_ValidRiskPct(t *testing.T) {
	err := ValidateSetting("RiskPct_r", "0.0075")
	assert.NoError(t, err)

	err = ValidateSetting("RiskPct_r", "0.01")
	assert.NoError(t, err)

	err = ValidateSetting("RiskPct_r", "1")
	assert.NoError(t, err)
}

func TestValidateSetting_InvalidRiskPct(t *testing.T) {
	err := ValidateSetting("RiskPct_r", "0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RiskPct_r must be between 0 and 1")

	err = ValidateSetting("RiskPct_r", "1.5")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RiskPct_r must be between 0 and 1")

	err = ValidateSetting("RiskPct_r", "-0.01")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RiskPct_r must be between 0 and 1")
}

func TestValidateSetting_ValidHeatCap(t *testing.T) {
	err := ValidateSetting("HeatCap_H_pct", "0.04")
	assert.NoError(t, err)

	err = ValidateSetting("HeatCap_H_pct", "0.05")
	assert.NoError(t, err)
}

func TestValidateSetting_InvalidHeatCap(t *testing.T) {
	err := ValidateSetting("HeatCap_H_pct", "1.5")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HeatCap_H_pct must be between 0 and 1")
}

func TestValidateSetting_ValidBucketHeatCap(t *testing.T) {
	err := ValidateSetting("BucketHeatCap_pct", "0.015")
	assert.NoError(t, err)

	err = ValidateSetting("BucketHeatCap_pct", "0.02")
	assert.NoError(t, err)
}

func TestValidateSetting_InvalidBucketHeatCap(t *testing.T) {
	err := ValidateSetting("BucketHeatCap_pct", "1.5")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BucketHeatCap_pct must be between 0 and 1")
}

func TestValidateSetting_ValidStopMultiple(t *testing.T) {
	err := ValidateSetting("StopMultiple_K", "2")
	assert.NoError(t, err)

	err = ValidateSetting("StopMultiple_K", "3")
	assert.NoError(t, err)

	err = ValidateSetting("StopMultiple_K", "2.5")
	assert.NoError(t, err)
}

func TestValidateSetting_InvalidStopMultiple(t *testing.T) {
	err := ValidateSetting("StopMultiple_K", "0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "StopMultiple_K must be positive")

	err = ValidateSetting("StopMultiple_K", "-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "StopMultiple_K must be positive")
}

func TestValidateSetting_InvalidKey(t *testing.T) {
	err := ValidateSetting("InvalidKey", "100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "setting not found: InvalidKey")
}

func TestValidateSetting_NonNumericValue(t *testing.T) {
	err := ValidateSetting("Equity_E", "abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be a number")

	err = ValidateSetting("RiskPct_r", "not-a-number")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value must be a number")
}

func TestValidateSetting_AllValidKeys(t *testing.T) {
	// Verify all valid keys are recognized
	validSettings := map[string]string{
		"Equity_E":           "10000",
		"RiskPct_r":          "0.0075",
		"HeatCap_H_pct":      "0.04",
		"BucketHeatCap_pct":  "0.015",
		"StopMultiple_K":     "2",
	}

	for key, value := range validSettings {
		err := ValidateSetting(key, value)
		assert.NoError(t, err, "Key %s should be valid", key)
	}
}
