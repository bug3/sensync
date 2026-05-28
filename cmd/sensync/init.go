package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed sensync.example.toml
var exampleConfigTOML []byte

func newInitCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write the example config to the user config directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := userConfigPath()
			if err != nil {
				return userErr(err)
			}
			if _, err := os.Stat(path); err == nil && !force {
				return userErr(fmt.Errorf("config already exists at %s; pass --force to overwrite", path))
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return hostErr(fmt.Errorf("create config dir: %w", err))
			}
			if err := os.WriteFile(path, exampleConfigTOML, 0o644); err != nil {
				return hostErr(fmt.Errorf("write config: %w", err))
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Wrote example config to %s\n", path)
			fmt.Fprintln(cmd.OutOrStdout(), "Edit the file, then run `sensync apply`.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite an existing config file")
	return cmd
}

func userConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", errors.New("could not determine user config directory")
	}
	return filepath.Join(dir, "sensync", "config.toml"), nil
}
