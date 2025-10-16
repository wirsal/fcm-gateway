package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig tests the LoadConfig function.
func TestLoadConfig(t *testing.T) {
	// t.Cleanup ensures viper's global state is reset after each test function (TestLoadConfig) completes.
	t.Cleanup(viper.Reset)

	t.Run("success - should load config from yaml file", func(t *testing.T) {
		// --- Setup ---
		viper.Reset() // Reset viper for this specific sub-test to ensure total isolation.
		tempDir := t.TempDir()

		// Define the content of the test config file with correct YAML indentation.
		// 'port' must be indented under 'server'.
		// 'credentials_file', 'scopes', and 'endpoint_url' must be indented under 'fcm'.
		configContent := `
server:
  port: "8081"
fcm:
  credentials_file: "test-credentials.json"
  scopes:
    - "https://www.googleapis.com/auth/firebase.messaging"
  endpoint_url: "http://localhost:8080/fcm/send"
`
		configPath := filepath.Join(tempDir, ".config.yaml")
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		require.NoError(t, err, "Failed to write temporary config file")

		// --- Execute ---
		cfg, err := LoadConfig(tempDir)

		// --- Assert ---
		require.NoError(t, err)
		require.NotNil(t, cfg, "Config should not be nil")

		// Assert server configuration.
		assert.Equal(t, "8081", cfg.Server.Port)

		// Assert FCM configuration.
		assert.Equal(t, "test-credentials.json", cfg.FCM.CredentialsFile)
		assert.Equal(t, []string{"https://www.googleapis.com/auth/firebase.messaging"}, cfg.FCM.Scopes)
		assert.Equal(t, "http://localhost:8080/fcm/send", cfg.FCM.EndpointURL)
	})

	t.Run("error - config file not found", func(t *testing.T) {
		// --- Setup ---
		viper.Reset() // Reset for isolation.
		tempDir := t.TempDir()

		// --- Execute ---
		cfg, err := LoadConfig(tempDir)

		// --- Assert ---
		require.Error(t, err, "Expected an error when config file is not found")
		assert.Nil(t, cfg, "Config should be nil on error")
	})

	t.Run("error - unmarshal error", func(t *testing.T) {
		// --- Setup ---
		viper.Reset() // Reset for isolation.
		tempDir := t.TempDir()

		// Create a config file with a type mismatch (port is an array, not a string).
		// Note the correct indentation.
		invalidConfigContent := `
server:
  port: ["8081"]
`
		configPath := filepath.Join(tempDir, ".config.yaml")
		err := os.WriteFile(configPath, []byte(invalidConfigContent), 0644)
		require.NoError(t, err)

		// --- Execute ---
		cfg, err := LoadConfig(tempDir)

		// --- Assert ---
		require.Error(t, err, "Expected an error during unmarshalling")
		assert.Contains(t, err.Error(), "decoding failed", "Error message should indicate a decoding problem")
		assert.Nil(t, cfg, "Config should be nil on error")
	})
}
