package registry

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command = &cobra.Command{
		Use:   "registry [command]",
		Short: "Commands for managing conshim registries.",
		Long:  `Various commands for the management of conshim registries`,
	}
)

func init() {
}

func Root() *cobra.Command {
	return rootCmd
}
