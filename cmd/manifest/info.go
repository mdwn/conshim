package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Gets info about a manifest.",
		Long:  "Prints out general information about a manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			fmt.Printf("Source: %s\n", m.Source)
			fmt.Printf("Number of shims: %d\n", len(m.Shims))
		},
	}
)

func init() {
	bindCommonManifestFlags(infoCmd)
}
