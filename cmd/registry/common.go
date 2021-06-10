package registry

import (
	"fmt"

	"github.com/meowfaceman/conshim/pkg/config"
	"github.com/meowfaceman/conshim/pkg/manifest"
	"github.com/meowfaceman/conshim/pkg/registry"
)

func addOrGetRegistry(registryName string) (*manifest.Manifest, error) {
	m, err := config.ReadManifestFromConfigDirectory(registryName)

	if err != nil {
		fmt.Printf("Error finding manifest for registry '%s', attempting to get it.", registryName)
		registry, err := registry.GetRegistry(registryName)

		if err != nil {
			return nil, err
		}

		if writeErr := config.WriteManifestToConfigDirectory(registry.GetManifest()); writeErr != nil {
			return nil, writeErr
		}

		m = registry.GetManifest()
	}

	return m, nil
}
