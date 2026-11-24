package testutil

import (
	"path/filepath"
	"testing"
)

// create temp folder and environment variables
// return tempDir and empty config struct
func TempHelper(t *testing.T) string {

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))

	return tempDir

}
