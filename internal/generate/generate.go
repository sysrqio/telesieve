package generate

import (
	"fmt"
	"strings"

	"github.com/sysrqio/telesieve/internal/scan"
)

// Options control OTel collector snippet generation.
type Options struct {
	DropHealth  bool
	MaskTokens  bool
}

// CollectorSnippet returns YAML for filter/transform processors based on scan patterns.
func CollectorSnippet(path string, opts Options) (string, error) {
	var healthHits, redisHits, debugHits int
	err := scan.ParseFile(path, func(raw, line string) error {
		text := raw
		if len(scan.ClassifyLine(text)) == 0 {
			text = line
		}
		for _, c := range scan.ClassifyLine(text) {
			switch c {
			case scan.CategoryKubeletHealth:
				healthHits++
			case scan.CategoryRedisPing:
				redisHits++
			case scan.CategoryDebugChatter:
				debugHits++
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("# telesieve-generated OpenTelemetry Collector snippet\n")
	b.WriteString("processors:\n")

	if opts.DropHealth {
		b.WriteString("  filter/drop_health:\n")
		b.WriteString("    error_mode: ignore\n")
		b.WriteString("    logs:\n")
		b.WriteString("      log_record:\n")
		b.WriteString("        - 'IsMatch(body, \"(?i)(/healthz|/livez|/readyz|kubelet.*health)\")'\n")
	}
	if redisHits > 0 || opts.DropHealth {
		b.WriteString("  filter/drop_redis_ping:\n")
		b.WriteString("    error_mode: ignore\n")
		b.WriteString("    logs:\n")
		b.WriteString("      log_record:\n")
		b.WriteString("        - 'IsMatch(body, \"(?i)(redis.*(PING|PONG)|\\\\bPING\\\\b.*redis)\")'\n")
	}
	b.WriteString("  filter/drop_debug:\n")
	b.WriteString("    error_mode: ignore\n")
	b.WriteString("    logs:\n")
	b.WriteString("      log_record:\n")
	b.WriteString("        - 'IsMatch(body, \"(?i)(debug|trace|verbose|heartbeat|keepalive)\")'\n")

	if opts.MaskTokens {
		b.WriteString("  transform/mask_pii:\n")
		b.WriteString("    error_mode: ignore\n")
		b.WriteString("    log_statements:\n")
		b.WriteString("      - context: log\n")
		b.WriteString("        statements:\n")
		b.WriteString("          - replace_pattern(body, \"eyJ[A-Za-z0-9_-]+\\\\.eyJ[A-Za-z0-9_-]+\\\\.[A-Za-z0-9_-]+\", \"[REDACTED_JWT]\")\n")
		b.WriteString("          - replace_pattern(body, \"(?i)Authorization:\\\\s*[^\\\\s]+\", \"Authorization: [REDACTED]\")\n")
	}

	b.WriteString("\nservice:\n")
	b.WriteString("  pipelines:\n")
	b.WriteString("    logs:\n")
	b.WriteString("      processors:\n")
	if opts.DropHealth {
		b.WriteString("        - filter/drop_health\n")
	}
	b.WriteString("        - filter/drop_redis_ping\n")
	b.WriteString("        - filter/drop_debug\n")
	if opts.MaskTokens {
		b.WriteString("        - transform/mask_pii\n")
	}

	if healthHits+redisHits+debugHits == 0 {
		b.WriteString(fmt.Sprintf("# note: scan found no waste lines in %q; filters still emitted as templates\n", path))
	}
	return b.String(), nil
}
