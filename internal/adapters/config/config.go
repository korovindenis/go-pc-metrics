package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/spf13/cobra"
)

type configAdapter struct {
	reportInterval  int
	pollInterval    int
	httpAddress     string
	logsLevel       string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func New(isServer bool) (*configAdapter, error) {
	adapter := configAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "metrics",
	}

	// get data from flags
	rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().StringVarP(&adapter.logsLevel, "logs", "l", "info", "log level")

	// server
	if isServer {
		rootCmd.Flags().IntVarP(&adapter.storeInterval, "store_interval", "i", 300, "Interval for save data to disk")
		rootCmd.Flags().StringVarP(&adapter.fileStoragePath, "file_storage_path", "f", "metrics-db.json", "Log file path")
		rootCmd.Flags().BoolVarP(&adapter.restore, "restore", "r", true, "Load prev. data from file")
	} else {
		// agent
		rootCmd.Flags().IntVarP(&adapter.reportInterval, "report", "r", 10, "Metrics report interval")
		rootCmd.Flags().IntVarP(&adapter.pollInterval, "poll", "p", 2, "Metrics poll interval")
	}
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	// if env var not empty
	// get data from env
	if envHTTPAddress, err := getEnvVariable("ADDRESS"); err == nil {
		adapter.httpAddress = envHTTPAddress
	}
	if reportInterval, err := getEnvVariable("REPORT_INTERVAL"); err == nil {
		adapter.reportInterval, _ = strconv.Atoi(reportInterval)
	}
	if pollInterval, err := getEnvVariable("POLL_INTERVAL"); err == nil {
		adapter.pollInterval, _ = strconv.Atoi(pollInterval)
	}
	if storeInterval, err := getEnvVariable("STORE_INTERVAL"); err == nil {
		adapter.storeInterval, _ = strconv.Atoi(storeInterval)
	}
	if fileStoragePath, err := getEnvVariable("FILE_STORAGE_PATH"); err == nil {
		adapter.fileStoragePath = fileStoragePath
	}
	if restore, err := getEnvVariable("RESTORE"); err == nil {
		adapter.restore, _ = strconv.ParseBool(restore)
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

func (f *configAdapter) GetStoreInterval() time.Duration {
	if f.storeInterval == 0 {
		return 1 * time.Second
	}
	return time.Duration(f.storeInterval) * time.Second
}

func (f *configAdapter) GetFileStoragePath() string {
	return filepath.Join(os.TempDir(), f.fileStoragePath)
}

func (f *configAdapter) GetRestore() bool {
	return f.restore
}

func getEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists && envVarValue != "" {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
