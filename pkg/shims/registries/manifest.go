package registries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	brotliReaderBufferSize = 1024
)

// Manifest is a manifest of shims that are housed externally.
type Manifest struct {
	// Source is the name of the source this manifest represents.
	Source string `json:"source"`

	// Shims is a list of shims described by this manifest. The key here is the name of the shim
	// which corresponds to the executable name for this shim.
	Shims map[string]Shim `json:"shims"`
}

// Shim is a descriptor of a shim in a manifest.
type Shim struct {
	// Version is the version of teh shim represented in the manifest.
	Version string `json:"version"`

	// Command is the shim comman.
	Command string `json:"command"`

	// Parameters are parameters that can be used for the shim command.
	Parameters []string `json:"parameters"`
}

// CreateManifest will create a new manifest with the given source.
func CreateManifest(source string) *Manifest {
	return &Manifest{
		Source: source,
		Shims:  map[string]Shim{},
	}
}

// ReadManifest will read the manifest from the reader and return it.
func ReadManifest(src io.Reader) (*Manifest, error) {
	bReader := brotli.NewReader(src)

	data := bytes.Buffer{}
	for {
		buf := make([]byte, brotliReaderBufferSize)
		numRead, err := bReader.Read(buf)

		if err != nil {
			// If the error is EOF, break out of the loop gracefully.
			if err == io.EOF {
				break
			}

			return nil, errors.Wrap(err, "error decompressing manifest")
		}

		if numRead == 0 {
			break
		}

		data.Write(buf[:numRead])
	}

	manifest := &Manifest{}
	if err := json.Unmarshal(data.Bytes(), manifest); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling decompressed manifest")
	}

	return manifest, nil
}

// WriteManifest will take the current manifest and write it to the writer.
func (m *Manifest) WriteManifest(dst io.Writer) error {
	data, err := json.Marshal(m)

	if err != nil {
		return errors.Wrap(err, "error marshaling manifest")
	}

	bWriter := brotli.NewWriter(dst)

	defer func() {
		if closeErr := bWriter.Close(); closeErr != nil {
			zap.S().Errorf("error closing brotli writer: %v", closeErr)
		}
	}()
	numWritten, err := bWriter.Write(data)

	if err != nil {
		return errors.Wrap(err, "error compressing marshaled manifest")
	}

	dataLength := len(data)
	if dataLength != numWritten {
		return fmt.Errorf("%d bytes written, expected %d during manifest compression", numWritten, dataLength)
	}

	return nil
}

// AddShim will add one or more shims to the manifest. If the shim already exists, this will error.
func (m *Manifest) AddShim(shimName string, shim Shim) error {
	if _, ok := m.Shims[shimName]; ok {
		return fmt.Errorf("shim '%s' already exists in the manifest", shimName)
	}

	m.Shims[shimName] = shim

	return nil
}

// RemoveShim will remove the shim with the given name from the manifest.
func (m *Manifest) RemoveShim(shimName string) error {
	if _, ok := m.Shims[shimName]; !ok {
		return fmt.Errorf("shim '%s' does not exist in the manifest", shimName)
	}

	delete(m.Shims, shimName)

	return nil
}
