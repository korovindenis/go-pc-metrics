package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/spf13/cobra"
)

type configAdapter struct {
	httpAddress    string
	reportInterval int
	pollInterval   int
}

func New() (*configAdapter, error) {
	adapter := configAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "metrics",
	}

	if envHTTPAddress, err := getEnvVariable("ADDRESS"); envHTTPAddress == "" || errors.Is(err, entity.ErrEnvVarNotFound) {
		rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	} else {
		adapter.httpAddress = envHTTPAddress
	}

	if reportInterval, err := getEnvVariable("REPORT_INTERVAL"); reportInterval == "" || errors.Is(err, entity.ErrEnvVarNotFound) {
		rootCmd.Flags().IntVarP(&adapter.reportInterval, "report", "r", 10, "Metrics report interval")
	} else {
		adapter.reportInterval, _ = strconv.Atoi(reportInterval)
	}

	if pollInterval, err := getEnvVariable("POLL_INTERVAL"); pollInterval == "" || errors.Is(err, entity.ErrEnvVarNotFound) {
		rootCmd.Flags().IntVarP(&adapter.pollInterval, "poll", "p", 2, "Metrics poll interval")
	} else {
		adapter.pollInterval, _ = strconv.Atoi(pollInterval)
	}

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	return &adapter, nil
}

func (f *configAdapter) GetHTTPAddress() string {
	return f.httpAddress
}

func (f *configAdapter) GetHTTPAddressWithScheme() string {
	return "http://" + f.httpAddress
}

func (f *configAdapter) GetReportInterval() time.Duration {
	return time.Duration(f.reportInterval) * time.Second
}

func (f *configAdapter) GetPollInterval() time.Duration {
	return time.Duration(f.pollInterval) * time.Second
}

func getEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
