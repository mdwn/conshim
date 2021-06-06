package manifest

import (
	"fmt"
	"os"

	"github.com/meowfaceman/conshim/pkg/shims/registries"
	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Gets info about a manifest.",
		Long:  "Prints out general information about a manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			manifestFile, err := os.Open(manifestFileName)
			cobra.CheckErr(err)

			defer func() {
				cobra.CheckErr(manifestFile.Close())
			}()

			m, err := registries.ReadManifest(manifestFile)
			cobra.CheckErr(err)

			fmt.Printf("Source: %s\n", m.Source)
			fmt.Printf("Number of shims: %d\n", len(m.Shims))
		},
	}
)

func init() {
	BindCommonManifestFlags(infoCmd)
}
