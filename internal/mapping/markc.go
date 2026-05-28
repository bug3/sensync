package mapping

import "math"

// MarkCFromMultiplier maps a unified sensitivity multiplier to the
// Windows HKCU\Control Panel\Mouse\MouseSensitivity value (1..20).
//
// Reference: HKCU MouseSensitivity is a documented lookup table; 10 is the
// only value at which Windows applies pure 1:1 input. Other values use the
// internal acceleration curve. We compute a linear scaling around 10 and
// clamp to the documented range.
func MarkCFromMultiplier(m float64) int {
	v := int(math.Round(10.0 * m))
	if v < 1 {
		return 1
	}
	if v > 20 {
		return 20
	}
	return v
}
