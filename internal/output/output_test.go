package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sysrqio/telesieve/internal/cost"
	"github.com/sysrqio/telesieve/internal/pii"
	"github.com/sysrqio/telesieve/internal/scan"
)

func sampleReport() Report {
	st := scan.NewStats(3)
	st.TotalLines = 10
	st.WasteLines = 4
	st.ByCategory[scan.CategoryKubeletHealth] = 2
	st.ByCategory[scan.CategoryDebugChatter] = 2
	st.SampleWaste = []string{"GET /healthz"}
	c := cost.Estimate(1024, st.WastePercent(), cost.DefaultRates(cost.ProviderDatadog), 30)
	findings := []pii.Finding{
		{Kind: "jwt_bearer", Severity: pii.SeverityCritical, Masked: "Bear****.sig", Count: 1},
	}
	return BuildReport("datadog", st, c, findings)
}

func TestRender_JSONFields(t *testing.T) {
	rep := sampleReport()
	var buf bytes.Buffer
	if err := Render(&buf, FormatJSON, rep); err != nil {
		t.Fatal(err)
	}

	var decoded Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.TotalLines != 10 || decoded.WasteLines != 4 {
		t.Errorf("lines: total=%d waste=%d", decoded.TotalLines, decoded.WasteLines)
	}
	if !decoded.PIICritical || decoded.PIIFindingCount != 1 {
		t.Errorf("pii: critical=%v count=%d", decoded.PIICritical, decoded.PIIFindingCount)
	}
	if decoded.Cost.TotalMonthly <= 0 {
		t.Error("expected positive cost estimate")
	}
	if len(decoded.PIISamples) != 1 || decoded.PIISamples[0].Kind != "jwt_bearer" {
		t.Errorf("pii samples: %+v", decoded.PIISamples)
	}
}

func TestRender_MarkdownFields(t *testing.T) {
	rep := sampleReport()
	var buf bytes.Buffer
	if err := Render(&buf, FormatMarkdown, rep); err != nil {
		t.Fatal(err)
	}
	body := buf.String()
	for _, want := range []string{
		"# telesieve scan report",
		"| Provider | datadog |",
		"| Total lines | 10 |",
		"| Waste lines | 4",
		"PII findings",
		"critical: true",
		"## Waste categories",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("markdown missing %q\n%s", want, body)
		}
	}
}

func TestRender_TerminalSmoke(t *testing.T) {
	rep := sampleReport()
	var buf bytes.Buffer
	if err := Render(&buf, FormatTerminal, rep); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "telesieve scan report") {
		t.Error("terminal header")
	}
}
