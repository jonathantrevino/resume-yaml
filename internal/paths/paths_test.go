package paths_test

import (
	"path/filepath"
	"testing"
)

func TestFindPaths(t *testing.T) {

	// fresh
	t.Run("fresh", func(t *testing.T) {

	})

	// after init
	t.Run("after init", func(t *testing.T) {

	})

	// with --file flag
	t.Run("with --file flag", func(t *testing.T) {

	})
}

// create temp folder and environment variables
// return tempDir and empty config struct
func pathsTempHelper(t *testing.T) string {

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))

	return tempDir

}
