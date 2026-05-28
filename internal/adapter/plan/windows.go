package plan

import (
	"fmt"
	"strconv"

	"github.com/bug3/sensync/internal/adapter/types"
	"github.com/bug3/sensync/internal/config"
	"github.com/bug3/sensync/internal/mapping"
)

const (
	regMouse          = `HKCU\Control Panel\Mouse`
	regDesktop        = `HKCU\Control Panel\Desktop`
	regPrecisionTouch = `HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\PrecisionTouchpad`
)

// Windows builds a Plan for a Windows host. All registry writes target HKCU,
// so no elevation is needed. After registry changes are written, a
// SystemParametersInfo broadcast (SPI_SETMOUSE = 0x0004) tells the running
// session to pick up the new values without logout.
func Windows(cfg config.Config) (types.Plan, error) {
	var warnings []string

	if cfg.Mouse.NaturalScroll || cfg.Trackpad.NaturalScroll {
		warnings = append(warnings,
			"windows: natural_scroll requires per-device registry edits and is not implemented in MVP; setting will be skipped")
	}

	accel := cfg.Mouse.Acceleration
	if cfg.Mouse.Sensitivity != 1.0 && !accel {
		warnings = append(warnings,
			"windows: acceleration=false is only honored at sensitivity=1.0; the OS curve will be used for non-unity sensitivity")
	}

	sensVal := mapping.MarkCFromMultiplier(cfg.Mouse.Sensitivity)
	mouseSpeed := "0"
	if accel || cfg.Mouse.Sensitivity != 1.0 {
		mouseSpeed = "1"
	}

	steps := []types.Step{
		regSet(regMouse, "MouseSensitivity", "REG_SZ", strconv.Itoa(sensVal)),
		regSet(regMouse, "MouseSpeed", "REG_SZ", mouseSpeed),
		regSet(regMouse, "MouseThreshold1", "REG_SZ", "0"),
		regSet(regMouse, "MouseThreshold2", "REG_SZ", "0"),
		regSet(regDesktop, "WheelScrollLines", "REG_SZ", scrollLinesFromSpeed(cfg.Mouse.ScrollSpeed)),
		regSet(regPrecisionTouch, "Sensitivity", "REG_DWORD", trackpadSensValue(cfg.Trackpad.Sensitivity)),
		{
			Kind:   types.StepSysCall,
			Target: "SPI_SETMOUSE",
			Args:   nil,
			Desc:   "SystemParametersInfo(SPI_SETMOUSE) to broadcast mouse changes",
		},
	}
	return types.Plan{Steps: steps, Warnings: warnings}, nil
}

func regSet(path, name, typ, value string) types.Step {
	return types.Step{
		Kind:   types.StepRegSet,
		Target: path,
		Args:   []string{name, typ, value},
		Desc:   fmt.Sprintf(`reg set %s\%s (%s) = %s`, path, name, typ, value),
	}
}

func scrollLinesFromSpeed(s float64) string {
	// Default Windows is 3 lines per notch. Multiplier scales linearly.
	lines := int(s * 3.0)
	if lines < 1 {
		lines = 1
	}
	if lines > 100 {
		lines = 100
	}
	return strconv.Itoa(lines)
}

// trackpadSensValue maps the unified multiplier to the Precision Touchpad
// `Sensitivity` DWORD. Microsoft exposes documented buckets:
//
//	0 = Most sensitive, 1 = High sensitivity, 2 = Medium sensitivity,
//	3 = Low sensitivity, 4 = Most low sensitivity.
//
// Unity multiplier maps to 2 (medium), with one bucket per ±0.5.
func trackpadSensValue(m float64) string {
	bucket := 2 - int((m-1.0)*2.0)
	if bucket < 0 {
		bucket = 0
	}
	if bucket > 4 {
		bucket = 4
	}
	return strconv.Itoa(bucket)
}
