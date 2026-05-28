package mapping

import "testing"

func TestMacOSScalingAtUnityAccelOff(t *testing.T) {
	if got := MacOSScalingFromMultiplier(1.0, false); got != -1.0 {
		t.Errorf("1.0/off -> %v, want -1.0", got)
	}
}

func TestMacOSScalingAtUnityAccelOn(t *testing.T) {
	if got := MacOSScalingFromMultiplier(1.0, true); got != 0.0 {
		t.Errorf("1.0/on -> %v, want 0.0", got)
	}
}

func TestMacOSScalingAboveUnity(t *testing.T) {
	if got := MacOSScalingFromMultiplier(1.5, true); got != 1.5 {
		t.Errorf("1.5 -> %v, want 1.5", got)
	}
}

func TestMacOSScalingClampsHigh(t *testing.T) {
	if got := MacOSScalingFromMultiplier(5.0, true); got != 3.0 {
		t.Errorf("5.0 -> %v, want 3.0", got)
	}
}

func TestMacOSScalingClampsLow(t *testing.T) {
	if got := MacOSScalingFromMultiplier(0.5, true); got != 0.0 {
		t.Errorf("0.5 (below unity, accel-on path) -> %v, want 0.0", got)
	}
}
