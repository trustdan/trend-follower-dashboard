package main

// getSettingWithDefault safely extracts a setting from the map with a default value
func getSettingWithDefault(settings map[string]string, key, defaultValue string) string {
	if val, exists := settings[key]; exists && val != "" {
		return val
	}
	return defaultValue
}
