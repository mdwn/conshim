package manifest

import (
	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/spf13/cobra"
)

var (
	updateShimCmd = &cobra.Command{
		Use:   "update-shim",
		Short: "Updates a shim in the manifest.",
		Long:  "Updates a shim entry in the manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			newShim := shim.Shim{
				Version:     shimVersion,
				Description: shimDescription,
				Parameters:  shimParameters,
				Command:     shimCommand,
			}

			cobra.CheckErr(m.UpdateShim(shimName, newShim))

			writeManifestFile(m)
		},
	}
)

func init() {
	bindCommonManifestFlags(updateShimCmd)
	bindShimFlags(updateShimCmd)
	bindShimModificationFlags(updateShimCmd)
}
