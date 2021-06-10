package config

import (
	"github.com/meowfaceman/conshim/pkg/manifest"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WriteManifestToConfigDirectory will write the manifest as a compressed, serialized file to the config directory.
func WriteManifestToConfigDirectory(m *manifest.Manifest) error {
	file, err := configDir.CreateRegistryFile(m.SourceHash())

	if err != nil {
		return errors.Wrap(err, "error creating registry file")
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			zap.S().Errorf("error closing registry file: %v", closeErr)
		}
	}()

	if writeErr := m.WriteManifest(file); writeErr != nil {
		return errors.Wrap(writeErr, "error writing manifest file")
	}

	return nil
}

// UpdateManifestInConfigDirectory will update the manifest file in the config directory.
func UpdateManifestInConfigDirectory(m *manifest.Manifest) error {
	file, err := configDir.UpdateRegistryFile(m.SourceHash())

	if err != nil {
		return errors.Wrap(err, "error updating registry file")
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			zap.S().Errorf("error closing registry file: %v", closeErr)
		}
	}()

	if writeErr := m.WriteManifest(file); writeErr != nil {
		return errors.Wrap(writeErr, "error writing manifest file")
	}

	return nil
}

// ReadManifestFromConfigDirectory will read the manifest from the config directory.
func ReadManifestFromConfigDirectory(sourceName string) (*manifest.Manifest, error) {
	return ReadManifestFromRawFileInConfigDirectory(manifest.SourceHash(sourceName))
}

// ReadManifestFromRawFileInConfigDirectory will read the manifest from the config directory given a raw file name.
func ReadManifestFromRawFileInConfigDirectory(filename string) (*manifest.Manifest, error) {
	file, closer, err := configDir.GetRegistryFile(filename)

	if err != nil {
		return nil, errors.Wrap(err, "error opening registry file")
	}

	defer closer()

	m, err := manifest.ReadManifest(file)
	if err != nil {
		return nil, errors.Wrap(err, "error reading manifest file")
	}

	return m, nil
}

// ListRegistryManifestFiles will return a list of registry manifest files from the config directory.
func ListRegistryManifestFiles() ([]string, error) {
	registryManifests, err := configDir.ListRegistryFiles()

	if err != nil {
		return nil, errors.Wrap(err, "error listing the registry manifests from config")
	}

	return registryManifests, nil
}
