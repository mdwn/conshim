package shim

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "shim [command]",
		Short: "Commands for managing shim files.",
		Long:  "Commands for managing shim files in the local config directory.",
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
}

func Root() *cobra.Command {
	return rootCmd
}
