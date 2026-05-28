package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/bug3dev/sensync/internal/adapter"
	"github.com/bug3dev/sensync/internal/config"
)

func newApplyCmd() *cobra.Command {
	var (
		dryRun     bool
		yes        bool
		configPath string
	)
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the sensync config to this host",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := resolveConfigPath(configPath)
			if err != nil {
				return userErr(err)
			}
			cfg, err := config.Load(path)
			if err != nil {
				return userErr(err)
			}

			a, err := pickAdapter()
			if err != nil {
				return hostErr(err)
			}

			plan, err := a.BuildPlan(cfg)
			if err != nil {
				return hostErr(err)
			}

			for _, w := range plan.Warnings {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning:", w)
			}
			if !yes && needsTrackpadPrompt(cfg) {
				if !isStdinTTY() {
					return userErr(errors.New("mouse/trackpad settings differ and stdin is not a TTY; pass --yes to skip prompt"))
				}
				if !confirmTTY(cmd, "Apply may overwrite mouse settings with trackpad values. Continue? [y/N]: ") {
					return userErr(errors.New("aborted by user"))
				}
			}

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] dry-run, %d step(s):\n", a.Name(), len(plan.Steps))
				for _, s := range plan.Steps {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", s.Desc)
				}
				return nil
			}

			res, err := a.Apply(plan, false)
			if err != nil {
				return hostErr(err)
			}
			printApplySummary(cmd, a.Name(), res)
			if code := res.ExitCode(); code != 0 {
				return cliError{code: code}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print planned steps without changing anything")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().StringVar(&configPath, "config", "", "Explicit path to sensync.toml")
	return cmd
}

func resolveConfigPath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if _, err := os.Stat("sensync.toml"); err == nil {
		return "sensync.toml", nil
	}
	cfgDir, err := os.UserConfigDir()
	if err == nil {
		candidate := filepath.Join(cfgDir, "sensync", "config.toml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", errors.New("no config found: pass --config or create ./sensync.toml")
}

func needsTrackpadPrompt(cfg config.Config) bool {
	return cfg.Mouse.Sensitivity != cfg.Trackpad.Sensitivity ||
		cfg.Mouse.Acceleration != cfg.Trackpad.Acceleration
}

func isStdinTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func confirmTTY(cmd *cobra.Command, prompt string) bool {
	fmt.Fprint(cmd.OutOrStdout(), prompt)
	r := bufio.NewReader(cmd.InOrStdin())
	line, err := r.ReadString('\n')
	if err != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	}
	return false
}

func printApplySummary(cmd *cobra.Command, name string, res adapter.Result) {
	fmt.Fprintf(cmd.OutOrStdout(), "[%s] applied %d step(s), failed %d\n", name, len(res.Applied), len(res.Failed))
	for _, f := range res.Failed {
		fmt.Fprintf(cmd.ErrOrStderr(), "  ! %s: %v\n", f.Step.Desc, f.Err)
	}
	if name == "macos" {
		fmt.Fprintln(cmd.OutOrStdout(), "macOS: log out / restart for full effect.")
	}
}

// cliError carries a process exit code from RunE up to main(). When err is
// nil, main() exits with the code without printing anything (the caller has
// already produced any user-facing output, e.g. printApplySummary).
type cliError struct {
	err  error
	code int
}

func (e cliError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e cliError) Unwrap() error { return e.err }

func userErr(err error) error { return cliError{err: err, code: 1} }
func hostErr(err error) error { return cliError{err: err, code: 2} }
