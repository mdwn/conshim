package manifest

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/spf13/cobra"
)

var (
	loadShimCmdShimName   string
	loadShimCmdParameters map[string]string
	loadShimCmdUpdate     bool

	loadShimCmd = &cobra.Command{
		Use:   "load-shim <name>",
		Short: "Loads a shim from the manifest.",
		Long:  "Loads a shim from the manifest into the local conshim config.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)

			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			loadShimCmdShimName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			if manifestShim, ok := m.GetShim(loadShimCmdShimName); ok {
				renderedShim, err := manifestShim.RenderShim(loadShimCmdParameters)
				cobra.CheckErr(err)

				if loadShimCmdUpdate {
					cobra.CheckErr(config.Directory().UpdateBinFile(loadShimCmdShimName, []byte(renderedShim)))
				} else {
					cobra.CheckErr(config.Directory().AddBinFile(loadShimCmdShimName, []byte(renderedShim)))
				}
			} else {
				fmt.Printf("No shim '%s' was found in the manifest.\n", shimName)
			}
		},
	}
)

func init() {
	bindCommonManifestFlags(loadShimCmd)

	loadShimCmd.Flags().StringToStringVarP(&loadShimCmdParameters, "parameters", "p", map[string]string{}, "parameters and values for the command")
	loadShimCmd.Flags().BoolVarP(&loadShimCmdUpdate, "update", "u", false, "update and overwrite an existing local shim")
}
