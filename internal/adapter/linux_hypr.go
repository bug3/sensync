//go:build linux

package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bug3dev/sensync/internal/adapter/plan"
	"github.com/bug3dev/sensync/internal/config"
)

type HyprlandAdapter struct {
	exec     Executor
	confPath string // e.g. $XDG_CONFIG_HOME/hypr/sensync.conf
}

// NewHyprlandAdapter constructs an adapter wired to ShellExecutor and the
// default config path under $XDG_CONFIG_HOME (or ~/.config).
func NewHyprlandAdapter() (*HyprlandAdapter, error) {
	if _, err := exec.LookPath("hyprctl"); err != nil {
		return nil, fmt.Errorf("hyprctl not found in PATH; sensync currently supports Hyprland only on Linux")
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("locate home dir: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return &HyprlandAdapter{
		exec:     ShellExecutor{},
		confPath: filepath.Join(base, "hypr", "sensync.conf"),
	}, nil
}

func (a *HyprlandAdapter) Name() string { return "hyprland" }

func (a *HyprlandAdapter) BuildPlan(cfg config.Config) (Plan, error) {
	if err := os.MkdirAll(filepath.Dir(a.confPath), 0o755); err != nil {
		return Plan{}, fmt.Errorf("ensure conf dir: %w", err)
	}
	return plan.Hyprland(cfg, a.confPath)
}

func (a *HyprlandAdapter) Apply(p Plan, dryRun bool) (Result, error) {
	var res Result
	if dryRun {
		res.Applied = append(res.Applied, p.Steps...)
		return res, nil
	}
	for _, s := range p.Steps {
		if err := a.exec.Do(s); err != nil {
			res.Failed = append(res.Failed, FailedStep{Step: s, Err: err})
			continue
		}
		res.Applied = append(res.Applied, s)
	}
	return res, nil
}

func (a *HyprlandAdapter) Get() (config.Config, error) {
	cfg := config.Default()
	if v, err := readHyprctlFloat("input:sensitivity"); err == nil {
		cfg.Mouse.Sensitivity = v + 1.0
		cfg.Trackpad.Sensitivity = v + 1.0
	}
	// Acceleration, natural_scroll, and scroll_speed are not read in MVP; the
	// other config fields stay at Default() values so `get` output is still
	// valid TOML.
	return cfg, nil
}

// readHyprctlFloat invokes `hyprctl -j getoption KEY` and extracts the
// `"float": X` field via a small string scan. We avoid encoding/json to keep
// dependencies minimal; the format is stable across Hyprland versions.
func readHyprctlFloat(key string) (float64, error) {
	out, err := exec.Command("hyprctl", "-j", "getoption", key).Output()
	if err != nil {
		return 0, fmt.Errorf("hyprctl getoption %s: %w", key, err)
	}
	raw := strings.ReplaceAll(string(out), "\n", " ")
	idx := strings.Index(raw, `"float":`)
	if idx < 0 {
		return 0, fmt.Errorf("hyprctl json: missing float field in %q", raw)
	}
	rest := raw[idx+len(`"float":`):]
	end := strings.IndexAny(rest, ",}")
	if end < 0 {
		return 0, fmt.Errorf("hyprctl json: malformed float field in %q", raw)
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(rest[:end]), 64)
	if err != nil {
		return 0, fmt.Errorf("hyprctl json float parse: %w", err)
	}
	return v, nil
}
