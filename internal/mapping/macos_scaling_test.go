package mapping

import (
	"math"
	"testing"
)

func approxEqual(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func TestMacOSScalingAtUnityAccelOff(t *testing.T) {
	if got := MacOSScalingFromMultiplier(1.0, false); got != -1.0 {
		t.Errorf("1.0/off -> %v, want -1.0", got)
	}
}

func TestMacOSScalingAtUnityAccelOnIsMacOSDefault(t *testing.T) {
	// At sensitivity=1.0 with acceleration on, we anchor to macOS's
	// trackpad default (0.875) so users get a usable feel, not the
	// slowest slider position.
	if got := MacOSScalingFromMultiplier(1.0, true); !approxEqual(got, 0.875, 1e-9) {
		t.Errorf("1.0/on -> %v, want 0.875", got)
	}
}

func TestMacOSScalingAtMin(t *testing.T) {
	if got := MacOSScalingFromMultiplier(0.1, true); !approxEqual(got, 0.0, 1e-9) {
		t.Errorf("0.1 -> %v, want 0.0", got)
	}
}

func TestMacOSScalingAtMax(t *testing.T) {
	if got := MacOSScalingFromMultiplier(3.0, true); !approxEqual(got, 3.0, 1e-9) {
		t.Errorf("3.0 -> %v, want 3.0", got)
	}
}

func TestMacOSScalingMidBelowUnity(t *testing.T) {
	// Piecewise linear from (0.1, 0) to (1.0, 0.875):
	// at 0.5 -> (0.5-0.1)/0.9 * 0.875 = 0.3889
	if got := MacOSScalingFromMultiplier(0.5, true); !approxEqual(got, 0.3889, 1e-3) {
		t.Errorf("0.5 -> %v, want ~0.3889", got)
	}
}

func TestMacOSScalingMidAboveUnity(t *testing.T) {
	// Piecewise linear from (1.0, 0.875) to (3.0, 3.0):
	// at 2.0 -> 0.875 + 1.0/2.0 * 2.125 = 1.9375
	if got := MacOSScalingFromMultiplier(2.0, true); !approxEqual(got, 1.9375, 1e-3) {
		t.Errorf("2.0 -> %v, want ~1.9375", got)
	}
}

func TestMacOSScalingClampsHigh(t *testing.T) {
	if got := MacOSScalingFromMultiplier(5.0, true); got != 3.0 {
		t.Errorf("5.0 -> %v, want 3.0", got)
	}
}

func TestMacOSScalingClampsLow(t *testing.T) {
	if got := MacOSScalingFromMultiplier(0.05, true); got != 0.0 {
		t.Errorf("0.05 -> %v, want 0.0", got)
	}
}

func TestMacOSMultiplierFromScalingRawIsUnityAccelOff(t *testing.T) {
	m, accel := MacOSMultiplierFromScaling(-1.0)
	if m != 1.0 || accel {
		t.Errorf("-1 -> (%v, %v), want (1.0, false)", m, accel)
	}
}

func TestMacOSMultiplierFromScalingDefaultIsUnity(t *testing.T) {
	m, accel := MacOSMultiplierFromScaling(0.875)
	if !approxEqual(m, 1.0, 1e-6) || !accel {
		t.Errorf("0.875 -> (%v, %v), want (1.0, true)", m, accel)
	}
}

func TestMacOSMultiplierFromScalingMin(t *testing.T) {
	m, accel := MacOSMultiplierFromScaling(0.0)
	if !approxEqual(m, 0.1, 1e-6) || !accel {
		t.Errorf("0.0 -> (%v, %v), want (0.1, true)", m, accel)
	}
}

func TestMacOSMultiplierFromScalingMax(t *testing.T) {
	m, accel := MacOSMultiplierFromScaling(3.0)
	if !approxEqual(m, 3.0, 1e-6) || !accel {
		t.Errorf("3.0 -> (%v, %v), want (3.0, true)", m, accel)
	}
}

func TestMacOSScalingRoundTrip(t *testing.T) {
	cases := []float64{0.1, 0.5, 0.875, 1.0, 1.5, 2.0, 2.5, 3.0}
	for _, want := range cases {
		v := MacOSScalingFromMultiplier(want, true)
		got, _ := MacOSMultiplierFromScaling(v)
		if !approxEqual(got, want, 1e-6) {
			t.Errorf("round-trip %v -> %v -> %v", want, v, got)
		}
	}
}
