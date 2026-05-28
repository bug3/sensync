//go:build darwin

package adapter

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bug3dev/sensync/internal/adapter/plan"
	"github.com/bug3dev/sensync/internal/config"
)

type MacOSAdapter struct {
	exec Executor
}

func NewMacOSAdapter() (*MacOSAdapter, error) {
	if _, err := exec.LookPath("defaults"); err != nil {
		return nil, fmt.Errorf("`defaults` not found in PATH; not a macOS host?")
	}
	return &MacOSAdapter{exec: ShellExecutor{}}, nil
}

func (a *MacOSAdapter) Name() string { return "macos" }

func (a *MacOSAdapter) BuildPlan(cfg config.Config) (Plan, error) {
	return plan.MacOS(cfg)
}

func (a *MacOSAdapter) Apply(p Plan, dryRun bool) (Result, error) {
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

func (a *MacOSAdapter) Get() (config.Config, error) {
	cfg := config.Default()
	if v, err := readDefaultsFloat("com.apple.mouse.scaling"); err == nil {
		switch {
		case v == -1.0:
			cfg.Mouse.Sensitivity = 1.0
			cfg.Mouse.Acceleration = false
		case v == 0.0:
			cfg.Mouse.Sensitivity = 1.0
			cfg.Mouse.Acceleration = true
		default:
			cfg.Mouse.Sensitivity = (v / 3.0) + 1.0
			cfg.Mouse.Acceleration = true
		}
	}
	if v, err := readDefaultsFloat("com.apple.trackpad.scaling"); err == nil {
		switch {
		case v == -1.0:
			cfg.Trackpad.Sensitivity = 1.0
			cfg.Trackpad.Acceleration = false
		case v == 0.0:
			cfg.Trackpad.Sensitivity = 1.0
			cfg.Trackpad.Acceleration = true
		default:
			cfg.Trackpad.Sensitivity = (v / 3.0) + 1.0
			cfg.Trackpad.Acceleration = true
		}
	}
	if v, err := readDefaultsBool("com.apple.swipescrolldirection"); err == nil {
		cfg.Mouse.NaturalScroll = v
		cfg.Trackpad.NaturalScroll = v
	}
	if v, err := readDefaultsFloat("com.apple.scrollwheel.scaling"); err == nil {
		cfg.Mouse.ScrollSpeed = v
	}
	return cfg, nil
}

func readDefaultsFloat(key string) (float64, error) {
	out, err := exec.Command("defaults", "read", "-g", key).Output()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}

func readDefaultsBool(key string) (bool, error) {
	out, err := exec.Command("defaults", "read", "-g", key).Output()
	if err != nil {
		return false, err
	}
	switch strings.TrimSpace(string(out)) {
	case "1", "true", "YES":
		return true, nil
	case "0", "false", "NO":
		return false, nil
	}
	return false, fmt.Errorf("unparseable bool from defaults: %q", string(out))
}
