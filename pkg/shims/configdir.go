package shims

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/meowfaceman/conshim/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// ConshimConfigDirectory is the location of conshim configs.
	ConshimConfigDirectory = "conshim.config.directory"
)

// ConfigDirectory is the configuration directory where all the shims and configuration live.
type ConfigDirectory struct {
	lock         *flock.Flock
	binPath      string
	registryPath string
}

var (
	configDir *ConfigDirectory
)

func init() {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(fmt.Sprintf("error getting home directory: %v", err))
	}

	utils.Must(viper.BindEnv(ConshimConfigDirectory, "CONSHIM_CONFIG_DIRECTORY"))
	viper.SetDefault(ConshimConfigDirectory, filepath.Join(home, ".conshim"))

	configDir, err = newConfigDirectory()

	if err != nil {
		panic(fmt.Sprintf("error creating config directory: %v", err))
	}
}

// newConfigDirectory will create a config directory object, creating the actual config directory
// if necessary.
func newConfigDirectory() (*ConfigDirectory, error) {
	configDirPath := viper.GetString(ConshimConfigDirectory)

	lockFile := filepath.Join(configDirPath, "lock")
	binPath := filepath.Join(configDirPath, "bin")
	registryPath := filepath.Join(configDirPath, "registries")

	// Make the bin directory if it doesn't already exist.
	if err := os.MkdirAll(binPath, 0700); err != nil {
		return nil, errors.Wrap(err, "error creating configuration directory")
	}

	// Make the registry directory if it doesn't already exist.
	if err := os.MkdirAll(registryPath, 0700); err != nil {
		return nil, errors.Wrap(err, "error creating configuration directory")
	}

	return &ConfigDirectory{
		lock:         flock.New(lockFile),
		binPath:      binPath,
		registryPath: registryPath,
	}, nil
}

// GetBinPath will return the bin path for this config directory.
func (c *ConfigDirectory) GetBinPath() string {
	return c.binPath
}

// AddBinFile will add an executable file to the bin directory.
func (c *ConfigDirectory) AddBinFile(filename string, data []byte) error {
	if err := c.getLock(); err != nil {
		return errors.Wrap(err, "error getting lock while adding bin file")
	}
	defer c.unlock()

	execFilePath := filepath.Join(c.binPath, filename)

	// Add if the file doesn't exist.
	if _, err := os.Stat(execFilePath); err == nil {
		return fmt.Errorf("can't add file '%s' because it already exists", filename)
	}

	if err := ioutil.WriteFile(execFilePath, data, 0700); err != nil {
		return errors.Wrap(err, "error adding bin file")
	}

	return nil
}

// AddBinFile will update an executable file to the bin directory.
func (c *ConfigDirectory) UpdateBinFile(filename string, data []byte) error {
	if err := c.getLock(); err != nil {
		return errors.Wrap(err, "error getting lock while adding bin file")
	}
	defer c.unlock()

	execFilePath := filepath.Join(c.binPath, filename)

	// Update if the file does exist.
	if _, err := os.Stat(execFilePath); err != nil {
		return fmt.Errorf("can't update file '%s' because it doesn't exist", filename)
	}

	if err := ioutil.WriteFile(execFilePath, data, 0700); err != nil {
		return errors.Wrap(err, "error updating bin file")
	}

	return nil
}

// GetBinFile will return the contents of a bin file.
func (c *ConfigDirectory) GetBinFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(c.binPath, filename))

	if err != nil {
		return nil, errors.Wrap(err, "error reading bin file")
	}

	return data, nil
}

// ListBin will return the list of the files in the bin directory.
func (c *ConfigDirectory) ListBinFiles() ([]string, error) {
	if err := c.getLock(); err != nil {
		return nil, errors.Wrap(err, "error getting lock while listing bin files")
	}
	defer c.unlock()

	var shims []string

	err := filepath.Walk(c.binPath, func(path string, info fs.FileInfo, err error) error {
		// If we find a directory, skip it. We only care about shim files here.
		if info.IsDir() {
			return nil
		}

		shims = append(shims, info.Name())

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "error while listing bin file directory")
	}

	return shims, nil
}

// CreateRegistryFile will create a registry file with the given name. It's up to the caller to close this file.
func (c *ConfigDirectory) CreateRegistryFile(name string) (*os.File, error) {
	file, err := os.Create(filepath.Join(c.registryPath, name))

	if err != nil {
		return nil, errors.Wrap(err, "error while creating registry file")
	}

	return file, nil
}

func (c *ConfigDirectory) getLock() error {
	locked, err := c.lock.TryLock()

	if err != nil {
		return errors.Wrap(err, "error while getting config directory lock")
	}

	if !locked {
		return errors.New("config directory lock was not successfully acquired")
	}

	return nil
}

func (c *ConfigDirectory) unlock() {
	if err := c.lock.Unlock(); err != nil {
		zap.S().Warnf("error unlocking config directory lock: %v", err)
	}
}
