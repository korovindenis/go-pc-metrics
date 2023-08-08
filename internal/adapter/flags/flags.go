package flags

import (
	"flag"
	"time"
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
	adapter := flagsAdapter{}
	// rootCmd := &cobra.Command{
	// 	Use:   "go-pc-metrics",
	// 	Short: "Application",
	// }

	// rootCmd.Flags().StringVarP(&adapter.httpAddress, "address", "a", "localhost:8080", "HTTP server address")
	// rootCmd.Flags().DurationVarP(&adapter.reportInterval, "report", "r", 10*time.Second, "Report interval")
	// rootCmd.Flags().DurationVarP(&adapter.pollInterval, "poll", "p", 2*time.Second, "Poll interval")

	// if err := rootCmd.Execute(); err != nil {
	// 	return nil, err
	// }

	flag.StringVar(&adapter.httpAddress, "a", "localhost:8080", "HTTP server address")
	flag.DurationVar(&adapter.reportInterval, "r", 10*time.Second, "Report interval")
	flag.DurationVar(&adapter.pollInterval, "p", 2*time.Second, "Poll interval")

	flag.Parse()

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
