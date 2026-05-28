package adapter

import "github.com/bug3dev/sensync/internal/config"

// StepKind classifies a single change a Plan will perform.
type StepKind string

const (
	StepWriteFile StepKind = "WRITE_FILE"   // write to a file (e.g., Hyprland source file)
	StepExec      StepKind = "EXEC"         // shell-out (e.g., hyprctl, defaults)
	StepRegSet    StepKind = "REG_SET"      // Windows registry value
	StepSysCall   StepKind = "SYS_CALL"     // Windows SystemParametersInfo broadcast
)

// Step is a single declarative change. Adapters generate Steps in the planning
// phase; the Executor performs them in the apply phase.
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

// Adapter is the per-OS contract. One adapter is selected at build time via
// build tags; each OS implements all three methods.
type Adapter interface {
	Name() string
	BuildPlan(cfg config.Config) (Plan, error)
	Apply(p Plan, dryRun bool) (Result, error)
	Get() (config.Config, error)
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

func (r Result) AllOK() bool   { return len(r.Failed) == 0 }
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
