package registry

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/spf13/cobra"
)

var (
	loadShimCmdRegistryName string
	loadShimCmdShimName     string
	loadShimCmdParameters   map[string]string
	loadShimCmdUpdate       bool

	loadShimCmd = &cobra.Command{
		Use:   "load-shim <registry> <shim>",
		Short: "Loads a shim from the given registry.",
		Long:  "Loads a shim from the given registry.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)

			if len(args) != 2 {
				return fmt.Errorf("expected 2 arguments, got %d", numArgs)
			}

			loadShimCmdRegistryName = args[0]
			loadShimCmdShimName = args[1]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			m, err := addOrGetRegistry(loadShimCmdRegistryName)
			cobra.CheckErr(err)

			if manifestShim, ok := m.GetShim(loadShimCmdShimName); ok {
				renderedShim, err := manifestShim.RenderShim(loadShimCmdParameters)
				cobra.CheckErr(err)

				if loadShimCmdUpdate {
					cobra.CheckErr(config.Directory().UpdateBinFile(loadShimCmdShimName, []byte(renderedShim)))
				} else {
					cobra.CheckErr(config.Directory().AddBinFile(loadShimCmdShimName, []byte(renderedShim)))
				}
			} else {
				fmt.Printf("No shim '%s' was found in registry %s.\n", loadShimCmdShimName, loadShimCmdRegistryName)
			}
		},
	}
)

func init() {
	loadShimCmd.Flags().StringToStringVarP(&loadShimCmdParameters, "parameters", "p", map[string]string{}, "parameters and values for the command")
	loadShimCmd.Flags().BoolVarP(&loadShimCmdUpdate, "update", "u", false, "update and overwrite an existing local shim")
}
