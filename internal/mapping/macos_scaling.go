package mapping

// MacOSScalingFromMultiplier maps the unified multiplier to the
// `com.apple.mouse.scaling` (or `com.apple.trackpad.scaling`) `defaults` value.
//
// Range and meaning per Apple:
//   -1     : acceleration disabled (raw 1:1)
//    0     : default acceleration curve, anchored at 1:1
//    0..3  : amplified curves; 3 is the maximum macOS exposes via defaults
//
// At unity multiplier (1.0), the result depends on whether acceleration is
// requested: -1 (off) or 0 (on). Non-unity multipliers cannot reproduce raw
// 1:1, so the acceleration curve is implicitly engaged regardless of the
// `acceleration` config field; the caller is responsible for the warning.
func MacOSScalingFromMultiplier(m float64, accel bool) float64 {
	if m == 1.0 {
		if accel {
			return 0.0
		}
		return -1.0
	}
	v := (m - 1.0) * 3.0
	if v < 0.0 {
		return 0.0
	}
	if v > 3.0 {
		return 3.0
	}
	return v
}
