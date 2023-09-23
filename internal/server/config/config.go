package config

import (
	"os"
	"strconv"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/spf13/cobra"
)

const (
	Bd   = "database"
	Disk = "disk"
)

type configAdapter struct {
	restore                  bool
	storeInterval            int
	httpAddress              string
	logsLevel                string
	databaseConnectionString string
	fileStoragePath          string
	storageType              string
	key                      string
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
	rootCmd.Flags().IntVarP(&adapter.storeInterval, "store_interval", "i", 300, "Interval for save data to disk")
	rootCmd.Flags().StringVarP(&adapter.fileStoragePath, "file_storage_path", "f", "./tmp/metrics-db.json", "Log file path")
	rootCmd.Flags().BoolVarP(&adapter.restore, "restore", "r", true, "Load prev. data from file")
	rootCmd.Flags().StringVarP(&adapter.databaseConnectionString, "database_dsn", "d", "host=127.0.0.1 user=go password=go dbname=go sslmode=disable", "Database connection string")
	rootCmd.Flags().StringVarP(&adapter.key, "key", "k", "", "Key string")

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}
	if rootCmd.Flags().Changed("file_storage_path") {
		adapter.storageType = Disk
	}
	if rootCmd.Flags().Changed("database_dsn") {
		adapter.storageType = Bd
	}

	// if env var not empty
	// get data from env
	if envHTTPAddress, err := getEnvVariable("ADDRESS"); err == nil {
		adapter.httpAddress = envHTTPAddress
	}
	if storeInterval, err := getEnvVariable("STORE_INTERVAL"); err == nil {
		adapter.storeInterval, err = strconv.Atoi(storeInterval)
		if err != nil {
			return nil, err
		}
	}
	if fileStoragePath, err := getEnvVariable("FILE_STORAGE_PATH"); err == nil {
		adapter.fileStoragePath = fileStoragePath
		adapter.storageType = Disk
	}
	if restore, err := getEnvVariable("RESTORE"); err == nil {
		adapter.restore, err = strconv.ParseBool(restore)
		if err != nil {
			return nil, err
		}
	}
	if databaseConnectionString, err := getEnvVariable("DATABASE_DSN"); err == nil {
		adapter.databaseConnectionString = databaseConnectionString
		adapter.storageType = Bd
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
	return f.fileStoragePath
}

func (f *configAdapter) GetRestore() bool {
	return f.restore
}

func (f *configAdapter) GetDatabaseConnectionString() string {
	return f.databaseConnectionString
}

func (f *configAdapter) GetStorageType() string {
	return f.storageType
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
