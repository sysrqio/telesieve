package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sysrqio/telesieve/internal/cost"
	"github.com/sysrqio/telesieve/internal/pii"
	"github.com/sysrqio/telesieve/internal/scan"
)

// Format selects render target.
type Format string

const (
	FormatTerminal Format = "terminal"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Report bundles scan output for rendering.
type Report struct {
	Provider         string             `json:"provider"`
	TotalLines       int                `json:"total_lines"`
	WasteLines       int                `json:"waste_lines"`
	WastePercent     float64            `json:"waste_percent"`
	ByCategory       map[string]int     `json:"by_category"`
	PIIFindingCount  int                `json:"pii_finding_count"`
	PIICritical      bool               `json:"pii_critical"`
	SampleWaste      []string           `json:"sample_waste,omitempty"`
	Cost             cost.Result        `json:"cost"`
	PIISamples       []piiFindingSample `json:"pii_samples,omitempty"`
}

type piiFindingSample struct {
	Kind   string `json:"kind"`
	Masked string `json:"masked"`
}

// Render writes the report in the chosen format.
func Render(w io.Writer, format Format, r Report) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	case FormatMarkdown:
		return renderMarkdown(w, r)
	default:
		return renderTerminal(w, r)
	}
}

func renderTerminal(w io.Writer, r Report) error {
	fmt.Fprintf(w, "telesieve scan report\n")
	fmt.Fprintf(w, "─────────────────────\n")
	fmt.Fprintf(w, "Provider:      %s\n", r.Provider)
	fmt.Fprintf(w, "Lines:         %d total, %d waste (%.1f%%)\n", r.TotalLines, r.WasteLines, r.WastePercent)
	if len(r.ByCategory) > 0 {
		fmt.Fprintf(w, "Waste by type:\n")
		for k, v := range r.ByCategory {
			fmt.Fprintf(w, "  - %s: %d\n", k, v)
		}
	}
	fmt.Fprintf(w, "Monthly est.:  $%.2f ingest + $%.2f index = $%.2f (30d projection)\n",
		r.Cost.IngestionCost, r.Cost.IndexingCost, r.Cost.TotalMonthly)
	fmt.Fprintf(w, "Saveable:      ~$%.2f/mo if waste removed (%.2f GB/mo waste)\n",
		r.Cost.PotentialSavings, r.Cost.WasteGB)
	if r.PIIFindingCount > 0 {
		fmt.Fprintf(w, "PII:           %d finding(s), critical=%v\n", r.PIIFindingCount, r.PIICritical)
		for _, s := range r.PIISamples {
			fmt.Fprintf(w, "  - %s: %s\n", s.Kind, s.Masked)
		}
	}
	if len(r.SampleWaste) > 0 {
		fmt.Fprintf(w, "Sample waste:\n")
		for _, s := range r.SampleWaste {
			fmt.Fprintf(w, "  • %s\n", s)
		}
	}
	fmt.Fprintf(w, "\nTip: tune --ingestion-rate and --indexing-rate to match your vendor contract.\n")
	return nil
}

func renderMarkdown(w io.Writer, r Report) error {
	var b strings.Builder
	b.WriteString("# telesieve scan report\n\n")
	fmt.Fprintf(&b, "| Metric | Value |\n|--------|-------|\n")
	fmt.Fprintf(&b, "| Provider | %s |\n", r.Provider)
	fmt.Fprintf(&b, "| Total lines | %d |\n", r.TotalLines)
	fmt.Fprintf(&b, "| Waste lines | %d (%.1f%%) |\n", r.WasteLines, r.WastePercent)
	fmt.Fprintf(&b, "| Monthly cost (est.) | $%.2f |\n", r.Cost.TotalMonthly)
	fmt.Fprintf(&b, "| Potential savings | $%.2f/mo |\n", r.Cost.PotentialSavings)
	if r.PIIFindingCount > 0 {
		fmt.Fprintf(&b, "| PII findings | %d (critical: %v) |\n", r.PIIFindingCount, r.PIICritical)
	}
	b.WriteString("\n")
	if len(r.ByCategory) > 0 {
		b.WriteString("## Waste categories\n\n")
		for k, v := range r.ByCategory {
			fmt.Fprintf(&b, "- **%s**: %d\n", k, v)
		}
		b.WriteString("\n")
	}
	_, err := w.Write([]byte(b.String()))
	return err
}

// BuildReport maps scan stats and cost into Report.
func BuildReport(provider string, st *scan.Stats, c cost.Result, findings []pii.Finding) Report {
	byCat := make(map[string]int)
	for k, v := range st.ByCategory {
		byCat[string(k)] = v
	}
	var samples []piiFindingSample
	seen := make(map[string]bool)
	for _, f := range findings {
		if len(samples) >= 5 {
			break
		}
		key := f.Kind + f.Masked
		if seen[key] {
			continue
		}
		seen[key] = true
		samples = append(samples, piiFindingSample{Kind: f.Kind, Masked: f.Masked})
	}
	critical := false
	for _, f := range findings {
		if f.Severity == pii.SeverityCritical {
			critical = true
			break
		}
	}
	return Report{
		Provider:        provider,
		TotalLines:      st.TotalLines,
		WasteLines:      st.WasteLines,
		WastePercent:    st.WastePercent(),
		ByCategory:      byCat,
		PIIFindingCount: len(findings),
		PIICritical:     critical,
		SampleWaste:     st.SampleWaste,
		Cost:            c,
		PIISamples:      samples,
	}
}
