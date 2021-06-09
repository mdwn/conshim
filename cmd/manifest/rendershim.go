package manifest

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	renderCmdShimName   string
	renderCmdParameters map[string]string

	renderShimCmd = &cobra.Command{
		Use:   "render-shim <name>",
		Short: "Renders a shim from the manifest.",
		Long:  "Renders a shim entry from the manifest.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)

			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			renderCmdShimName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			m, closeFunc := readManifestFile()
			defer closeFunc()

			if manifestShim, ok := m.GetShim(renderCmdShimName); ok {
				renderedShim, err := manifestShim.RenderShim(renderCmdParameters)
				cobra.CheckErr(err)
				fmt.Println(renderedShim)
			} else {
				fmt.Printf("No shim '%s' was found in the manifest.\n", shimName)
			}
		},
	}
)

func init() {
	bindCommonManifestFlags(renderShimCmd)

	renderShimCmd.Flags().StringToStringVarP(&renderCmdParameters, "parameters", "p", map[string]string{}, "parameters and values for the command")
}
