package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var ErrNotInitialized = errors.New("resume.yaml not found; run 'resume init' to set up your vault")
var ErrAlreadyExists = errors.New("already exists")

type Config struct {
	VaultPath string `yaml:"vault_path"`
}

// return path of config file
func ConfigFilePath() (string, error) {

	// get users config directory
	config, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get config directory: %w", err)
	}

	configFilePath := filepath.Join(config, "resume-yaml", "config.yaml")

	return configFilePath, nil

}

// return path of home directory
func DefaultVaultPath() (string, error) {

	// get users home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	vaultPath := filepath.Join(home, "resumes")
	return vaultPath, nil
}

// ensureAbsent returns nil if nothing exists at path, or an error
func ensureAbsent(path, name string) error {
	info, err := os.Stat(path)
	switch {
	case err == nil && info.IsDir():
		return fmt.Errorf("%s is a directory, expected a %s file", path, name)
	case err == nil:
		return fmt.Errorf("%s %w at %s, use --force to overwrite", name, ErrAlreadyExists, path)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		return fmt.Errorf("check %s: %w", path, err)
	}
}

// check in order which file (resume) to use
func resolvePath(file string) (string, error) {

	// 1. user has set --file / -f flag
	if file != "" {
		return file, nil
	}
	// 2. RESUME_YAML_FILE environment variable is set
	if envFile := os.Getenv("RESUME_YAML_FILE"); envFile != "" {
		return envFile, nil
	}

	// 3. Use default vault path from config.yaml

	configFilePath, err := ConfigFilePath()
	if err != nil {
		return "", fmt.Errorf("get config file: %w", err)
	}

	// load config
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("missing config at %s, %w", configFilePath, err)
		}
		return "", fmt.Errorf("reading config file: %w", err)
	}

	// unmarshal data
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return "", fmt.Errorf("unmarshalling config data: %w", err)
	}

	if c.VaultPath == "" {
		return "", fmt.Errorf("vault path not set in config %w", ErrNotInitialized)
	}

	templateFilePath := filepath.Join(c.VaultPath, "resume.yaml")

	return templateFilePath, nil
}
