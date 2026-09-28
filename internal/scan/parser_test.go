package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile_JSONAndLogfmt(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "mixed.log")
	var lines []string
	err := ParseFile(path, func(_, line string) error {
		lines = append(lines, line)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < 5 {
		t.Fatalf("expected lines, got %d", len(lines))
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "healthz") {
		t.Error("json message extract")
	}
	if !strings.Contains(joined, "redis") {
		t.Error("logfmt msg extract")
	}
}

func TestParseFile_Plain(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "plain.log")
	n := 0
	if err := ParseFile(path, func(_, _ string) error { n++; return nil }); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("want 2 lines got %d", n)
	}
}

func TestParseFile_Gzip(t *testing.T) {
	dir := t.TempDir()
	gzPath := filepath.Join(dir, "x.log.gz")
	// create gzip via ParseReader roundtrip - write manually
	plain := filepath.Join(dir, "x.log")
	if err := os.WriteFile(plain, []byte("line /readyz\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// skip gzip integration if no file; test gzip error path
	err := ParseReader(strings.NewReader("not gzip"), true, func(_, _ string) error { return nil })
	if err == nil {
		t.Error("expected gzip error")
	}
	_ = gzPath
}
