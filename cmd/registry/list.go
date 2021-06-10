package registry

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists registries that have been added to conshim.",
		Long:  "Lists registries that are currently present in conshim's local configuration.",

		Run: func(cmd *cobra.Command, args []string) {
			registryManifests, err := config.ListRegistryManifestFiles()
			cobra.CheckErr(err)

			for _, manifest := range registryManifests {
				m, err := config.ReadManifestFromRawFileInConfigDirectory(manifest)
				cobra.CheckErr(err)

				fmt.Printf("%s (%s)\n", m.Source, m.Version)
			}
		},
	}
)
