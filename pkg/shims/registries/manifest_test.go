package registries

import (
	"bytes"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

const (
	testSourceName = "dummy"
)

type shimNameAndInfo struct {
	name string
	shim Shim
}

func TestAddShim(t *testing.T) {
	tests := []struct {
		name        string
		shimsToAdd  []shimNameAndInfo
		expectedErr bool
	}{
		{
			name: "add shim",
			shimsToAdd: []shimNameAndInfo{
				{
					name: "new-shim",
					shim: Shim{
						Version:    "1234",
						Command:    "my-command",
						Parameters: []string{"param1", "param2"},
					},
				},
			},
		},
		{
			name: "add two shims",
			shimsToAdd: []shimNameAndInfo{
				{
					name: "new-shim1",
					shim: Shim{
						Version:    "1234",
						Command:    "my-command1",
						Parameters: []string{"param1", "param2"},
					},
				},
				{
					name: "new-shim2",
					shim: Shim{
						Version:    "2345",
						Command:    "my-command2",
						Parameters: []string{"param1", "param2"},
					},
				},
			},
		},
		{
			name: "add duplicate shim name",
			shimsToAdd: []shimNameAndInfo{
				{
					name: "new-shim1",
					shim: Shim{
						Version:    "1234",
						Command:    "my-command1",
						Parameters: []string{"param1", "param2"},
					},
				},
				{
					name: "new-shim1",
					shim: Shim{
						Version:    "2345",
						Command:    "my-command2",
						Parameters: []string{"param1", "param2"},
					},
				},
			},
			expectedErr: true,
		},
	}

	for _, test := range tests {
		m := CreateManifest("dummy")

		// Accumulate errors through adding shims
		err := &multierror.Error{}
		for _, shimToAdd := range test.shimsToAdd {
			err = multierror.Append(err, m.AddShim(shimToAdd.name, shimToAdd.shim))
		}

		didErr := err.ErrorOrNil() != nil
		assert.Equal(t, test.expectedErr, didErr, "error states should equal")

		if !test.expectedErr {
			assert.Len(t, m.Shims, len(test.shimsToAdd), "should have the same number of shims")

			for _, shimToAdd := range test.shimsToAdd {
				if diff := deep.Equal(m.Shims[shimToAdd.name], shimToAdd.shim); diff != nil {
					t.Errorf("%s %s", test.name, diff)
				}
			}
		}
	}
}

func TestRemoveShim(t *testing.T) {
	tests := []struct {
		name           string
		shimsToRemove  []string
		expectedLength int
		expectedErr    bool
	}{
		{
			name:           "remove one shim",
			shimsToRemove:  []string{"new-shim1"},
			expectedLength: 2,
		},
		{
			name:           "remove two shim",
			shimsToRemove:  []string{"new-shim1", "new-shim2"},
			expectedLength: 1,
		},
		{
			name:          "remove non-existent shim",
			shimsToRemove: []string{"non-existent-shim"},
			expectedErr:   true,
		},
	}

	for _, test := range tests {
		m := createTestingManifest(t)

		// Accumulate errors through adding shims
		err := &multierror.Error{}
		for _, shimToRemove := range test.shimsToRemove {
			err = multierror.Append(err, m.RemoveShim(shimToRemove))
		}

		didErr := err.ErrorOrNil() != nil
		assert.Equal(t, test.expectedErr, didErr, "error states should equal")

		if !test.expectedErr {
			assert.Len(t, m.Shims, test.expectedLength, "should have the same number of shims")
		}
	}
}

func TestManifestReadAndWrite(t *testing.T) {
	writer := &bytes.Buffer{}

	m := createTestingManifest(t)
	assert.NoError(t, m.WriteManifest(writer), "should have no error writing the manifest")

	readM, err := ReadManifest(bytes.NewReader(writer.Bytes()))
	assert.NoError(t, err, "reading manifest should have no error")

	if diff := deep.Equal(m, readM); diff != nil {
		t.Error(diff)
	}
}

func createTestingManifest(t *testing.T) *Manifest {
	shimsToAdd := []shimNameAndInfo{
		{
			name: "new-shim1",
			shim: Shim{
				Version:    "1234",
				Command:    "my-command1",
				Parameters: []string{"param1", "param2"},
			},
		},
		{
			name: "new-shim2",
			shim: Shim{
				Version:    "2345",
				Command:    "my-command2",
				Parameters: []string{"param1", "param2"},
			},
		},
		{
			name: "new-shim3",
			shim: Shim{
				Version:    "3456",
				Command:    "my-command3",
				Parameters: []string{"param1", "param2"},
			},
		},
	}

	m := CreateManifest(testSourceName)

	for _, shimToAdd := range shimsToAdd {
		assert.NoError(t, m.AddShim(shimToAdd.name, shimToAdd.shim), "should not error")
	}

	return m
}
