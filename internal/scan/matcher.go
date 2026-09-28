package scan

import (
	"regexp"
	"strings"
)

// WasteCategory identifies a class of low-value log noise.
type WasteCategory string

const (
	CategoryKubeletHealth WasteCategory = "kubelet_health"
	CategoryRedisPing     WasteCategory = "redis_ping"
	CategoryDebugChatter  WasteCategory = "debug_chatter"
)

var (
	kubeletHealth = regexp.MustCompile(`(?i)(/healthz|/livez|/readyz|kubelet.*health)`)
	redisPingPong = regexp.MustCompile(`(?i)(redis.*\b(PING|PONG)\b|\bPING\b.*redis|\bPONG\b)`)
	debugChatter  = regexp.MustCompile(`(?i)(debug|trace|verbose|dumping|heartbeat|keepalive|noop)`)
)

// ClassifyLine returns waste categories matched by the line, or nil if not waste.
func ClassifyLine(line string) []WasteCategory {
	if strings.TrimSpace(line) == "" {
		return nil
	}
	var out []WasteCategory
	if kubeletHealth.MatchString(line) {
		out = append(out, CategoryKubeletHealth)
	}
	if redisPingPong.MatchString(line) {
		out = append(out, CategoryRedisPing)
	}
	if debugChatter.MatchString(line) {
		out = append(out, CategoryDebugChatter)
	}
	return out
}

// IsWaste returns true if the line matches any waste pattern.
func IsWaste(line string) bool {
	return len(ClassifyLine(line)) > 0
}
