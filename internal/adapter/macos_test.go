//go:build darwin

package adapter

import (
	"testing"

	"github.com/bug3dev/sensync/internal/config"
)

func TestMacOSAdapterApplyExecutesAllSteps(t *testing.T) {
	rec := &RecordingExecutor{}
	a := &MacOSAdapter{exec: rec}
	plan, err := a.BuildPlan(config.Default())
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	res, err := a.Apply(plan, false)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := len(rec.Calls), len(plan.Steps); got != want {
		t.Errorf("executor calls: got %d want %d", got, want)
	}
	if !res.AllOK() {
		t.Errorf("AllOK: %+v", res)
	}
}

func TestMacOSAdapterDryRunDoesNothing(t *testing.T) {
	rec := &RecordingExecutor{}
	a := &MacOSAdapter{exec: rec}
	plan, _ := a.BuildPlan(config.Default())
	if _, err := a.Apply(plan, true); err != nil {
		t.Fatalf("dry: %v", err)
	}
	if len(rec.Calls) != 0 {
		t.Errorf("dry-run should not invoke executor; got %d calls", len(rec.Calls))
	}
}
