package plan

import (
	"fmt"
	"strconv"

	"github.com/bug3/sensync/internal/adapter/types"
	"github.com/bug3/sensync/internal/config"
	"github.com/bug3/sensync/internal/mapping"
)

// MacOS builds a Plan for a macOS host. All changes go through `defaults`;
// changes apply to new processes, so the caller should print a "log out for
// full effect" reminder after Apply succeeds.
func MacOS(cfg config.Config) (types.Plan, error) {
	var warnings []string

	if cfg.Mouse.NaturalScroll != cfg.Trackpad.NaturalScroll {
		warnings = append(warnings,
			"macos: com.apple.swipescrolldirection is a single global setting; trackpad and mouse natural_scroll cannot diverge. Using trackpad value.")
	}
	if cfg.Mouse.Sensitivity != 1.0 && !cfg.Mouse.Acceleration {
		warnings = append(warnings,
			"macos: acceleration=false is only honored at sensitivity=1.0; non-unity sensitivity uses the OS acceleration curve")
	}

	mouseScale := mapping.MacOSScalingFromMultiplier(cfg.Mouse.Sensitivity, cfg.Mouse.Acceleration)
	trackpadScale := mapping.MacOSScalingFromMultiplier(cfg.Trackpad.Sensitivity, cfg.Trackpad.Acceleration)

	steps := []types.Step{
		defaultsFloat("com.apple.mouse.scaling", mouseScale),
		defaultsFloat("com.apple.trackpad.scaling", trackpadScale),
		defaultsBool("com.apple.swipescrolldirection", cfg.Trackpad.NaturalScroll),
		defaultsFloat("com.apple.scrollwheel.scaling", cfg.Mouse.ScrollSpeed),
	}
	return types.Plan{Steps: steps, Warnings: warnings}, nil
}

func defaultsFloat(key string, val float64) types.Step {
	v := strconv.FormatFloat(val, 'f', -1, 64)
	return types.Step{
		Kind:   types.StepExec,
		Target: "defaults",
		Args:   []string{"write", "-g", key, "-float", v},
		Desc:   fmt.Sprintf("defaults write -g %s -float %s", key, v),
	}
}

func defaultsBool(key string, val bool) types.Step {
	v := strconv.FormatBool(val)
	return types.Step{
		Kind:   types.StepExec,
		Target: "defaults",
		Args:   []string{"write", "-g", key, "-bool", v},
		Desc:   fmt.Sprintf("defaults write -g %s -bool %s", key, v),
	}
}
