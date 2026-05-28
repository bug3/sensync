// Package types holds the value types shared between adapter implementations
// and the plan generators. It is dependency-free of the surrounding adapter
// package, which lets `adapter` import `adapter/plan` (and `adapter/plan`
// import `adapter/types`) without creating an import cycle.
package types

// StepKind classifies a single change a Plan will perform.
type StepKind string

const (
	StepWriteFile StepKind = "WRITE_FILE" // write to a file (e.g., Hyprland source file)
	StepExec      StepKind = "EXEC"       // shell-out (e.g., hyprctl, defaults)
	StepRegSet    StepKind = "REG_SET"    // Windows registry value
	StepSysCall   StepKind = "SYS_CALL"   // Windows SystemParametersInfo broadcast
)

// Step is a single declarative change. Plan generators emit Steps; adapter
// executors perform them.
type Step struct {
	Kind   StepKind
	Target string   // file path, command name, registry path, syscall name
	Args   []string // command args, registry [name, type, value], etc.
	Desc   string   // human-readable line for --dry-run output
}

// Plan is a flat ordered list of Steps for one Apply.
type Plan struct {
	Steps    []Step
	Warnings []string // non-fatal, printed before steps run
}

// Result reports per-step outcome from Apply.
type Result struct {
	Applied []Step
	Failed  []FailedStep
}

type FailedStep struct {
	Step Step
	Err  error
}

func (r Result) AllOK() bool { return len(r.Failed) == 0 }

func (r Result) ExitCode() int {
	switch {
	case r.AllOK() && len(r.Applied) > 0:
		return 0
	case len(r.Applied) == 0 && len(r.Failed) > 0:
		return 2
	default:
		return 3 // partial
	}
}
