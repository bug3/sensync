package mapping

// macOS exposes `com.apple.{mouse,trackpad}.scaling` as a single knob:
//
//   -1     : acceleration disabled (raw 1:1)
//    0..3  : Apple's UI tracking-speed slider; 0 is the slowest position,
//            ~0.875 is the trackpad default, and 3.0 is the maximum.
//
// We anchor sensync's unity multiplier (sensitivity=1.0) to the macOS
// trackpad default of 0.875 so a fresh config produces a usable feel
// out of the box, and we map the full sensync range [0.1, 3.0] to the
// full macOS range [0, 3] piecewise-linearly through that anchor.
const macOSDefaultScaling = 0.875

// MacOSScalingFromMultiplier maps the unified multiplier to the
// `com.apple.{mouse,trackpad}.scaling` `defaults` value.
//
// At unity multiplier with acceleration disabled we return -1 so macOS
// uses raw 1:1 input; with acceleration enabled we return the macOS
// trackpad default (0.875). Non-unity multipliers always engage the
// macOS acceleration curve regardless of the `acceleration` config field;
// the caller is responsible for warning the user.
func MacOSScalingFromMultiplier(m float64, accel bool) float64 {
	if m == 1.0 && !accel {
		return -1.0
	}
	if m <= 0.1 {
		return 0.0
	}
	if m >= 3.0 {
		return 3.0
	}
	if m <= 1.0 {
		return (m - 0.1) / (1.0 - 0.1) * macOSDefaultScaling
	}
	return macOSDefaultScaling + (m-1.0)/(3.0-1.0)*(3.0-macOSDefaultScaling)
}

// MacOSMultiplierFromScaling inverts MacOSScalingFromMultiplier so
// `sensync get` round-trips through `sensync apply`. A scaling of -1
// reports raw 1:1 (acceleration off); every other value reports
// acceleration on.
func MacOSMultiplierFromScaling(v float64) (multiplier float64, acceleration bool) {
	if v == -1.0 {
		return 1.0, false
	}
	if v <= 0.0 {
		return 0.1, true
	}
	if v >= 3.0 {
		return 3.0, true
	}
	if v <= macOSDefaultScaling {
		return 0.1 + v/macOSDefaultScaling*(1.0-0.1), true
	}
	return 1.0 + (v-macOSDefaultScaling)/(3.0-macOSDefaultScaling)*(3.0-1.0), true
}
