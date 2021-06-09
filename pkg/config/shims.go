package config

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	shebangMissingErrorMessage  = "unexpected EOF while skipping shebang line"
	metadataMissingErrorMessage = "unexpected EOF while reading metadata"
	commandMissingErrorMessage  = "unexpected EOF while reading command"
)

var (
	sourceVersionRegex = regexp.MustCompile(`^#\s*source:\s*([^\s]+)\s*version:\s*([^\s]+)\s*(parameters:\s*([^\s]+)\s*)?$`)
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
		shims = append(shims, getShimFromFile(shimFile, fullPath))
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

func getShimFromFile(shimFile, fileName string) shim.Shim {
	shimInfo := shim.Shim{
		Name:    shimFile,
		Source:  "???",
		Version: "???",
	}

	contents, readErr := ioutil.ReadFile(fileName)

	if readErr != nil {
		zap.S().Debugf("error reading contents of shim '%s': %v", shimFile, readErr)
		shimInfo.Command = "error reading contents"

		return shimInfo
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	// Shims should have three lines: a shebang header, a metadata comment line, and the actual command.
	if !scanner.Scan() {
		shimInfo.Command = shebangMissingErrorMessage

		return shimInfo
	}

	if !scanner.Scan() {
		shimInfo.Command = metadataMissingErrorMessage

		return shimInfo
	}

	sourceVersion := scanner.Text()
	matches := sourceVersionRegex.FindStringSubmatch(sourceVersion)

	numMatches := len(matches)
	if numMatches > 4 {
		shimInfo.Source = matches[1]
		shimInfo.Version = matches[2]

		if matches[4] != "" {
			shimInfo.Parameters = strings.Split(matches[4], ",")
		}
	}

	if !scanner.Scan() {
		shimInfo.Command = commandMissingErrorMessage
	} else {
		shimInfo.Command = scanner.Text()
	}

	return shimInfo
}
