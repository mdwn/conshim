package manifest

import (
	"fmt"
	"os"

	"github.com/meowfaceman/conshim/pkg/manifest"
	"github.com/spf13/cobra"
)

var (
	manifestFileName string

	shimName        string
	shimVersion     string
	shimDescription string
	shimParameters  []string
	shimCommand     string
)

// bindCommonManifestFlags will bind flags that are common to all manifest commands.
func bindCommonManifestFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&manifestFileName, "manifest-file", "m", "manifest.br", "manifest file name")
}

// bindShimFlags will bind flags that are common to shim commands.
func bindShimFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&shimName, "shim-name", "n", "", "the name of the shim")
}

// bindShimModificationFlags will bind flags that are common to shim modification commands.
func bindShimModificationFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&shimVersion, "shim-version", "v", "", "the version of the shim")
	cmd.Flags().StringVarP(&shimDescription, "shim-description", "d", "", "the description of the shim")
	cmd.Flags().StringSliceVarP(&shimParameters, "shim-parameters", "p", []string{}, "the parameters that can be adjusted for the shim")
	cmd.Flags().StringVarP(&shimCommand, "shim-command", "c", "", "the command executed by the shim")
}

// readManifestFile will read the configured manifest file and return it along with a close function.
func readManifestFile() (*manifest.Manifest, func()) {
	manifestFile, err := os.Open(manifestFileName)
	cobra.CheckErr(err)

	m, err := manifest.ReadManifest(manifestFile)
	cobra.CheckErr(err)

	return m, func() {
		cobra.CheckErr(manifestFile.Close())
	}

}

// writeManifestFile will write the configured manifest file. If it exists already, it will be re-written.
func writeManifestFile(m *manifest.Manifest) {
	manifestFile, err := os.OpenFile(manifestFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	cobra.CheckErr(err)

	defer func() {
		cobra.CheckErr(manifestFile.Close())
	}()

	cobra.CheckErr(m.WriteManifest(manifestFile))

	fmt.Printf("Manifest written to '%s'\n", manifestFileName)
}
