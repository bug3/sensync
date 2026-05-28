package plan

import (
	"strings"
	"testing"

	"github.com/bug3dev/sensync/internal/adapter"
	"github.com/bug3dev/sensync/internal/config"
)

func TestHyprlandUnityPlan(t *testing.T) {
	cfg := config.Default()
	got, err := Hyprland(cfg, "/home/u/.config/hypr/sensync.conf")
	if err != nil {
		t.Fatalf("Hyprland: %v", err)
	}

	var hasFileWrite, hasSensCmd, hasAccelCmd bool
	for _, s := range got.Steps {
		switch {
		case s.Kind == adapter.StepWriteFile && s.Target == "/home/u/.config/hypr/sensync.conf":
			hasFileWrite = true
			if !strings.Contains(s.Args[0], "sensitivity = 0") {
				t.Errorf("expected `sensitivity = 0` in file, got:\n%s", s.Args[0])
			}
			if !strings.Contains(s.Args[0], "accel_profile = flat") {
				t.Errorf("expected `accel_profile = flat` in file, got:\n%s", s.Args[0])
			}
		case s.Kind == adapter.StepExec && s.Target == "hyprctl" && stepArgsContain(s, "input:sensitivity"):
			hasSensCmd = true
		case s.Kind == adapter.StepExec && s.Target == "hyprctl" && stepArgsContain(s, "input:accel_profile"):
			hasAccelCmd = true
		}
	}
	if !hasFileWrite {
		t.Error("missing file-write step for sensync.conf")
	}
	if !hasSensCmd {
		t.Error("missing hyprctl input:sensitivity step")
	}
	if !hasAccelCmd {
		t.Error("missing hyprctl input:accel_profile step")
	}
}

func TestHyprlandWarnsWhenMouseAndTrackpadDiffer(t *testing.T) {
	cfg := config.Default()
	cfg.Trackpad.Sensitivity = 1.5
	got, err := Hyprland(cfg, "/tmp/x.conf")
	if err != nil {
		t.Fatalf("Hyprland: %v", err)
	}
	if len(got.Warnings) == 0 {
		t.Error("expected a warning about global Hyprland sensitivity, got none")
	}
}

func stepArgsContain(s adapter.Step, needle string) bool {
	for _, a := range s.Args {
		if strings.Contains(a, needle) {
			return true
		}
	}
	return false
}
