//go:build windows

package adapter

import (
	"fmt"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/bug3dev/sensync/internal/adapter/plan"
	"github.com/bug3dev/sensync/internal/config"
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
			cfg.Mouse.Sensitivity = float64(n) / 10.0
		}
	}
	if v, err := readRegistryString(`Control Panel\Mouse`, "MouseSpeed"); err == nil {
		cfg.Mouse.Acceleration = v != "0"
	}
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
