package main

import "github.com/spf13/cobra"

const appVersion = "0.0.1-dev"

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "sensync",
		Short:         "Cross-platform mouse sensitivity sync",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newVersionCmd())
	root.AddCommand(newApplyCmd())
	root.AddCommand(newGetCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print sensync version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("sensync %s\n", appVersion)
		},
	}
}
