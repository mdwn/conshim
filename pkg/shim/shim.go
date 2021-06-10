package shim

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/meowfaceman/conshim/assets"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// Only support bash templates for now
	bashTemplate = "bash"

	shebangMissingErrorMessage  = "unexpected EOF while skipping shebang line"
	metadataMissingErrorMessage = "unexpected EOF while reading metadata"
	commandMissingErrorMessage  = "unexpected EOF while reading command"
)

var (
	templates          = map[string]*template.Template{}
	sourceVersionRegex = regexp.MustCompile(`^#\s*source:\s*([^\s]+)\s*version:\s*(.+)\s*$`)
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

// Shim is a descriptor of a shim.
type Shim struct {
	// Name is the name of the shim.
	Name string `json:"name,omitempty"`

	// Source is the source of the shim.
	Source string `json:"source,omitempty"`

	// Version is the version of the shim represented in the manifest.
	Version string `json:"version"`

	// Description is the description string for the shim.
	Description string `json:"description,omitempty"`

	// Parameters are parameters that can be used for the shim command.
	Parameters []string `json:"parameters"`

	// Command is the shim comman.
	Command string `json:"command"`
}

// String will return a string representation of the shim.
func (s Shim) String() string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("     Source: %s\n", s.Source))
	builder.WriteString(fmt.Sprintf("       Name: %s\n", s.Name))
	builder.WriteString(fmt.Sprintf("    Version: %s\n", s.Version))
	builder.WriteString(fmt.Sprintf("Description: %s\n", s.Description))

	if len(s.Parameters) > 0 {
		builder.WriteString(fmt.Sprintf(" Parameters: %s\n", strings.Join(s.Parameters, ",")))
	}

	builder.WriteString(fmt.Sprintf("    Command: %s", s.Command))

	return builder.String()
}

// RenderShim will render the shim and replace the parameters with the provided values.
func (s Shim) RenderShim(parameters map[string]string) (string, error) {
	renderedShim := &bytes.Buffer{}

	if err := templates[bashTemplate].Execute(renderedShim, s); err != nil {
		return "", errors.Wrap(err, "error rendering template for add")
	}

	replacerArgs := []string{}

	for parameter, value := range parameters {
		replacerArgs = append(replacerArgs, fmt.Sprintf("{{%s}}", parameter), value)
	}

	return strings.NewReplacer(replacerArgs...).Replace(renderedShim.String()), nil
}

// ParseShimFromReader will parse a shim object from a file. Right now this only supports bash templates.
func ParseShimFromReader(shimFile string, reader io.Reader) Shim {
	shimInfo := Shim{
		Name:    shimFile,
		Source:  "???",
		Version: "???",
	}

	contents, readErr := io.ReadAll(reader)

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
	if numMatches == 3 {
		shimInfo.Source = matches[1]
		shimInfo.Version = matches[2]
	}

	if !scanner.Scan() {
		shimInfo.Command = commandMissingErrorMessage
	} else {
		shimInfo.Command = scanner.Text()
	}

	return shimInfo
}

// ShimsListToString will take a list of shims and create a string.
func ShimsListToString(shims []Shim) string {
	entries := []string{}
	for _, shim := range shims {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("       Name: %s\n", shim.Name))
		if shim.Source != "" {
			builder.WriteString(fmt.Sprintf("     Source: %s\n", shim.Source))
		}
		builder.WriteString(fmt.Sprintf("    Version: %s\n", shim.Version))

		if shim.Description != "" {
			builder.WriteString(fmt.Sprintf("Description: %s\n", shim.Description))
		}

		if len(shim.Parameters) > 0 {
			builder.WriteString(fmt.Sprintf(" Parameters: %s\n", strings.Join(shim.Parameters, ",")))
		}

		builder.WriteString(fmt.Sprintf("    Command: %s\n", shim.Command))
		entries = append(entries, builder.String())
	}

	return strings.Join(entries, "-------\n")
}
