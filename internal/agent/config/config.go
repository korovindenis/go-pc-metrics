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

type ConfigAdapter struct {
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	HTTPAddress    string `env:"ADDRESS" json:"address"`
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKeyPath  string `env:"CRYPTO_KEY" json:"crypto_key"`
	logsLevel      string
	key            string
	useCryptoKey   bool
	configFilePath string
}

func New() (*ConfigAdapter, error) {
	adapter := ConfigAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "metrics",
	}

	// get data from flags
	rootCmd.Flags().StringVarP(&adapter.HTTPAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().StringVarP(&adapter.logsLevel, "logs", "i", "info", "log level")
	rootCmd.Flags().IntVarP(&adapter.ReportInterval, "report", "r", 10, "Metrics report interval")
	rootCmd.Flags().IntVarP(&adapter.PollInterval, "poll", "p", 2, "Metrics poll interval")
	rootCmd.Flags().StringVarP(&adapter.key, "key", "k", "", "Key string")
	rootCmd.Flags().IntVarP(&adapter.RateLimit, "limit", "l", 1, "Limit http reg")
	rootCmd.Flags().StringVarP(&adapter.CryptoKeyPath, "crypto-key", "y", "", "Path to key file")
	rootCmd.Flags().StringVarP(&adapter.configFilePath, "config", "o", "", "Path to config file")

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	// if env var not empty
	// get data from env
	if envHTTPAddress, err := getEnvVariable("ADDRESS"); err == nil {
		adapter.HTTPAddress = envHTTPAddress
	}
	if reportInterval, err := getEnvVariable("REPORT_INTERVAL"); err == nil {
		adapter.ReportInterval, err = strconv.Atoi(reportInterval)
		if err != nil {
			return nil, err
		}
	}
	if pollInterval, err := getEnvVariable("POLL_INTERVAL"); err == nil {
		adapter.PollInterval, err = strconv.Atoi(pollInterval)
		if err != nil {
			return nil, err
		}
	}
	if envKey, err := getEnvVariable("KEY"); err == nil {
		adapter.key = envKey
	}
	if rateLimit, err := getEnvVariable("RATE_LIMIT"); err == nil {
		adapter.RateLimit, err = strconv.Atoi(rateLimit)
		if err != nil {
			return nil, err
		}
	}
	if pathKey, err := getEnvVariable("CRYPTO_KEY"); err == nil {
		adapter.CryptoKeyPath = pathKey
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
	return f.HTTPAddress
}

func (f *ConfigAdapter) GetServerAddressWithScheme() string {
	return "http://" + f.GetServerAddress()
}

func (f *ConfigAdapter) GetReportInterval() time.Duration {
	return time.Duration(f.ReportInterval) * time.Second
}

func (f *ConfigAdapter) GetPollInterval() time.Duration {
	return time.Duration(f.PollInterval) * time.Second
}

func (f *ConfigAdapter) GetLogsLevel() string {
	return f.logsLevel
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

func (f *ConfigAdapter) GetRateLimit() int {
	return f.RateLimit
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
