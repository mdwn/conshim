package manifest

import (
	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/spf13/cobra"
)

var (
	updateShimCmd = &cobra.Command{
		Use:   "remove-shim",
		Short: "Removes a shim from the manifest.",
		Long:  "Removes a shim entry from the manifest.",

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
