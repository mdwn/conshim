package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	getCmdShimName string

	getShimCmd = &cobra.Command{
		Use:   "get-shim <name>",
		Short: "Gets a shim from the manifest.",
		Long:  "Gets a shim entry from the manifest.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)

			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			getCmdShimName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			if manifestShim, ok := m.GetShim(getCmdShimName); ok {
				fmt.Println(manifestShim)
			} else {
				fmt.Printf("No shim '%s' was found in the manifest.\n", shimName)
			}
		},
	}
)

func init() {
	bindCommonManifestFlags(getShimCmd)
}
