package cmdroot

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRoot builds the telesieve CLI.
func NewRoot(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "telesieve",
		Short: "Local telemetry cost and PII scanner for log dumps",
		Long:  "Stream-parse log dumps (JSON lines, logfmt, plain text, optional gzip) to find noisy waste and sensitive data—entirely offline.",
	}
	root.AddCommand(newScanCmd())
	root.AddCommand(newGenerateCmd())
	root.AddCommand(newVersionCmd(version))
	return root
}

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("telesieve %s\n", version)
		},
	}
}

func exitCode(err error, code int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(code)
}
