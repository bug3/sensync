package main

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Print live system mouse settings in sensync config format",
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := pickAdapter()
			if err != nil {
				return hostErr(err)
			}
			cfg, err := a.Get()
			if err != nil {
				return hostErr(err)
			}
			return toml.NewEncoder(cmd.OutOrStdout()).Encode(cfg)
		},
	}
}
