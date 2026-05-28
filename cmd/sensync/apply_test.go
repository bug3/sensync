package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/bug3/sensync/internal/adapter"
	"github.com/bug3/sensync/internal/adapter/types"
	"github.com/bug3/sensync/internal/config"
)

func TestResolveConfigPathExplicit(t *testing.T) {
	got, err := resolveConfigPath("/explicit/path.toml")
	if err != nil {
		t.Fatalf("resolveConfigPath: %v", err)
	}
	if got != "/explicit/path.toml" {
		t.Errorf("got %q, want explicit path", got)
	}
}

func TestResolveConfigPathPicksCwd(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "sensync.toml")
	if err := os.WriteFile(cfg, []byte("version=1\n"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}
	prev, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	got, err := resolveConfigPath("")
	if err != nil {
		t.Fatalf("resolveConfigPath: %v", err)
	}
	if got != "sensync.toml" {
		t.Errorf("got %q, want sensync.toml", got)
	}
}

func TestResolveConfigPathMissing(t *testing.T) {
	dir := t.TempDir()
	prev, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	// Force UserConfigDir to a definitely-empty location so the fallback misses.
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, "no-such-dir"))
	t.Setenv("HOME", filepath.Join(dir, "no-home"))
	t.Setenv("APPDATA", filepath.Join(dir, "no-appdata"))
	if _, err := resolveConfigPath(""); err == nil {
		t.Fatal("expected error when no config exists, got nil")
	}
}

func TestNeedsTrackpadPrompt(t *testing.T) {
	base := config.Default()
	if needsTrackpadPrompt(base) {
		t.Error("default config should not need prompt")
	}
	diverged := base
	diverged.Trackpad.Sensitivity = 1.5
	if !needsTrackpadPrompt(diverged) {
		t.Error("diverged sensitivity should need prompt")
	}
	diverged = base
	diverged.Trackpad.Acceleration = true
	if !needsTrackpadPrompt(diverged) {
		t.Error("diverged acceleration should need prompt")
	}
}

func TestPrintApplySummaryWritesFailuresToStderr(t *testing.T) {
	cmd := &cobra.Command{}
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	res := adapter.Result{
		Applied: []types.Step{{Desc: "ok"}},
		Failed: []types.FailedStep{
			{Step: types.Step{Desc: "fail1"}, Err: errors.New("boom")},
		},
	}
	printApplySummary(cmd, "hyprland", res)
	stdout := out.String()
	stderr := errBuf.String()
	if !strings.Contains(stdout, "[hyprland] applied 1 step(s), failed 1") {
		t.Errorf("expected summary line on stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "fail1: boom") {
		t.Errorf("expected failure detail on stderr, got %q", stderr)
	}
}

func TestPrintApplySummaryMacOSReminder(t *testing.T) {
	cmd := &cobra.Command{}
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	printApplySummary(cmd, "macos", adapter.Result{Applied: []types.Step{{Desc: "x"}}})
	if !strings.Contains(out.String(), "log out / restart") {
		t.Errorf("expected macOS logout reminder, got %q", out.String())
	}
	printApplySummary(cmd, "hyprland", adapter.Result{Applied: []types.Step{{Desc: "x"}}})
	// Reminder appears once for macos, not for other adapter names.
	if strings.Count(out.String(), "log out / restart") != 1 {
		t.Errorf("logout reminder should be macOS-only, got %q", out.String())
	}
}

func TestCLIErrorEmptyMessage(t *testing.T) {
	ce := cliError{code: 3}
	if got := ce.Error(); got != "" {
		t.Errorf("nil-err cliError should produce empty string, got %q", got)
	}
}
