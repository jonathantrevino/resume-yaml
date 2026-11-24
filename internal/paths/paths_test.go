package paths_test

import (
	"errors"
	"path/filepath"
	"resume-yaml/internal/config"
	"resume-yaml/internal/paths"
	"testing"
)

func TestFindPaths(t *testing.T) {

	// fresh (no --file flag)
	t.Run("fresh", func(t *testing.T) {

		// initialize
		tempDir := pathsTempHelper(t)
		configDir, _ := config.ConfigFilePath()

		// expected path locations
		vaultFolderPath := filepath.Join(tempDir, "resumes")
		resumeFilePath := filepath.Join(tempDir, "resumes", "resume.yaml")

		emptyFile := ""
		p, err := paths.FindPaths(emptyFile)
		if err != nil {
			t.Fatalf("find paths: %v", err)
		}

		// verify config location
		if p.ConfigFile.Path != configDir {
			t.Errorf("expected config location %s got %s", configDir, p.ConfigFile.Path)
		}

		// check if config is already created
		if p.ConfigFile.Exists {
			t.Fatalf("config exists")
		}

		// verify vault location
		if p.VaultFolder.Path != vaultFolderPath {
			t.Errorf("expected vault location %s got %s", vaultFolderPath, p.VaultFolder.Path)
		}

		// check if vault is already created
		if p.VaultFolder.Exists {
			t.Fatalf("vault exists")
		}

		// make sure it is the default path
		if !p.VaultFolder.IsDefault {
			t.Errorf("vault location not default")

		}

		// verify resume location
		if p.ResumeFile.Path != resumeFilePath {
			t.Errorf("expected resume location %s got %s", resumeFilePath, p.ResumeFile.Path)
		}

		// make sure resume does not exist
		if p.ResumeFile.Exists {
			t.Fatalf("resume exists")
		}

		// verify set err is custom ErrNotInitialized
		if !errors.Is(p.ResumeFile.Err, config.ErrNotInitialized) {
			t.Errorf("expected ErrNotInitialized got %s", p.ResumeFile.Err)
		}

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
