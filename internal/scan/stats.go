package scan

// Stats aggregates scan results.
type Stats struct {
	TotalLines   int
	WasteLines   int
	ByCategory   map[WasteCategory]int
	PIIFindings  int
	PIICritical  bool
	SampleWaste  []string
	maxSamples   int
}

// NewStats creates an empty Stats with sample cap.
func NewStats(maxSamples int) *Stats {
	if maxSamples <= 0 {
		maxSamples = 5
	}
	return &Stats{
		ByCategory: make(map[WasteCategory]int),
		maxSamples: maxSamples,
	}
}

// RecordLine updates counters for one log line.
func (s *Stats) RecordLine(line string, waste []WasteCategory) {
	s.TotalLines++
	if len(waste) == 0 {
		return
	}
	s.WasteLines++
	seen := make(map[WasteCategory]bool)
	for _, c := range waste {
		if !seen[c] {
			s.ByCategory[c]++
			seen[c] = true
		}
	}
	if len(s.SampleWaste) < s.maxSamples {
		s.SampleWaste = append(s.SampleWaste, truncate(line, 120))
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

// WastePercent returns waste line ratio 0–100.
func (s *Stats) WastePercent() float64 {
	if s.TotalLines == 0 {
		return 0
	}
	return float64(s.WasteLines) / float64(s.TotalLines) * 100
}
