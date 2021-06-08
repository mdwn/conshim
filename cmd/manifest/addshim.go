package manifest

import (
	"github.com/spf13/cobra"
)

var (
	addShimCmd = &cobra.Command{
		Use:   "add-shim",
		Short: "Adds a shim to the manifest.",
		Long:  "Adds a shim entry to the manifest.",

		Args: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement this.
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement this.
		},
	}
)

func init() {
	BindCommonManifestFlags(addShimCmd)
}
