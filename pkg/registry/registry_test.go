package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMungeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ssh URL",
			input:    "git@github.com:some/repo",
			expected: "git@github.com:some/repo",
		},
		{
			name:     "http URL",
			input:    "http://github.com/some/repo",
			expected: "http://github.com/some/repo",
		},
		{
			name:     "https URL",
			input:    "https://github.com/some/repo",
			expected: "https://github.com/some/repo",
		},
		{
			name:     "mystery schema URL",
			input:    "myst://github.com/some/repo",
			expected: "myst://github.com/some/repo",
		},
		{
			name:     "no schema URL",
			input:    "github.com/some/repo",
			expected: "https://github.com/some/repo",
		},
		{
			name:     "no schema but ://",
			input:    "://github.com/some/repo",
			expected: "://github.com/some/repo",
		},
	}

	for _, test := range tests {
		output := mungeURL(test.input)
		assert.Equal(t, test.expected, output, "URLs should be the same")
	}
}
