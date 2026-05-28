//go:build linux

package adapter

import (
	"testing"

	"github.com/bug3dev/sensync/internal/config"
)

func TestHyprlandAdapterApplyDryRunRecordsAllSteps(t *testing.T) {
	rec := &RecordingExecutor{}
	a := &HyprlandAdapter{exec: rec, confPath: "/tmp/sensync.conf"}
	plan, err := a.BuildPlan(config.Default())
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	res, err := a.Apply(plan, true)
	if err != nil {
		t.Fatalf("Apply dry: %v", err)
	}
	if len(rec.Calls) != 0 {
		t.Errorf("dry-run should not invoke executor; got %d calls", len(rec.Calls))
	}
	if !res.AllOK() {
		t.Errorf("dry-run AllOK should be true")
	}
}

func TestHyprlandAdapterApplyExecutesEveryStep(t *testing.T) {
	rec := &RecordingExecutor{}
	a := &HyprlandAdapter{exec: rec, confPath: "/tmp/sensync.conf"}
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
		t.Errorf("AllOK: got %+v", res)
	}
}

func TestHyprlandAdapterReportsFailedSteps(t *testing.T) {
	rec := &RecordingExecutor{Errors: map[string]error{"hyprctl": errFake}}
	a := &HyprlandAdapter{exec: rec, confPath: "/tmp/sensync.conf"}
	plan, _ := a.BuildPlan(config.Default())
	res, _ := a.Apply(plan, false)
	if len(res.Failed) == 0 {
		t.Error("expected failures, got none")
	}
}
