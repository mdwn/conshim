package manifest

import (
	"github.com/spf13/cobra"
)

var (
	removeShimCmd = &cobra.Command{
		Use:   "remove-shim",
		Short: "Removes a shim from the manifest.",
		Long:  "Removes a shim entry from the manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			cobra.CheckErr(m.RemoveShim(shimName))

			writeManifestFile(m)
		},
	}
)

func init() {
	bindCommonManifestFlags(removeShimCmd)
	bindShimFlags(removeShimCmd)
}
