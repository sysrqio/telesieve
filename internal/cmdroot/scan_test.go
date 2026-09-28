package cmdroot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sysrqio/telesieve/internal/output"
)

func testdataPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("..", "..", "testdata", name)
}

func TestScan_MissingFile_Exit1(t *testing.T) {
	root := NewRoot("test")
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{"scan", filepath.Join(t.TempDir(), "does-not-exist.log")})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	var exit *ExitError
	if errors.As(err, &exit) {
		t.Fatalf("expected ordinary error, got exit status %d", exit.Code)
	}
}

func TestScan_JWTWithMaskPII_Exit2(t *testing.T) {
	root := NewRoot("test")
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(io.Discard)
	root.SetArgs([]string{"scan", testdataPath(t, "jwt-fixture.log"), "--mask-pii=true"})

	err := root.Execute()
	var exit *ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if exit.Code != 2 {
		t.Fatalf("expected exit code 2, got %d", exit.Code)
	}
	if !strings.Contains(out.String(), "PII") {
		t.Fatalf("expected PII in terminal output, got:\n%s", out.String())
	}
}

func TestScan_JWTWithMaskPII_JSONOutput(t *testing.T) {
	root := NewRoot("test")
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"scan", testdataPath(t, "jwt-fixture.log"),
		"--mask-pii=true", "--output", "json",
	})

	err := root.Execute()
	var exit *ExitError
	if !errors.As(err, &exit) || exit.Code != 2 {
		t.Fatalf("expected exit 2, err=%v", err)
	}

	var rep output.Report
	if err := json.Unmarshal(out.Bytes(), &rep); err != nil {
		t.Fatalf("json decode: %v\nraw: %s", err, out.String())
	}
	if !rep.PIICritical {
		t.Error("expected pii_critical true")
	}
	if rep.PIIFindingCount == 0 {
		t.Error("expected pii findings")
	}
	if rep.Provider != "datadog" {
		t.Errorf("provider=%q", rep.Provider)
	}
}

func TestScan_JWTWithMaskPII_MarkdownOutput(t *testing.T) {
	root := NewRoot("test")
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"scan", testdataPath(t, "jwt-fixture.log"),
		"--mask-pii=true", "--output", "markdown",
	})

	err := root.Execute()
	var exit *ExitError
	if !errors.As(err, &exit) || exit.Code != 2 {
		t.Fatalf("expected exit 2, err=%v", err)
	}
	body := out.String()
	for _, want := range []string{"# telesieve scan report", "PII findings", "critical: true"} {
		if !strings.Contains(body, want) {
			t.Errorf("markdown missing %q\n%s", want, body)
		}
	}
}

func TestScan_PlainLog_Success(t *testing.T) {
	root := NewRoot("test")
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(io.Discard)
	root.SetArgs([]string{"scan", testdataPath(t, "plain.log"), "--mask-pii=false"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "telesieve scan report") {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
