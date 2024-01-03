// Getting the settings for the application

package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigAdapter(t *testing.T) {
	tests := []struct {
		name                       string
		envVarValues               map[string]string
		expectedHTTPAddress        string
		expectedLogsLevel          string
		expectedStoreInterval      time.Duration
		expectedFileStoragePath    string
		expectedRestore            bool
		expectedDatabaseConnString string
		expectedStorageType        string
		expectedKey                string
	}{
		{
			name: "All Values Set",
			envVarValues: map[string]string{
				"ADDRESS":           "localhost:8080",
				"STORE_INTERVAL":    "300",
				"FILE_STORAGE_PATH": "./tmp/metrics-db.json",
				"RESTORE":           "true",
				"DATABASE_DSN":      "host=127.0.0.1 user=go password=go dbname=go sslmode=disable",
				"KEY":               "some_key",
			},
			expectedHTTPAddress:        "localhost:8080",
			expectedLogsLevel:          "info",
			expectedStoreInterval:      300 * time.Second,
			expectedFileStoragePath:    "./tmp/metrics-db.json",
			expectedRestore:            true,
			expectedDatabaseConnString: "host=127.0.0.1 user=go password=go dbname=go sslmode=disable",
			expectedStorageType:        "database",
			expectedKey:                "some_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVarValues {
				err := os.Setenv(key, value)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			// Act
			config, err := New()
			assert.NoError(t, err)

			// Assert
			assert.Equal(t, tt.expectedHTTPAddress, config.GetServerAddress())
			assert.Equal(t, tt.expectedLogsLevel, config.GetLogsLevel())
			assert.Equal(t, tt.expectedStoreInterval, config.GetStoreInterval())
			assert.Equal(t, tt.expectedFileStoragePath, config.GetFileStoragePath())
			assert.Equal(t, tt.expectedRestore, config.GetRestore())
			assert.Equal(t, tt.expectedDatabaseConnString, config.GetDatabaseConnectionString())
			assert.Equal(t, tt.expectedStorageType, config.GetStorageType())
			assert.Equal(t, tt.expectedKey, config.GetKey())
		})
	}
}
