package manifest

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "manifest [command]",
		Short: "Commands for managing conshim manifests.",
		Long:  `Various commands for the creation and management of conshim manifests.`,
	}
)

func init() {
	rootCmd.AddCommand(addShimCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(infoCmd)
}

func Root() *cobra.Command {
	return rootCmd
}