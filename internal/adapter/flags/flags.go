package flags

import (
	"time"

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
	reportInterval time.Duration
	pollInterval   time.Duration
}

func New() (IFlags, error) {
	var reportInterval, pollInterval int
	var httpAddress string
	adapter := flagsAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "Application",
	}

	rootCmd.Flags().StringVarP(&httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().IntVarP(&reportInterval, "report", "r", 10, "Report interval")
	rootCmd.Flags().IntVarP(&pollInterval, "poll", "p", 2, "Poll interval")

	adapter.httpAddress = httpAddress
	adapter.reportInterval = time.Duration(reportInterval) * time.Second
	adapter.pollInterval = time.Duration(pollInterval) * time.Second

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
	return f.reportInterval
}

func (f *flagsAdapter) GetPollInterval() time.Duration {
	return f.pollInterval
}
