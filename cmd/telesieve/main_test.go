package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func buildBinary(t *testing.T, dir string) string {
	t.Helper()
	bin := filepath.Join(dir, "telesieve")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(repoRoot(t), "cmd", "telesieve")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return bin
}

func TestCLI_E2E_ExitCodes(t *testing.T) {
	bin := buildBinary(t, t.TempDir())
	jwtLog := filepath.Join(repoRoot(t), "testdata", "jwt-fixture.log")
	missing := filepath.Join(t.TempDir(), "missing.log")

	t.Run("jwt_mask_pii_exit_2", func(t *testing.T) {
		cmd := exec.Command(bin, "scan", jwtLog, "--mask-pii=true", "--output=json")
		err := cmd.Run()
		if code := exitCode(err); code != 2 {
			t.Fatalf("exit=%d err=%v", code, err)
		}
	})

	t.Run("missing_file_exit_1", func(t *testing.T) {
		cmd := exec.Command(bin, "scan", missing)
		err := cmd.Run()
		if code := exitCode(err); code != 1 {
			t.Fatalf("exit=%d err=%v", code, err)
		}
	})
}

func TestCLI_E2E_Version(t *testing.T) {
	bin := buildBinary(t, t.TempDir())
	out, err := exec.Command(bin, "version").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "telesieve") {
		t.Fatalf("output: %s", out)
	}
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}
