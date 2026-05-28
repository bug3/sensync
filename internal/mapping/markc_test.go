package mapping

import "testing"

func TestMarkCAtUnity(t *testing.T) {
	if got := MarkCFromMultiplier(1.0); got != 10 {
		t.Errorf("1.0 -> %d, want 10", got)
	}
}

func TestMarkCClampsLow(t *testing.T) {
	if got := MarkCFromMultiplier(0.1); got < 1 || got > 1 {
		t.Errorf("0.1 -> %d, want 1", got)
	}
}

func TestMarkCClampsHigh(t *testing.T) {
	if got := MarkCFromMultiplier(3.0); got != 20 {
		t.Errorf("3.0 -> %d, want 20", got)
	}
}

func TestMarkCMonotonic(t *testing.T) {
	prev := MarkCFromMultiplier(0.1)
	for _, v := range []float64{0.5, 0.8, 1.0, 1.2, 1.5, 2.0, 2.5, 3.0} {
		got := MarkCFromMultiplier(v)
		if got < prev {
			t.Errorf("non-monotonic at %v: prev=%d got=%d", v, prev, got)
		}
		prev = got
	}
}
