package manifest

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/spf13/cobra"
)

var (
	addShimCmd = &cobra.Command{
		Use:   "add-shim",
		Short: "Adds a shim to the manifest.",
		Long:  "Adds a shim entry to the manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			func() {
				defer closeFunc()

				newShim := shim.Shim{
					Version:     shimVersion,
					Description: shimDescription,
					Parameters:  shimParameters,
					Command:     shimCommand,
				}

				cobra.CheckErr(m.AddShim(shimName, newShim))

				fmt.Printf("Added shim '%s' to manifest %s.\n", shimName, m.Source)
			}()

			writeManifestFile(m)
		},
	}
)

func init() {
	bindCommonManifestFlags(addShimCmd)
	bindShimFlags(addShimCmd)
	bindShimModificationFlags(addShimCmd)
}
