package registry

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/meowfaceman/conshim/pkg/registry"
	"github.com/spf13/cobra"
)

var (
	addRegistryName string
	addCmd          = &cobra.Command{
		Use:   "add",
		Short: "Adds a registry to conshim.",
		Long:  "Adds a registry to the local conshim configuration.",

		Args: func(cmd *cobra.Command, args []string) error {
			numArgs := len(args)
			if numArgs != 1 {
				return fmt.Errorf("expected 1 argument, got %d", numArgs)
			}

			addRegistryName = args[0]

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			registry, err := registry.GetRegistry(addRegistryName)
			cobra.CheckErr(err)

			cobra.CheckErr(config.WriteManifestToConfigDirectory(registry.GetManifest()))
		},
	}
)
