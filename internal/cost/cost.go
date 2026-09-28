package cost

// Provider identifies telemetry billing vendor assumptions.
type Provider string

const (
	ProviderDatadog   Provider = "datadog"
	ProviderSplunk    Provider = "splunk"
	ProviderDynatrace Provider = "dynatrace"
	ProviderCustom    Provider = "custom"
)

// Rates are USD per GB for ingestion and indexing (observability list prices vary; override via flags).
type Rates struct {
	IngestionPerGB float64
	IndexingPerGB  float64
}

// DefaultRates returns list-style defaults for known providers (ingestion, indexing USD/GB).
func DefaultRates(p Provider) Rates {
	switch p {
	case ProviderSplunk:
		return Rates{IngestionPerGB: 0.15, IndexingPerGB: 0.90}
	case ProviderDynatrace:
		return Rates{IngestionPerGB: 0.20, IndexingPerGB: 0.50}
	case ProviderDatadog:
		return Rates{IngestionPerGB: 0.10, IndexingPerGB: 1.70}
	default:
		return Rates{IngestionPerGB: 0.10, IndexingPerGB: 1.70}
	}
}

// Estimate projects monthly cost from daily dump stats.
// totalBytes: size of analyzed dump; wastePercent: 0–100; daysInMonth scales daily → monthly.
func Estimate(totalBytes int64, wastePercent float64, rates Rates, daysInMonth float64) Result {
	if daysInMonth <= 0 {
		daysInMonth = 30
	}
	gbPerDay := float64(totalBytes) / (1024 * 1024 * 1024)
	if gbPerDay < 0 {
		gbPerDay = 0
	}
	monthlyGB := gbPerDay * daysInMonth
	wasteFrac := wastePercent / 100
	if wasteFrac < 0 {
		wasteFrac = 0
	}
	if wasteFrac > 1 {
		wasteFrac = 1
	}
	wasteGB := monthlyGB * wasteFrac
	usefulGB := monthlyGB - wasteGB
	ingestCost := monthlyGB * rates.IngestionPerGB
	indexCost := monthlyGB * rates.IndexingPerGB
	totalCost := ingestCost + indexCost
	savings := wasteGB * (rates.IngestionPerGB + rates.IndexingPerGB)
	return Result{
		MonthlyGB:       monthlyGB,
		WasteGB:         wasteGB,
		UsefulGB:        usefulGB,
		WastePercent:    wastePercent,
		IngestionCost:   ingestCost,
		IndexingCost:    indexCost,
		TotalMonthly:    totalCost,
		PotentialSavings: savings,
	}
}

// Result holds cost projection output.
type Result struct {
	MonthlyGB        float64
	WasteGB          float64
	UsefulGB         float64
	WastePercent     float64
	IngestionCost    float64
	IndexingCost     float64
	TotalMonthly     float64
	PotentialSavings float64
}
