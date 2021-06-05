package shims

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/meowfaceman/conshim/assets"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// Only support bash templates for now
	bashTemplate = "bash"
)

var (
	templates = map[string]*template.Template{}

	sourceVersionRegex = regexp.MustCompile(`^#\s*source:\s*([^\s]+)\s*version:\s(.*[^\s])\s*$`)
)

func init() {
	shimsTemplatesPath := "shims"

	dirEntries, err := assets.ShimTemplates.ReadDir(shimsTemplatesPath)

	if err != nil {
		panic(fmt.Sprintf("error reading shim templates: %v", err))
	}

	for _, dirEntry := range dirEntries {
		templateName := dirEntry.Name()

		data, readErr := assets.ShimTemplates.ReadFile(filepath.Join(shimsTemplatesPath, templateName))

		if readErr != nil {
			panic(fmt.Sprintf("error reading shim template '%s': %v", templateName, readErr))
		}

		templates[templateName] = template.Must(template.New(templateName).Parse(string(data)))
	}
}

// ShimInfo is a collection of information about a shim.
type ShimInfo struct {
	// Name is the name of the shim.
	Name string

	// Version is the version of the shim.
	Version string

	// Source is the source of the shim.
	Source string

	// Command is the command associated with the shim.
	Command string
}

// Add will add a shim that calls a separate command. The intent is that
// this is used for container commands, though it's not strictly necessary.
func Add(shim string, command string) error {
	renderedShim := &bytes.Buffer{}

	if err := templates[bashTemplate].Execute(renderedShim, []string{time.Now().String(), command}); err != nil {
		return errors.Wrap(err, "error rendering template for add")
	}

	if err := configDir.AddBinFile(shim, renderedShim.Bytes()); err != nil {
		return errors.Wrap(err, "error adding shim")
	}

	return nil
}

// BinPath returns the bin path where the shims are located.
func BinPath() string {
	return configDir.GetBinPath()
}

// List will return a list of all currently managed shims.
func List() ([]ShimInfo, error) {
	shimFiles, err := configDir.ListBinFiles()

	if err != nil {
		return nil, errors.Wrap(err, "error listing shims from bin directory")
	}

	var shims []ShimInfo
	for _, shimFile := range shimFiles {
		shims = append(shims, getShimInfo(shimFile))
	}

	return shims, nil
}

// Update will update an existing shim with a new command.
func Update(shim string, command string) error {
	renderedShim := &bytes.Buffer{}

	if err := templates[bashTemplate].Execute(renderedShim, []string{time.Now().String(), command}); err != nil {
		return errors.Wrap(err, "error rendering template for update")
	}

	if err := configDir.UpdateBinFile(shim, renderedShim.Bytes()); err != nil {
		return errors.Wrap(err, "error updating shim")
	}

	return nil
}

func getShimInfo(shimFile string) ShimInfo {
	shimInfo := ShimInfo{
		Name: shimFile,
	}

	contents, readErr := configDir.GetBinFile(shimFile)

	if readErr != nil {
		fmt.Printf("%v", readErr)
		zap.S().Debugf("error reading contents of shim '%s': %v", shimFile, readErr)
		shimInfo.Command = "error reading contents"

		return shimInfo
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	// Shims should have three lines: a shebang header, a source/version, and the actual command.
	if !scanner.Scan() {
		shimInfo.Command = "unexpected EOF while skipping shebang line"
	}

	if !scanner.Scan() {
		shimInfo.Command = "unexpected EOF while reading source/version"
	}

	sourceVersion := scanner.Text()
	matches := sourceVersionRegex.FindStringSubmatch(sourceVersion)

	if len(matches) != 3 {
		shimInfo.Source = "???"
		shimInfo.Version = "???"
	} else {
		shimInfo.Source = matches[1]
		shimInfo.Version = matches[2]
	}

	if !scanner.Scan() {
		shimInfo.Command = "unexpected EOF while reading command"
	}

	shimInfo.Command = scanner.Text()

	return shimInfo
}
