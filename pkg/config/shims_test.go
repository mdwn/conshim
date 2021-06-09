package config

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShimInfo(t *testing.T) {
	tests := []struct {
		name               string
		shimFileName       string
		shimFileContents   string
		expectedSource     string
		expectedVersion    string
		expectedParameters []string
		expectedCommand    string
	}{
		{
			name:         "shim file has all parts",
			shimFileName: "test",
			shimFileContents: `#!/usr/bin/env bash
# source: some-source version: 1234567 parameters: a,b,c
docker run container "$@"`,
			expectedSource:     "some-source",
			expectedVersion:    "1234567",
			expectedParameters: []string{"a", "b", "c"},
			expectedCommand:    "docker run container \"$@\"",
		},
		{
			name:         "shim file has no parameters",
			shimFileName: "test",
			shimFileContents: `#!/usr/bin/env bash
# source: some-source version: 1234567
docker run container "$@"`,
			expectedSource:  "some-source",
			expectedVersion: "1234567",
			expectedCommand: "docker run container \"$@\"",
		},
		{
			name:             "empty file",
			shimFileName:     "test",
			shimFileContents: ``,
			expectedSource:   "???",
			expectedVersion:  "???",
			expectedCommand:  shebangMissingErrorMessage,
		},
		{
			name:             "metadata missing",
			shimFileName:     "test",
			shimFileContents: `#!/usr/bin/env bash`,
			expectedSource:   "???",
			expectedVersion:  "???",
			expectedCommand:  metadataMissingErrorMessage,
		},
		{
			name:         "command missing",
			shimFileName: "test",
			shimFileContents: `#!/usr/bin/env bash
# source: some-source version: 1234567`,
			expectedSource:  "some-source",
			expectedVersion: "1234567",
			expectedCommand: commandMissingErrorMessage,
		},
	}

	for _, test := range tests {
		func() {
			tmpFileNeedsClose := false
			tmpFile, err := ioutil.TempFile("", "")
			assert.NoError(t, err, "shoudl be no error when creating temp file")

			tmpFileName := tmpFile.Name()

			defer func() {
				if tmpFileNeedsClose {
					assert.NoError(t, tmpFile.Close(), "should be no error when closing temp file")
				}
				assert.NoError(t, os.Remove(tmpFileName), "should be no error when removing temp file")
			}()

			_, err = io.WriteString(tmpFile, test.shimFileContents)
			assert.NoError(t, err, "should be no error when writing shim contents to temp file")

			_, err = tmpFile.Seek(0, 0)
			assert.NoError(t, err, "should be no error when resetting position of temp file")

			assert.NoError(t, tmpFile.Close(), "should be no error when closing temp file")

			shim := getShimFromFile(test.shimFileName, tmpFileName)

			assert.Equal(t, test.shimFileName, shim.Name, "sources should match")
			assert.Equal(t, test.expectedSource, shim.Source, "sources should match")
			assert.Equal(t, test.expectedVersion, shim.Version, "versions should match")
			assert.Equal(t, test.expectedParameters, shim.Parameters, "versions should match")
			assert.Equal(t, test.expectedCommand, shim.Command, "commands should match")
		}()
	}
}
