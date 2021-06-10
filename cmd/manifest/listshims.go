package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listShimsCmd = &cobra.Command{
		Use:   "list-shims",
		Short: "List shims from the manifest.",
		Long:  "List shim entries from the manifest.",

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			fmt.Print(m.ShimsToString())
		},
	}
)

func init() {
	bindCommonManifestFlags(listShimsCmd)
}
