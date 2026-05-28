//go:build windows

package adapter

import (
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/bug3/sensync/internal/adapter/plan"
	"github.com/bug3/sensync/internal/config"
)

const (
	spiSetMouse       uintptr = 0x0004
	spifSendChange    uintptr = 0x0002
	spifUpdateIniFile uintptr = 0x0001
)

type WindowsAdapter struct {
	exec Executor // tests inject RecordingExecutor; production uses nativeWindowsExecutor
}

func NewWindowsAdapter() (*WindowsAdapter, error) {
	return &WindowsAdapter{exec: nativeWindowsExecutor{}}, nil
}

func (a *WindowsAdapter) Name() string { return "windows" }

func (a *WindowsAdapter) BuildPlan(cfg config.Config) (Plan, error) {
	return plan.Windows(cfg)
}

func (a *WindowsAdapter) Apply(p Plan, dryRun bool) (Result, error) {
	var res Result
	if dryRun {
		res.Applied = append(res.Applied, p.Steps...)
		return res, nil
	}
	for _, s := range p.Steps {
		if err := a.exec.Do(s); err != nil {
			res.Failed = append(res.Failed, FailedStep{Step: s, Err: err})
			continue
		}
		res.Applied = append(res.Applied, s)
	}
	return res, nil
}

func (a *WindowsAdapter) Get() (config.Config, error) {
	cfg := config.Default()
	if v, err := readRegistryString(`Control Panel\Mouse`, "MouseSensitivity"); err == nil {
		if n, err2 := strconv.Atoi(v); err2 == nil {
			m := float64(n) / 10.0
			if m < 0.1 {
				m = 0.1
			}
			if m > 3.0 {
				m = 3.0
			}
			cfg.Mouse.Sensitivity = m
		}
	}
	if v, err := readRegistryString(`Control Panel\Mouse`, "MouseSpeed"); err == nil {
		cfg.Mouse.Acceleration = v != "0"
	}
	if v, err := readRegistryString(`Control Panel\Desktop`, "WheelScrollLines"); err == nil {
		if n, err2 := strconv.Atoi(v); err2 == nil && n > 0 {
			s := float64(n) / 3.0 // inverse of scrollLinesFromSpeed
			if s < 0.1 {
				s = 0.1
			}
			if s > 5.0 {
				s = 5.0
			}
			cfg.Mouse.ScrollSpeed = s
		}
	}
	if n, err := readRegistryDWord(`SOFTWARE\Microsoft\Windows\CurrentVersion\PrecisionTouchpad`, "Sensitivity"); err == nil {
		// Inverse of trackpadSensValue: bucket 0 = most sensitive, 4 = least.
		// Unity multiplier (1.0) maps to bucket 2.
		m := 1.0 - (float64(n)-2.0)/2.0
		if m < 0.1 {
			m = 0.1
		}
		if m > 3.0 {
			m = 3.0
		}
		cfg.Trackpad.Sensitivity = m
	}
	// natural_scroll is not read on Windows; MVP does not write it either.
	return cfg, nil
}

// nativeWindowsExecutor performs the actual registry writes and SPI broadcast.
type nativeWindowsExecutor struct{}

func (nativeWindowsExecutor) Do(s Step) error {
	switch s.Kind {
	case StepRegSet:
		return writeRegistryValue(s)
	case StepSysCall:
		return broadcastMouseChange()
	default:
		// File writes and Exec aren't expected on Windows for MVP. Defer to shell.
		return ShellExecutor{}.Do(s)
	}
}

func broadcastMouseChange() error {
	user32 := windows.NewLazySystemDLL("user32.dll")
	proc := user32.NewProc("SystemParametersInfoW")
	ret, _, errno := proc.Call(spiSetMouse, 0, uintptr(unsafe.Pointer(nil)), spifSendChange|spifUpdateIniFile)
	if ret == 0 {
		return fmt.Errorf("SystemParametersInfoW: %w", errno)
	}
	return nil
}

func readRegistryString(subkey, name string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, subkey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	v, _, err := k.GetStringValue(name)
	return v, err
}

func readRegistryDWord(subkey, name string) (uint32, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, subkey, registry.QUERY_VALUE)
	if err != nil {
		return 0, err
	}
	defer k.Close()
	v, _, err := k.GetIntegerValue(name)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}
