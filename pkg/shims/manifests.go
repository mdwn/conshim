package shims

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/meowfaceman/conshim/pkg/shims/registries"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WriteManifestToConfigDirectory will write the manifest as a compressed, serialized file to the config directory.
func WriteManifestToConfigDirectory(m *registries.Manifest) error {
	hash := sha256.New()
	hash.Write([]byte(m.Source))
	hashedFilename := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	file, err := configDir.CreateRegistryFile(hashedFilename)

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
