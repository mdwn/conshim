package registry

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/meowfaceman/conshim/pkg/registry"
	"github.com/spf13/cobra"
)

var (
	updateRegistryName string
	updateCmd          = &cobra.Command{
		Use:   "update <registry-name>",
		Short: "Updates a registry in conshim.",
		Long:  "Updates a registry already present in the local conshim configuration.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)
			if numArgs != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			updateRegistryName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			registry, err := registry.GetRegistry(updateRegistryName)
			cobra.CheckErr(err)

			cobra.CheckErr(config.UpdateManifestInConfigDirectory(registry.GetManifest()))
		},
	}
)
