package manifest

import "github.com/spf13/cobra"

var manifestFileName string

func BindCommonManifestFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&manifestFileName, "manifest-file", "m", "manifest.br", "manifest file name")
}
