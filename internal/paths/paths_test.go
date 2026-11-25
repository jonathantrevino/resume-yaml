package paths_test

import (
	"errors"
	"os"
	"path/filepath"
	"resume-yaml/internal/config"
	"resume-yaml/internal/paths"
	"resume-yaml/internal/testutil"
	"testing"
)

func TestFindPaths(t *testing.T) {

	// fresh (no --file flag)
	t.Run("fresh", func(t *testing.T) {

		// initialize
		tempDir := pathsTempHelper(t)

		// expected path locations
		configFilePath, _ := config.ConfigFilePath()
		vaultFolderPath := filepath.Join(tempDir, "resumes")
		resumeFilePath := filepath.Join(tempDir, "resumes", "resume.yaml")

		emptyFile := ""
		p, err := paths.FindPaths(emptyFile)
		if err != nil {
			t.Fatalf("find paths: %v", err)
		}

		// verify config location
		if p.ConfigFile.Path != configFilePath {
			t.Errorf("expected config location %s got %s", configFilePath, p.ConfigFile.Path)
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

		// TODO: Add checks for IsDefault

		// initialize
		tempDir := testutil.TempHelper(t)
		force := false
		vaultFlag := ""

		config.Init(force, vaultFlag)

		// expected path locations
		configFilePath, _ := config.ConfigFilePath()
		vaultFolderPath := filepath.Join(tempDir, "resumes")
		resumeFilePath := filepath.Join(tempDir, "resumes", "resume.yaml")

		emptyFile := ""
		p, err := paths.FindPaths(emptyFile)
		if err != nil {
			t.Fatalf("find paths: %v", err)
		}

		// verify config location
		if p.ConfigFile.Path != configFilePath {
			t.Errorf("expected config location %s got %s", configFilePath, p.ConfigFile.Path)
		}

		// check if config is created
		if !p.ConfigFile.Exists {
			t.Errorf("config not initialized")
		}

		// verify vault location
		if p.VaultFolder.Path != vaultFolderPath {
			t.Errorf("expected vault location %s got %s", vaultFolderPath, p.VaultFolder.Path)
		}

		// check if vault is created
		if !p.VaultFolder.Exists {
			t.Errorf("vault not initialized")
		}

		// verify resume location
		if p.ResumeFile.Path != resumeFilePath {
			t.Errorf("expected resume location %s got %s", vaultFolderPath, p.VaultFolder.Path)
		}

		// check if resume.yaml was created
		if !p.ResumeFile.Exists {
			t.Errorf("resume not initialized")
		}

		// check if resume is default location
		if !p.ResumeFile.IsDefault {
			t.Errorf("resume is not default location")
		}
	})

	// with --file flag
	t.Run("with --file flag", func(t *testing.T) {

		// initialize
		tempDir := testutil.TempHelper(t)
		force := false
		vaultPath := ""

		config.Init(force, vaultPath)

		// expected path locations
		configFilePath, _ := config.ConfigFilePath()
		vaultFolderPath := filepath.Join(tempDir, "resumes")
		file := "custom-resume.yaml"
		resumeFilePath := filepath.Join(tempDir, "resumes", file)

		// create resume file
		if err := os.WriteFile(resumeFilePath, []byte("testing"), 0755); err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		// pass in resumeFilePath instead of file because FindPath expects a full path
		p, err := paths.FindPaths(resumeFilePath)

		if err != nil {
			t.Fatalf("find paths: %v", err)
		}

		// verify config location
		if p.ConfigFile.Path != configFilePath {
			t.Errorf("expected config location %s got %s", configFilePath, p.ConfigFile.Path)
		}

		// check if config is created
		if !p.ConfigFile.Exists {
			t.Errorf("config is not initialized")

		}

		// verify vault location
		if p.VaultFolder.Path != vaultFolderPath {
			t.Errorf("expected vault location %s got %s", vaultFolderPath, p.VaultFolder.Path)
		}

		// check if vault is created
		if !p.VaultFolder.Exists {
			t.Errorf("vault is not initialized")

		}

		// verify resume location
		if p.ResumeFile.Path != resumeFilePath {
			t.Errorf("expected resume location %s got %s", resumeFilePath, p.ResumeFile.Path)
		}

		// check if resume is created
		if !p.ResumeFile.Exists {
			t.Errorf("resume is not initialized")
		}

		// check if resume is default location
		if p.ResumeFile.IsDefault {
			t.Errorf("resume is default location (not set --file)")
		}
	})

	// TODO: Test cases to add:
	// - after init with flag to missing file
	// - after init but delete config manually
}

// create temp folder and environment variables
// return tempDir and empty config struct
func pathsTempHelper(t *testing.T) string {

	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))

	return tempDir

}
