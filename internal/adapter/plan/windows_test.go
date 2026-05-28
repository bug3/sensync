package plan

import (
	"strings"
	"testing"

	"github.com/bug3dev/sensync/internal/adapter/types"
	"github.com/bug3dev/sensync/internal/config"
)

func TestWindowsUnityPlanSetsSensitivity10AndAccelOff(t *testing.T) {
	cfg := config.Default() // sensitivity=1.0, acceleration=false
	got, err := Windows(cfg)
	if err != nil {
		t.Fatalf("Windows: %v", err)
	}
	var sawSens10, sawSpeed0, sawThr1Zero, sawThr2Zero, sawSysCall bool
	for _, s := range got.Steps {
		if s.Kind == types.StepRegSet {
			joined := strings.Join(s.Args, " ")
			switch {
			case strings.Contains(joined, "MouseSensitivity") && strings.Contains(joined, "10"):
				sawSens10 = true
			case strings.Contains(joined, "MouseSpeed") && strings.HasSuffix(joined, " 0"):
				sawSpeed0 = true
			case strings.Contains(joined, "MouseThreshold1") && strings.HasSuffix(joined, " 0"):
				sawThr1Zero = true
			case strings.Contains(joined, "MouseThreshold2") && strings.HasSuffix(joined, " 0"):
				sawThr2Zero = true
			}
		}
		if s.Kind == types.StepSysCall && strings.Contains(s.Target, "SPI_SETMOUSE") {
			sawSysCall = true
		}
	}
	if !sawSens10 {
		t.Error("missing MouseSensitivity=10")
	}
	if !sawSpeed0 {
		t.Error("missing MouseSpeed=0")
	}
	if !sawThr1Zero || !sawThr2Zero {
		t.Error("missing MouseThreshold1/2=0")
	}
	if !sawSysCall {
		t.Error("missing SPI_SETMOUSE broadcast step")
	}
}

func TestWindowsWarnsAboutNaturalScroll(t *testing.T) {
	cfg := config.Default()
	cfg.Mouse.NaturalScroll = true
	got, _ := Windows(cfg)
	if len(got.Warnings) == 0 {
		t.Error("expected a warning about Windows natural_scroll being unsupported in MVP")
	}
}
