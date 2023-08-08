package flags

import (
	"strconv"
	"time"

	"github.com/korovindenis/go-pc-metrics/internal/adapter/env"
	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
	"github.com/spf13/cobra"
)

// functions
type IFlags interface {
	GetHTTPAddress() string
	GetHTTPAddressWithScheme() string
	GetReportInterval() time.Duration
	GetPollInterval() time.Duration
}

type flagsAdapter struct {
	httpAddress    string
	reportInterval int
	pollInterval   int
}

func New(envObj env.IEnv) (IFlags, error) {
	adapter := flagsAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "Application",
	}

	if envHTTPAddress, err := envObj.GetEnvVariable("ADDRESS"); envHTTPAddress == "" || err == entity.ErrEnvVarNotFound {
		rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	} else {
		adapter.httpAddress = envHTTPAddress
	}

	if reportInterval, err := envObj.GetEnvVariable("REPORT_INTERVAL"); reportInterval == "" || err == entity.ErrEnvVarNotFound {
		rootCmd.Flags().IntVarP(&adapter.reportInterval, "report", "r", 10, "Report interval")
	} else {
		adapter.reportInterval, _ = strconv.Atoi(reportInterval)
	}

	if pollInterval, err := envObj.GetEnvVariable("POLL_INTERVAL"); pollInterval == "" || err == entity.ErrEnvVarNotFound {
		rootCmd.Flags().IntVarP(&adapter.pollInterval, "poll", "p", 2, "Poll interval")
	} else {
		adapter.pollInterval, _ = strconv.Atoi(pollInterval)
	}

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	return &adapter, nil
}

func (f *flagsAdapter) GetHTTPAddress() string {
	return f.httpAddress
}

func (f *flagsAdapter) GetHTTPAddressWithScheme() string {
	return "http://" + f.httpAddress
}

func (f *flagsAdapter) GetReportInterval() time.Duration {
	return time.Duration(f.reportInterval) * time.Second
}

func (f *flagsAdapter) GetPollInterval() time.Duration {
	return time.Duration(f.pollInterval) * time.Second
}
