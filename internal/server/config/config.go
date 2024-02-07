// Getting the settings for the application
package config

import (
	"bytes"
	"encoding/json"
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

type ConfigAdapter struct {
	Restore                  bool   `env:"RESTORE" json:"restore"`
	StoreInterval            int    `env:"STORE_INTERVAL" json:"store_interval"`
	HttpAddress              string `env:"ADDRESS" json:"address"`
	logsLevel                string
	DatabaseConnectionString string `env:"DATABASE_DSN" json:"database_dsn"`
	FileStoragePath          string `env:"STORE_PATH" json:"store_path"`
	storageType              string
	key                      string
	CryptoKeyPath            string `env:"CRYPTO_KEY" json:"crypto_key"`
	useCryptoKey             bool
	configFilePath           string
	TrustedSubnet            string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	isGrpc                   bool
}

func New() (*ConfigAdapter, error) {
	adapter := ConfigAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "metrics",
	}

	// get data from flags
	rootCmd.Flags().StringVarP(&adapter.HttpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().StringVarP(&adapter.logsLevel, "logs", "l", "info", "log level")
	rootCmd.Flags().IntVarP(&adapter.StoreInterval, "store_interval", "i", 300, "Interval for save data to disk")
	rootCmd.Flags().StringVarP(&adapter.FileStoragePath, "file_storage_path", "f", "./tmp/metrics-db.json", "Log file path")
	rootCmd.Flags().BoolVarP(&adapter.Restore, "restore", "r", true, "Load prev. data from file")
	rootCmd.Flags().StringVarP(&adapter.DatabaseConnectionString, "database_dsn", "d", "host=127.0.0.1 user=go password=go dbname=go sslmode=disable", "Database connection string")
	rootCmd.Flags().StringVarP(&adapter.key, "key", "k", "", "Key string")
	rootCmd.Flags().StringVarP(&adapter.CryptoKeyPath, "crypto-key", "y", "", "Path to key file")
	rootCmd.Flags().StringVarP(&adapter.TrustedSubnet, "trusted-subnet", "t", "127.0.0.1/24", "Trusted subnet")
	rootCmd.Flags().BoolVarP(&adapter.isGrpc, "grpc", "g", false, "Enable grpc server")

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
		adapter.HttpAddress = envHTTPAddress
	}
	if storeInterval, err := getEnvVariable("STORE_INTERVAL"); err == nil {
		adapter.StoreInterval, err = strconv.Atoi(storeInterval)
		if err != nil {
			return nil, err
		}
	}
	if fileStoragePath, err := getEnvVariable("FILE_STORAGE_PATH"); err == nil {
		adapter.FileStoragePath = fileStoragePath
		adapter.storageType = Disk
	}
	if restore, err := getEnvVariable("RESTORE"); err == nil {
		adapter.Restore, err = strconv.ParseBool(restore)
		if err != nil {
			return nil, err
		}
	}
	if databaseConnectionString, err := getEnvVariable("DATABASE_DSN"); err == nil {
		adapter.DatabaseConnectionString = databaseConnectionString
		adapter.storageType = Bd
	}
	if envKey, err := getEnvVariable("KEY"); err == nil {
		adapter.key = envKey
	}
	if pathKey, err := getEnvVariable("CRYPTO_KEY"); err == nil {
		adapter.CryptoKeyPath = pathKey
	}

	if trustedSubnet, err := getEnvVariable("TRUSTED_SUBNET"); err == nil {
		adapter.TrustedSubnet = trustedSubnet
	}

	// get data from config
	if adapter.configFilePath != "" {
		cfgFile, err := adapter.readConfig()
		if err == nil {
			return &cfgFile, nil
		}
	}
	return &adapter, nil
}

func (f *ConfigAdapter) GetServerAddress() string {
	return f.HttpAddress
}

func (f *ConfigAdapter) GetServerAddressWithScheme() string {
	return "http://" + f.GetServerAddress()
}

func (f *ConfigAdapter) GetLogsLevel() string {
	return f.logsLevel
}

func (f *ConfigAdapter) GetStoreInterval() time.Duration {
	if f.StoreInterval == 0 {
		return 1 * time.Second
	}
	return time.Duration(f.StoreInterval) * time.Second
}

func (f *ConfigAdapter) GetFileStoragePath() string {
	return f.FileStoragePath
}

func (f *ConfigAdapter) GetRestore() bool {
	return f.Restore
}

func (f *ConfigAdapter) IsGrpc() bool {
	return f.isGrpc
}

func (f *ConfigAdapter) GetDatabaseConnectionString() string {
	return f.DatabaseConnectionString
}

func (f *ConfigAdapter) GetStorageType() string {
	return f.storageType
}

func (f *ConfigAdapter) GetKey() string {
	if f.CryptoKeyPath != "" {
		if keyContent, err := os.ReadFile(f.CryptoKeyPath); err != nil {
			f.useCryptoKey = true
			return string(keyContent)
		}
	}
	return f.key
}

func (f *ConfigAdapter) UseCryptoKey() bool {
	return f.useCryptoKey
}
func (f *ConfigAdapter) GetTrustedSubnet() string {
	return f.TrustedSubnet
}

func getEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists && envVarValue != "" {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}

func (f *ConfigAdapter) readConfig() (ConfigAdapter, error) {
	flags := new(ConfigAdapter)

	data, err := os.ReadFile(f.configFilePath)
	if err != nil {
		return *flags, err
	}
	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(&flags); err != nil {
		return *flags, err
	}

	return *flags, nil
}
