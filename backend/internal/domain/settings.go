package domain

import (
	"fmt"
	"strconv"
)

// SettingKey represents a valid setting key
type SettingKey string

const (
	SettingEquity        SettingKey = "Equity_E"
	SettingRiskPct       SettingKey = "RiskPct_r"
	SettingHeatCap       SettingKey = "HeatCap_H_pct"
	SettingBucketHeatCap SettingKey = "BucketHeatCap_pct"
	SettingStopMultiple  SettingKey = "StopMultiple_K"
)

// ValidSettingKeys lists all valid setting keys
var ValidSettingKeys = []SettingKey{
	SettingEquity,
	SettingRiskPct,
	SettingHeatCap,
	SettingBucketHeatCap,
	SettingStopMultiple,
}

// ValidateSetting validates a setting key and value
//
// Validation Rules:
//   - Key must be one of the 5 valid settings
//   - Value must be numeric
//   - Equity_E must be positive
//   - RiskPct_r must be between 0 and 1
//   - HeatCap_H_pct must be between 0 and 1
//   - BucketHeatCap_pct must be between 0 and 1
//   - StopMultiple_K must be positive
func ValidateSetting(key, value string) error {
	// Check if key is valid
	valid := false
	for _, validKey := range ValidSettingKeys {
		if string(validKey) == key {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("setting not found: %s", key)
	}

	// Parse value as float
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("value must be a number, got: %s", value)
	}

	// Key-specific validation
	switch SettingKey(key) {
	case SettingEquity:
		if floatVal <= 0 {
			return fmt.Errorf("Equity_E must be positive, got %.2f", floatVal)
		}

	case SettingRiskPct:
		if floatVal <= 0 || floatVal > 1 {
			return fmt.Errorf("RiskPct_r must be between 0 and 1, got %.4f", floatVal)
		}

	case SettingHeatCap:
		if floatVal <= 0 || floatVal > 1 {
			return fmt.Errorf("HeatCap_H_pct must be between 0 and 1, got %.4f", floatVal)
		}

	case SettingBucketHeatCap:
		if floatVal <= 0 || floatVal > 1 {
			return fmt.Errorf("BucketHeatCap_pct must be between 0 and 1, got %.4f", floatVal)
		}

	case SettingStopMultiple:
		if floatVal <= 0 {
			return fmt.Errorf("StopMultiple_K must be positive, got %.2f", floatVal)
		}
	}

	return nil
}
