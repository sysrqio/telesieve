package scan

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_GzipRoundtrip(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "app.log")
	content := "probe /livez ok\napp ok\n"
	if err := os.WriteFile(plain, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	gzPath := filepath.Join(dir, "app.log.gz")
	f, err := os.Create(gzPath)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	_, _ = gw.Write([]byte(content))
	gw.Close()
	f.Close()

	n := 0
	if err := ParseFile(gzPath, func(_, _ string) error { n++; return nil }); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("lines %d", n)
	}
}
