package scan

import "testing"

func TestClassifyLine_Kubelet(t *testing.T) {
	cats := ClassifyLine("GET /healthz 200")
	if len(cats) == 0 {
		t.Fatal("expected waste")
	}
	found := false
	for _, c := range cats {
		if c == CategoryKubeletHealth {
			found = true
		}
	}
	if !found {
		t.Fatalf("got %v", cats)
	}
}

func TestClassifyLine_Redis(t *testing.T) {
	if !IsWaste(`redis: PING`) {
		t.Error("redis ping")
	}
}

func TestClassifyLine_Debug(t *testing.T) {
	if !IsWaste("DEBUG verbose trace dump") {
		t.Error("debug")
	}
}

func TestClassifyLine_Clean(t *testing.T) {
	if IsWaste("payment capture succeeded order=9912") {
		t.Error("should not be waste")
	}
}

func TestClassifyLine_Empty(t *testing.T) {
	if ClassifyLine("  ") != nil {
		t.Error("empty")
	}
}
