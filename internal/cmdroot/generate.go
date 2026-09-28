package cmdroot

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sysrqio/telesieve/internal/generate"
)

func newGenerateCmd() *cobra.Command {
	var dropHealth, maskTokens bool
	cmd := &cobra.Command{
		Use:   "generate [file]",
		Short: "Generate OpenTelemetry Collector filter/transform YAML from scan patterns",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippet, err := generate.CollectorSnippet(args[0], generate.Options{
				DropHealth: dropHealth,
				MaskTokens: maskTokens,
			})
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(os.Stdout, snippet)
			return err
		},
	}
	cmd.Flags().BoolVar(&dropHealth, "drop-health", true, "include health/live/ready probe filters")
	cmd.Flags().BoolVar(&maskTokens, "mask-tokens", true, "include JWT/Authorization masking transforms")
	return cmd
}
