package resume

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {

	var r Resume

	fileLocations := []struct {
		name     string
		file     string
		expected error
	}{
		{
			name:     "file exists",
			file:     "valid-template.yaml",
			expected: nil,
		},
		{
			name:     "file does not exist",
			file:     "missing-template.yaml",
			expected: fs.ErrNotExist,
		},
	}

	for _, tc := range fileLocations {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", tc.file)
			result := r.Load(path)
			if !errors.Is(result, tc.expected) {
				t.Errorf("expected: %d, got: %d", tc.expected, result)
			}
		})
	}
}
