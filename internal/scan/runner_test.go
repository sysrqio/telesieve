package scan

import (
	"path/filepath"
	"testing"
)

func TestRunScan_MixedLog(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "mixed.log")
	st := NewStats(5)
	_, _, err := RunScan(path, nil, st)
	if err != nil {
		t.Fatal(err)
	}
	if st.WasteLines < 4 {
		t.Fatalf("expected >=4 waste lines, got %d cats=%v", st.WasteLines, st.ByCategory)
	}
}
