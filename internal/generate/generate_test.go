package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectorSnippet(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "mixed.log")
	out, err := CollectorSnippet(path, Options{DropHealth: true, MaskTokens: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "filter/drop_health") {
		t.Error("health filter")
	}
	if !strings.Contains(out, "transform/mask_pii") {
		t.Error("mask")
	}
}
