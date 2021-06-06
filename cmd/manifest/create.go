package manifest

import (
	"fmt"
	"os"

	"github.com/meowfaceman/conshim/pkg/shims/registries"
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

			manifestFile, err := os.OpenFile(manifestFileName, os.O_CREATE|os.O_WRONLY, 0644)
			cobra.CheckErr(err)

			defer func() {
				cobra.CheckErr(manifestFile.Close())
			}()

			m := registries.CreateManifest(createCmdSourceName)
			cobra.CheckErr(m.WriteManifest(manifestFile))

			fmt.Printf("Manifest wrote to '%s'\n", manifestFileName)
		},
	}
)

func init() {
	BindCommonManifestFlags(createCmd)
}
