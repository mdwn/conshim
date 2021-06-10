package manifest

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/meowfaceman/conshim/pkg/shim"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	brotliReaderBufferSize = 1024
)

// Manifest is a manifest of shims that are housed externally.
type Manifest struct {
	// Source is the URL of the registry that this manifest belongs to.
	Source string `json:"source"`

	// Version is the version of the manifest.
	Version string `json:"version"`

	// Shims is a list of shims described by this manifest. The key here is the name of the shim
	// which corresponds to the executable name for this shim.
	Shims map[string]shim.Shim `json:"shims"`
}

// CreateManifest will create a new manifest with the given source.
func CreateManifest(source string) *Manifest {
	return &Manifest{
		Source: source,
		Shims:  map[string]shim.Shim{},
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
func (m *Manifest) AddShim(shimName string, shim shim.Shim) error {
	if _, ok := m.Shims[shimName]; ok {
		return fmt.Errorf("shim '%s' already exists in the manifest", shimName)
	}

	m.Shims[shimName] = shim

	return nil
}

// GetShim will get the shim from the manifest. The boolean will be true if the shim was found in the manifest.
func (m *Manifest) GetShim(shimName string) (shim.Shim, bool) {
	if _, ok := m.Shims[shimName]; !ok {
		return shim.Shim{}, false
	}

	shim := m.Shims[shimName]

	// Inject the source and shim name into the shim object.
	shim.Source = m.Source
	shim.Name = shimName

	return shim, true
}

// UpdateShim will update the shim in the manifest.
func (m *Manifest) UpdateShim(shimName string, shim shim.Shim) error {
	if _, ok := m.Shims[shimName]; !ok {
		return fmt.Errorf("shim '%s' does not exist in the manifest", shimName)
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

// SourceHash returns a hash of the source name attached to this manifest.
func (m *Manifest) SourceHash() string {
	return SourceHash(m.Source)
}

func (m *Manifest) ShimsToString() string {
	shims := []shim.Shim{}

	for shimName, shim := range m.Shims {
		shim.Name = shimName
		shims = append(shims, shim)
	}

	return shim.ShimsListToString(shims)
}

// SourceHash generates a hash of the given source name.
func SourceHash(sourceName string) string {
	hash := sha256.New()
	hash.Write([]byte(sourceName))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
