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
	reportInterval int
	pollInterval   int
}

func New() (IFlags, error) {
	adapter := flagsAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "Application",
	}

	rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().IntVarP(&adapter.reportInterval, "report", "r", 10, "Report interval")
	rootCmd.Flags().IntVarP(&adapter.pollInterval, "poll", "p", 2, "Poll interval")

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
