package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatHuman OutputFormat = "human"
	FormatJSON  OutputFormat = "json"
)

// GetOutputFormat retrieves the output format from command flags
func GetOutputFormat(cmd *cobra.Command) OutputFormat {
	format, _ := cmd.Flags().GetString("format")
	if format == "json" {
		return FormatJSON
	}
	return FormatHuman
}

// PrintJSON outputs JSON to stdout (for both human and json formats)
func PrintJSON(data interface{}) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// PrintHuman outputs human-readable text (only for human format)
func PrintHuman(format OutputFormat, text string) {
	if format == FormatHuman {
		fmt.Println(text)
	}
}

// PrintHumanf outputs formatted human-readable text (only for human format)
func PrintHumanf(format OutputFormat, formatStr string, args ...interface{}) {
	if format == FormatHuman {
		fmt.Printf(formatStr, args...)
	}
}
