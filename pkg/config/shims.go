package config

import (
	"os"

	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// AddShim will add a shim that calls a separate command. The intent is that
// this is used for container commands, though it's not strictly necessary.
func AddShim(shimName, source, version string, params []string, command string) error {
	newShim := shim.Shim{
		Source:     source,
		Name:       shimName,
		Version:    version,
		Parameters: params,
		Command:    command,
	}

	renderedShim, err := newShim.RenderShim(map[string]string{})

	if err != nil {
		return errors.Wrap(err, "error rendering shim")
	}

	if err := configDir.AddBinFile(shimName, []byte(renderedShim)); err != nil {
		return errors.Wrap(err, "error adding shim")
	}

	return nil
}

// BinPath returns the bin path where the shims are located.
func BinPath() string {
	return configDir.GetBinPath()
}

// List will return a list of all currently managed shims.
func ListShims() ([]shim.Shim, error) {
	shimFiles, err := configDir.ListBinFiles()

	if err != nil {
		return nil, errors.Wrap(err, "error listing shims from bin directory")
	}

	var shims []shim.Shim
	for _, shimFile := range shimFiles {
		fullPath := configDir.GetBinFileName(shimFile)

		f, readErr := os.Open(fullPath)

		if readErr != nil {
			return nil, errors.Wrap(err, "error reading shim file")
		}

		func() {
			defer func() {
				if closeErr := f.Close(); closeErr != nil {
					zap.S().Errorf("error clsoing shim file '%s': %v", fullPath, closeErr)
				}
			}()
			shims = append(shims, shim.ParseShimFromReader(shimFile, f))
		}()
	}

	return shims, nil
}

// UpdateShim will update an existing shim with a new command.
func UpdateShim(shimName, source, version string, params []string, command string) error {
	newShim := shim.Shim{
		Source:     source,
		Name:       shimName,
		Version:    version,
		Parameters: params,
		Command:    command,
	}

	renderedShim, err := newShim.RenderShim(map[string]string{})

	if err != nil {
		return errors.Wrap(err, "error rendering shim")
	}

	if err := configDir.UpdateBinFile(shimName, []byte(renderedShim)); err != nil {
		return errors.Wrap(err, "error updating shim")
	}

	return nil
}
