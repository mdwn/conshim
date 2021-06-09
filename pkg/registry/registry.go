package registry

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// ManifestFilename is the default manifest filename.
	ManifestFilename = "manifest.br"
)

var (
	isSSHRegex = regexp.MustCompile(`^[A-Za-z_][\w-]*\$?@.*$`)
	hasSchema  = regexp.MustCompile(`^[A-Za-z]*://.*$`)
)

// Registry is a git repository containing a conshim manifest.
type Registry struct {
	// manifest is the manifest associated with the registry. This manifest is expected to be titled `manifest.br` in the root of the directory.
	manifest *Manifest
}

// GetRegistry will attempt to get a manifest file from the given registry.
func GetRegistry(url string) (*Registry, error) {
	tmpDir, err := ioutil.TempDir("", "")

	if err != nil {
		return nil, errors.Wrap(err, "error creating temporary directory while retrieving registry")
	}

	defer func() {
		if removeErr := os.RemoveAll(tmpDir); removeErr != nil {
			zap.S().Errorf("error removing temporary directory during registry retrieval: %v", removeErr)
		}
	}()

	// Clone the repository and try to get the manifest file from the root.
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: mungeURL(url),
	})

	if err != nil {
		return nil, errors.Wrap(err, "error cloning git repository during registry retrieval")
	}

	worktree, err := repo.Worktree()

	if err != nil {
		return nil, errors.Wrap(err, "error getting worktree during registry retrieval")
	}

	manifestFile, err := worktree.Filesystem.Open(ManifestFilename)

	if err != nil {
		return nil, errors.Wrap(err, "error opening manifest from git repository during registry retrieval")
	}

	defer func() {
		if closeErr := manifestFile.Close(); closeErr != nil {
			zap.S().Errorf("error closing manifest file: %v", err)
		}
	}()

	manifest, err := ReadManifest(manifestFile)

	if err != nil {
		return nil, errors.Wrap(err, "error reading manifest from git repository during registry retrieval")
	}

	return &Registry{
		manifest: manifest,
	}, nil
}

// GetManifest will return the manifest from the registry.
func (r *Registry) GetManifest() *Manifest {
	return r.manifest
}

// mungeURL will attempt to detect if a URL is missing a schema. If it is, it will prepend "https" to it.
func mungeURL(url string) string {
	if isSSHRegex.MatchString(url) {
		return url
	} else if !hasSchema.MatchString(url) {
		return "https://" + url
	}

	return url
}
