package adapter

import (
	"fmt"
	"os"
	"os/exec"
)

// Executor performs a single Step. Tests substitute a recording fake; the
// production adapter uses ShellExecutor (Linux/macOS) or the OS-native
// executor declared in the Windows-tagged file.
type Executor interface {
	Do(s Step) error
}

// ShellExecutor runs StepExec via os/exec and StepWriteFile via os.WriteFile.
// Other kinds return an error indicating they belong to a platform-specific
// executor.
type ShellExecutor struct{}

func (ShellExecutor) Do(s Step) error {
	switch s.Kind {
	case StepExec:
		if s.Target == "" {
			return fmt.Errorf("exec step missing Target")
		}
		cmd := exec.Command(s.Target, s.Args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("exec %s %v: %w (output: %s)", s.Target, s.Args, err, string(out))
		}
		return nil
	case StepWriteFile:
		if s.Target == "" || len(s.Args) != 1 {
			return fmt.Errorf("write_file step needs Target and Args[0]=contents")
		}
		return os.WriteFile(s.Target, []byte(s.Args[0]), 0o644)
	default:
		return fmt.Errorf("ShellExecutor cannot perform step kind %q (wrong platform?)", s.Kind)
	}
}

// RecordingExecutor is a test-only fake. It records every Step it sees and
// can be configured to return errors for specific Targets.
type RecordingExecutor struct {
	Calls  []Step
	Errors map[string]error
}

func (r *RecordingExecutor) Do(s Step) error {
	r.Calls = append(r.Calls, s)
	if r.Errors != nil {
		if err, ok := r.Errors[s.Target]; ok {
			return err
		}
	}
	return nil
}
