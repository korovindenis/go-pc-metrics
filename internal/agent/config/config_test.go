package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigAdapter(t *testing.T) {
	tests := []struct {
		name               string
		setEnvVariables    func()
		expectedServerAddr string
		expectedSchemeAddr string
		expectedReportInt  time.Duration
		expectedPollInt    time.Duration
		expectedLogsLevel  string
		expectedKey        string
		expectedRateLimit  int
	}{
		{
			name:               "Default Values",
			setEnvVariables:    func() {},
			expectedServerAddr: "localhost:8080",
			expectedSchemeAddr: "http://localhost:8080",
			expectedReportInt:  10 * time.Second,
			expectedPollInt:    2 * time.Second,
			expectedLogsLevel:  "info",
			expectedKey:        "",
			expectedRateLimit:  1,
		},
		{
			name: "Custom Values",
			setEnvVariables: func() {
				os.Setenv("ADDRESS", "customAddress")
				os.Setenv("REPORT_INTERVAL", "5")
				os.Setenv("POLL_INTERVAL", "3")
				os.Setenv("KEY", "customKey")
				os.Setenv("RATE_LIMIT", "2")
			},
			expectedServerAddr: "customAddress",
			expectedSchemeAddr: "http://customAddress",
			expectedReportInt:  5 * time.Second,
			expectedPollInt:    3 * time.Second,
			expectedLogsLevel:  "info",
			expectedKey:        "customKey",
			expectedRateLimit:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnvVariables()

			adapter, err := New()
			assert.NoError(t, err, "New should not return an error")

			assert.Equal(t, tt.expectedServerAddr, adapter.GetServerAddress())
			assert.Equal(t, tt.expectedSchemeAddr, adapter.GetServerAddressWithScheme())
			assert.Equal(t, tt.expectedReportInt, adapter.GetReportInterval())
			assert.Equal(t, tt.expectedPollInt, adapter.GetPollInterval())
			assert.Equal(t, tt.expectedLogsLevel, adapter.GetLogsLevel())
			assert.Equal(t, tt.expectedKey, adapter.GetKey())
			assert.Equal(t, tt.expectedRateLimit, adapter.GetRateLimit())
		})
	}
}
