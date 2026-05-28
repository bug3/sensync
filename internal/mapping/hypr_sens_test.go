package mapping

import "testing"

func TestHyprlandAtUnity(t *testing.T) {
	if got := HyprlandSensitivityFromMultiplier(1.0); got != 0.0 {
		t.Errorf("1.0 -> %v, want 0.0", got)
	}
}

func TestHyprlandAbove(t *testing.T) {
	if got := HyprlandSensitivityFromMultiplier(1.5); got != 0.5 {
		t.Errorf("1.5 -> %v, want 0.5", got)
	}
}

func TestHyprlandBelow(t *testing.T) {
	if got := HyprlandSensitivityFromMultiplier(0.5); got != -0.5 {
		t.Errorf("0.5 -> %v, want -0.5", got)
	}
}

func TestHyprlandClampsHigh(t *testing.T) {
	if got := HyprlandSensitivityFromMultiplier(3.0); got != 1.0 {
		t.Errorf("3.0 -> %v, want 1.0", got)
	}
}

func TestHyprlandClampsLow(t *testing.T) {
	if got := HyprlandSensitivityFromMultiplier(0.0); got != -1.0 {
		t.Errorf("0.0 -> %v, want -1.0", got)
	}
}
