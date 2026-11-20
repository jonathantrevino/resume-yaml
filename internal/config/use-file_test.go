package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"resume-yaml/internal/config"
	"resume-yaml/internal/resume"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestUseFile(t *testing.T) {

	// simulate sending --file "resume-custom.yaml" to right location
	t.Run("use flag, file exists", func(t *testing.T) {

		tempDir := tempHelper(t)

		// create custom resume-custom.yaml for --file flag
		r := resume.Resume{
			Name: resume.Name{
				FirstName:  "Little",
				MiddleName: "Red",
				LastName:   "Fox",
			},
		}

		// marshal resume-config.yaml data
		data, err := yaml.Marshal(r)
		if err != nil {
			t.Fatalf("marshalling resume-config.yaml data failed")
		}

		customFilePath := filepath.Join(tempDir, "resume-custom.yaml")

		os.WriteFile(customFilePath, data, 0755)

		resultPath, err := config.UseFile(customFilePath)
		if err != nil {
			t.Errorf("use file: %v", err)
		}

		if resultPath != customFilePath {
			t.Errorf("expected %s got %s", customFilePath, resultPath)
		}

	})

	// simulate sending --file "missing-file.yaml" which is a wrong location
	t.Run("use flag, file missing", func(t *testing.T) {

		tempDir := tempHelper(t)
		_, err := config.UseFile(filepath.Join(tempDir, "missing-file.yaml"))

		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected fs.ErrNotExist got: %v", err)

		}
	})

	// use default resume.yaml after init
	t.Run("no flag, after init", func(t *testing.T) {

		tempDir := tempHelper(t)
		force := false
		vaultFlag := ""

		// init
		config.Init(force, vaultFlag)

		defaultFilePath := filepath.Join(tempDir, "resumes", "resume.yaml")

		resultPath, err := config.UseFile("")

		if err != nil {
			t.Errorf("use file: %v", err)
		}

		if resultPath != defaultFilePath {
			t.Errorf("expected %s got %s", defaultFilePath, resultPath)
		}
	})

	// don't initialize and don't send a flag expect run init
	t.Run("no flag, never initialized", func(t *testing.T) {
		_ = tempHelper(t)
		_, err := config.UseFile("")
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected fs.ErrNotExist got: %v", err)
		}
	})

	// run init, but delete resume.yaml and try usefile with no flag
	t.Run("no flag, resume deleted", func(t *testing.T) {
		tempDir := tempHelper(t)

		force := false
		vaultFlag := ""

		// init
		config.Init(force, vaultFlag)

		defaultFilePath := filepath.Join(tempDir, "resumes", "resume.yaml")

		if err := os.Remove(defaultFilePath); err != nil {
			t.Fatalf("deleting resume.yaml failed: %v", err)
		}

		_, err := config.UseFile("")

		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected fs.ErrNotExist got: %v", err)
		}

		// TODO: error message when resume.yaml was deleted
		// should include that the user must now run --force because config.yaml still exists
	})
}
