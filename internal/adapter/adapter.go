package adapter

import (
	"github.com/bug3dev/sensync/internal/adapter/types"
	"github.com/bug3dev/sensync/internal/config"
)

// Re-export shared value types so existing callers within the adapter
// package can use `adapter.Step` style names. Plan generators in the
// `adapter/plan` subpackage use the canonical `types` import directly.
type (
	StepKind   = types.StepKind
	Step       = types.Step
	Plan       = types.Plan
	Result     = types.Result
	FailedStep = types.FailedStep
)

const (
	StepWriteFile = types.StepWriteFile
	StepExec      = types.StepExec
	StepRegSet    = types.StepRegSet
	StepSysCall   = types.StepSysCall
)

// Adapter is the per-OS contract. One adapter is selected at build time via
// build tags; each OS implements all four methods.
type Adapter interface {
	Name() string
	BuildPlan(cfg config.Config) (Plan, error)
	Apply(p Plan, dryRun bool) (Result, error)
	Get() (config.Config, error)
}
