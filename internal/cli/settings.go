package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/domain"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewGetSettingsCommand creates the get-settings command
func NewGetSettingsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-settings",
		Short: "Get configuration settings",
		Long: `Get all configuration settings or a specific setting by key.

Examples:
  # Get all settings
  tf-engine get-settings

  # Get a specific setting
  tf-engine get-settings --key Equity_E`,
		RunE: runGetSettings,
	}

	cmd.Flags().String("key", "", "Specific setting key to retrieve")

	return cmd
}

func runGetSettings(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	log.Info("Getting settings")

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	key, _ := cmd.Flags().GetString("key")

	if key != "" {
		// Get single setting
		value, err := db.GetSetting(key)
		if err != nil {
			log.WithError(err).WithField("key", key).Error("Failed to get setting")
			return fmt.Errorf("failed to get setting: %w", err)
		}

		result := map[string]string{key: value}
		jsonResult, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonResult))

		log.WithField("key", key).WithField("value", value).Info("Retrieved setting")
	} else {
		// Get all settings
		settings, err := db.GetAllSettings()
		if err != nil {
			log.WithError(err).Error("Failed to get settings")
			return fmt.Errorf("failed to get settings: %w", err)
		}

		jsonResult, _ := json.MarshalIndent(settings, "", "  ")
		fmt.Println(string(jsonResult))

		log.WithField("count", len(settings)).Info("Retrieved all settings")
	}

	return nil
}

// NewSetSettingCommand creates the set-setting command
func NewSetSettingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-setting",
		Short: "Update a configuration setting",
		Long: `Update a configuration setting with validation.

Valid settings:
  - Equity_E: Account equity (must be positive)
  - RiskPct_r: Risk per trade as decimal (must be between 0 and 1)
  - HeatCap_H_pct: Portfolio heat cap as decimal (must be between 0 and 1)
  - BucketHeatCap_pct: Bucket heat cap as decimal (must be between 0 and 1)
  - StopMultiple_K: ATR stop multiple (must be positive)

Examples:
  # Update account equity
  tf-engine set-setting --key Equity_E --value 20000

  # Update risk percent
  tf-engine set-setting --key RiskPct_r --value 0.01`,
		RunE: runSetSetting,
	}

	cmd.Flags().String("key", "", "Setting key (required)")
	cmd.Flags().String("value", "", "New value (required)")

	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")

	return cmd
}

func runSetSetting(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	key, _ := cmd.Flags().GetString("key")
	value, _ := cmd.Flags().GetString("value")

	log.WithField("key", key).WithField("value", value).Info("Setting configuration value")

	// Validate setting
	if err := domain.ValidateSetting(key, value); err != nil {
		log.WithError(err).Error("Setting validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}

	// Open database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Update setting
	if err := db.SetSetting(key, value); err != nil {
		log.WithError(err).Error("Failed to update setting")
		return fmt.Errorf("failed to update setting: %w", err)
	}

	log.Info("Setting updated successfully")

	fmt.Printf("✓ Setting updated: %s = %s\n", key, value)

	return nil
}
