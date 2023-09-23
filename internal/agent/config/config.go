package config

import (
	"os"
	"strconv"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/spf13/cobra"
)

type configAdapter struct {
	reportInterval int
	pollInterval   int
	httpAddress    string
	logsLevel      string
	key            string
}

func New() (*configAdapter, error) {
	adapter := configAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "metrics",
	}

	// get data from flags
	rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().StringVarP(&adapter.logsLevel, "logs", "l", "info", "log level")
	rootCmd.Flags().IntVarP(&adapter.reportInterval, "report", "r", 10, "Metrics report interval")
	rootCmd.Flags().IntVarP(&adapter.pollInterval, "poll", "p", 2, "Metrics poll interval")
	rootCmd.Flags().StringVarP(&adapter.key, "key", "k", "", "Key string")
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	// if env var not empty
	// get data from env
	if envHTTPAddress, err := getEnvVariable("ADDRESS"); err == nil {
		adapter.httpAddress = envHTTPAddress
	}
	if reportInterval, err := getEnvVariable("REPORT_INTERVAL"); err == nil {
		adapter.reportInterval, err = strconv.Atoi(reportInterval)
		if err != nil {
			return nil, err
		}
	}
	if pollInterval, err := getEnvVariable("POLL_INTERVAL"); err == nil {
		adapter.pollInterval, err = strconv.Atoi(pollInterval)
		if err != nil {
			return nil, err
		}
	}
	if envKey, err := getEnvVariable("KEY"); err == nil {
		adapter.key = envKey
	}
	return &adapter, nil
}

func (f *configAdapter) GetServerAddress() string {
	return f.httpAddress
}

func (f *configAdapter) GetServerAddressWithScheme() string {
	return "http://" + f.GetServerAddress()
}

func (f *configAdapter) GetReportInterval() time.Duration {
	return time.Duration(f.reportInterval) * time.Second
}

func (f *configAdapter) GetPollInterval() time.Duration {
	return time.Duration(f.pollInterval) * time.Second
}

func (f *configAdapter) GetLogsLevel() string {
	return f.logsLevel
}

func (f *configAdapter) GetKey() string {
	return f.key
}

func getEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists && envVarValue != "" {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
