package flags

import (
	"time"

	"github.com/spf13/cobra"
)

// functions
type IFlags interface {
	GetHttpAddress() string
	GetReportInterval() time.Duration
	GetPollInterval() time.Duration
}

type flagsAdapter struct {
	httpAddress    string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func New() (IFlags, error) {
	adapter := flagsAdapter{}
	rootCmd := &cobra.Command{
		Use:   "go-pc-metrics",
		Short: "Application",
	}

	rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	rootCmd.Flags().DurationVarP(&adapter.reportInterval, "report", "r", 10*time.Second, "Report interval")
	rootCmd.Flags().DurationVarP(&adapter.pollInterval, "poll", "p", 2*time.Second, "Poll interval")

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	return &adapter, nil
}

func (f *flagsAdapter) GetHttpAddress() string {
	return f.httpAddress
}

func (f *flagsAdapter) GetReportInterval() time.Duration {
	return f.reportInterval
}

func (f *flagsAdapter) GetPollInterval() time.Duration {
	return f.pollInterval
}