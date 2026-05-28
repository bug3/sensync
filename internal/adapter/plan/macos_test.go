package plan

import (
	"strings"
	"testing"

	"github.com/bug3dev/sensync/internal/adapter"
	"github.com/bug3dev/sensync/internal/config"
)

func TestMacOSUnityPlanWritesScalingMinusOne(t *testing.T) {
	cfg := config.Default() // sensitivity=1.0, acceleration=false
	got, err := MacOS(cfg)
	if err != nil {
		t.Fatalf("MacOS: %v", err)
	}
	var sawMouseScaling bool
	for _, s := range got.Steps {
		if s.Kind != adapter.StepExec || s.Target != "defaults" {
			continue
		}
		joined := strings.Join(s.Args, " ")
		if strings.Contains(joined, "com.apple.mouse.scaling") && strings.Contains(joined, "-1") {
			sawMouseScaling = true
		}
	}
	if !sawMouseScaling {
		t.Errorf("expected `defaults write -g com.apple.mouse.scaling -float -1` step, plan was:\n%+v", got.Steps)
	}
}

func TestMacOSWarnsWhenNaturalScrollDiffers(t *testing.T) {
	cfg := config.Default()
	cfg.Mouse.NaturalScroll = false
	cfg.Trackpad.NaturalScroll = true
	got, err := MacOS(cfg)
	if err != nil {
		t.Fatalf("MacOS: %v", err)
	}
	found := false
	for _, w := range got.Warnings {
		if strings.Contains(w, "swipescrolldirection") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning about global swipescrolldirection, got: %v", got.Warnings)
	}
}
