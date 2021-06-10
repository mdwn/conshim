package manifest

import (
	"fmt"
	"os"

	"github.com/meowfaceman/conshim/pkg/manifest"
	"github.com/spf13/cobra"
)

var (
	createCmdSourceName string

	createCmd = &cobra.Command{
		Use:   "create <source-name>",
		Short: "Creates a manifest.",
		Long:  "Creates a manifest with the given source name.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)
			if numArgs != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			createCmdSourceName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			if _, err := os.Stat(manifestFileName); err == nil {
				cobra.CheckErr(fmt.Sprintf("manifest file '%s' already exists", manifestFileName))
			}

			m := manifest.CreateManifest(createCmdSourceName)
			writeManifestFile(m)
		},
	}
)

func init() {
	bindCommonManifestFlags(createCmd)
}
