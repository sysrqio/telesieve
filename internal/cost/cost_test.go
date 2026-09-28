package cost

import "testing"

func TestDefaultRates(t *testing.T) {
	d := DefaultRates(ProviderDatadog)
	if d.IngestionPerGB != 0.10 || d.IndexingPerGB != 1.70 {
		t.Fatalf("datadog rates %+v", d)
	}
	s := DefaultRates(ProviderSplunk)
	if s.IngestionPerGB != 0.15 || s.IndexingPerGB != 0.90 {
		t.Fatalf("splunk %+v", s)
	}
	dt := DefaultRates(ProviderDynatrace)
	if dt.IngestionPerGB != 0.20 || dt.IndexingPerGB != 0.50 {
		t.Fatalf("dynatrace %+v", dt)
	}
	c := DefaultRates(ProviderCustom)
	if c.IngestionPerGB != 0.10 || c.IndexingPerGB != 1.70 {
		t.Fatalf("custom default %+v", c)
	}
	c2 := DefaultRates(Provider("unknown"))
	if c2.IndexingPerGB != 1.70 {
		t.Fatalf("unknown %+v", c2)
	}
}

func TestEstimate_Waste(t *testing.T) {
	const oneGB = 1024 * 1024 * 1024
	r := Estimate(oneGB, 50, Rates{IngestionPerGB: 0.10, IndexingPerGB: 1.00}, 30)
	if r.MonthlyGB != 30 {
		t.Fatalf("monthly gb %v", r.MonthlyGB)
	}
	if r.WasteGB != 15 {
		t.Fatalf("waste gb %v", r.WasteGB)
	}
	if r.PotentialSavings != 15*(0.10+1.00) {
		t.Fatalf("savings %v", r.PotentialSavings)
	}
}

func TestEstimate_ZeroBytes(t *testing.T) {
	r := Estimate(0, 0, DefaultRates(ProviderDatadog), 30)
	if r.TotalMonthly != 0 {
		t.Fatalf("%v", r.TotalMonthly)
	}
}

func TestEstimate_ClampsWaste(t *testing.T) {
	r := Estimate(1000, 150, Rates{IngestionPerGB: 1, IndexingPerGB: 1}, 30)
	if r.WasteGB > r.MonthlyGB {
		t.Fatal("waste clamp")
	}
}
