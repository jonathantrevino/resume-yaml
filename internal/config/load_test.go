package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"resume-yaml/internal/config"
	"resume-yaml/internal/testutil"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestLoad(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		// initialize
		tempDir := testutil.TempHelper(t)
		configFilePath := filepath.Join(tempDir, "config.yaml") // example config.yaml location
		vaultPath := filepath.Join(tempDir, "resumes")

		// set vault path
		initialConfig := config.Config{
			VaultPath: vaultPath,
		}

		// marshal config
		data, err := yaml.Marshal(&initialConfig)
		if err != nil {
			t.Fatalf("could not marshal config: %v", err)
		}

		// write yaml to config.yaml
		if err := os.WriteFile(configFilePath, data, 0600); err != nil {
			t.Fatalf("could not write config: %v", err)
		}

		// test load on fresh config
		var testConfig config.Config

		// test load
		if err := testConfig.Load(configFilePath); err != nil {
			t.Fatalf("load failed: %v", err)
		}

		// verify VaultPath
		if testConfig.VaultPath != initialConfig.VaultPath {

			t.Errorf("expected %s got %s", initialConfig.VaultPath, testConfig.VaultPath)
		}

	})

	t.Run("missing", func(t *testing.T) {

		// initialize
		var c config.Config
		tempDir := testutil.TempHelper(t)
		dir := filepath.Join(tempDir, "missing")

		// load random path
		err := c.Load(dir)

		// fail if it does not return ErrNotExist
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("not fs.ErrNotExist: %v", err)
		}

	})

	t.Run("broken", func(t *testing.T) {

		// initialize
		tempDir := testutil.TempHelper(t)
		configFilePath := filepath.Join(tempDir, "config.yaml") // example config.yaml location
		var c config.Config

		invalidYaml := []byte("invalid yaml")

		// write invalid yaml to config.yaml
		if err := os.WriteFile(configFilePath, invalidYaml, 0600); err != nil {
			t.Fatalf("could not write invalid config: %v", err)
		}

		err := c.Load(configFilePath)

		if err == nil {
			// fail if Load() on broken yaml did not break
			t.Fatalf("Load() did not break")
		}

		// fail if config file does not exist
		if errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("config intialization failed: %v", err)
		}

	})
}
