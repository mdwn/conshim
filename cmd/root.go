package cmd

import (
	"github.com/meowfaceman/conshim/cmd/manifest"
	"github.com/meowfaceman/conshim/cmd/registry"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "conshim [command]",
		Short: "Tool for managing container shims.",
		Long: `conshim is a tool that manages small shims that call containers instead of relying on
the local environment.`,
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(binPathCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)

	rootCmd.AddCommand(manifest.Root())
	rootCmd.AddCommand(registry.Root())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Errorf("error running command: %v", err)
	}
}
