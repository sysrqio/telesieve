package cmdroot

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/sysrqio/telesieve/internal/cost"
	"github.com/sysrqio/telesieve/internal/output"
	"github.com/sysrqio/telesieve/internal/pii"
	"github.com/sysrqio/telesieve/internal/scan"
)

func newScanCmd() *cobra.Command {
	var (
		provider      string
		ingestionRate float64
		indexingRate  float64
		outFmt        string
		maskPII       bool
	)
	cmd := &cobra.Command{
		Use:   "scan [file]",
		Short: "Scan a log dump for waste and PII",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			st := scan.NewStats(5)
			scanner := pii.NewScanner(maskPII)
			size, findings, err := scan.RunScan(path, scanner, st)
			if err != nil {
				return err
			}
			var rates cost.Rates
			if provider == "custom" {
				rates = cost.Rates{IngestionPerGB: ingestionRate, IndexingPerGB: indexingRate}
			} else {
				rates = cost.DefaultRates(cost.Provider(provider))
				if cmd.Flags().Changed("ingestion-rate") {
					rates.IngestionPerGB = ingestionRate
				}
				if cmd.Flags().Changed("indexing-rate") {
					rates.IndexingPerGB = indexingRate
				}
			}

			costRes := cost.Estimate(size, st.WastePercent(), rates, 30)
			rep := output.BuildReport(provider, st, costRes, findings)
			format := output.Format(outFmt)
			if err := output.Render(os.Stdout, format, rep); err != nil {
				return err
			}
			if maskPII && st.PIICritical {
				os.Exit(2)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "datadog", "cost model: datadog|splunk|dynatrace|custom")
	cmd.Flags().Float64Var(&ingestionRate, "ingestion-rate", 0.10, "USD per GB ingested")
	cmd.Flags().Float64Var(&indexingRate, "indexing-rate", 1.70, "USD per GB indexed")
	cmd.Flags().StringVar(&outFmt, "output", "terminal", "output format: terminal|json|markdown")
	cmd.Flags().BoolVar(&maskPII, "mask-pii", true, "mask PII in output; exit 2 if critical PII found")
	return cmd
}
