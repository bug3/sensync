package adapter

import "testing"

func TestRecordingExecutorCapturesCalls(t *testing.T) {
	r := &RecordingExecutor{}
	steps := []Step{
		{Kind: StepExec, Target: "hyprctl", Args: []string{"keyword", "input:sensitivity", "0"}},
		{Kind: StepWriteFile, Target: "/tmp/x.conf", Args: []string{"contents"}},
	}
	for _, s := range steps {
		if err := r.Do(s); err != nil {
			t.Fatalf("Do: %v", err)
		}
	}
	if got, want := len(r.Calls), 2; got != want {
		t.Fatalf("Calls: got %d, want %d", got, want)
	}
	if r.Calls[0].Kind != StepExec || r.Calls[0].Target != "hyprctl" {
		t.Errorf("first call wrong: %+v", r.Calls[0])
	}
}

func TestRecordingExecutorReturnsConfiguredErrors(t *testing.T) {
	r := &RecordingExecutor{
		Errors: map[string]error{"hyprctl": errFake},
	}
	err := r.Do(Step{Kind: StepExec, Target: "hyprctl"})
	if err != errFake {
		t.Errorf("got %v, want errFake", err)
	}
}

var errFake = stringErr("fake failure")

type stringErr string

func (e stringErr) Error() string { return string(e) }
