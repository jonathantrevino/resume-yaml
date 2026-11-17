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

//go:embed resume-template.yaml
var defaultTemplate []byte

type Config struct {
	VaultPath string `yaml:"vault_path"`
}

// Create ~/.config/resume-yaml/config.yaml for config file
// Create ~/resumes folder (vault) for storing template and generated resumes
func Init(force bool, vault string) error {

	// get users home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}
	config, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	configPath := filepath.Join(config, ".config", "resume-yaml")
	configFilePath := filepath.Join(configPath, "config.yaml")
	vaultPath := filepath.Join(home, "resumes")

	// conditionally set vault path if user set --vault location
	trueVaultPath := vaultPath
	if vault != "" {
		trueVaultPath = vault
	}
	templatePath := filepath.Join(trueVaultPath, "resume.yaml")

	// if --force flag is set to false
	// check if config file and resume template are
	// already created to avoid overwriting and alert user
	if !force {
		if err := ensureAbsent(configFilePath, "config.yaml"); err != nil {
			return err
		}
		if err := ensureAbsent(templatePath, "resume.yaml"); err != nil {
			return err
		}
	}

	cfg := Config{
		VaultPath: trueVaultPath,
	}

	fmt.Printf("Creating config folder at %s\n", configPath)
	// create config folder
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("create config file: %w", err)
	}

	// marshall config to yaml
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshall config: %w", err)
	}

	fmt.Printf("Writing config file at %s\n", configFilePath)
	// write config file with defaults
	if err := os.WriteFile(configFilePath, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("Creating resumes folder at %s\n", trueVaultPath)
	// create resumes folder (vault)
	if err := os.MkdirAll(trueVaultPath, 0755); err != nil {
		return fmt.Errorf("create vault path: %w", err)
	}

	fmt.Printf("Copying default resume template to %s\n", templatePath)
	// write default resume template in vault
	if err := os.WriteFile(templatePath, defaultTemplate, 0755); err != nil {
		return fmt.Errorf("write template file: %w", err)
	}

	// create output folder
	// TODO: add config setting for output folder name
	outputPath := filepath.Join(trueVaultPath, "output")
	fmt.Printf("Creating output folder at %s\n", outputPath)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("create output path: %w", err)
	}

	return nil
}

// ensureAbsent returns nil if nothing exists at path, or an error
func ensureAbsent(path, name string) error {
	info, err := os.Stat(path)
	switch {
	case err == nil && info.IsDir():
		return fmt.Errorf("%s is a directory, expected a %s file", path, name)
	case err == nil:
		return fmt.Errorf("%s already exists at %s, use --force to overwrite", name, path)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		return fmt.Errorf("check %s: %w", path, err)
	}
}
