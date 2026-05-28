package mapping

// HyprlandSensitivityFromMultiplier maps a unified multiplier to Hyprland's
// `input { sensitivity = X }` value, which is libinput's accel speed setpoint
// in the range [-1.0, 1.0] when accel_profile = flat.
//
// 1.0 (unity) -> 0.0 (Hyprland default, raw 1:1 when paired with flat profile).
// Other values are offset around zero and clamped.
func HyprlandSensitivityFromMultiplier(m float64) float64 {
	v := m - 1.0
	if v < -1.0 {
		return -1.0
	}
	if v > 1.0 {
		return 1.0
	}
	return v
}
