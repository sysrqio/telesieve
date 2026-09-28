package scan

import "testing"

func TestStats_WastePercent(t *testing.T) {
	st := NewStats(2)
	st.RecordLine("a", nil)
	st.RecordLine("b", []WasteCategory{CategoryDebugChatter})
	if st.WastePercent() != 50 {
		t.Fatalf("%v", st.WastePercent())
	}
}
