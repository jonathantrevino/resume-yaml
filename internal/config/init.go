package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

//go:embed resume.yaml
var defaultTemplate []byte

// Create ~/.config/resume-yaml/config.yaml for config file
// Create ~/resumes folder (vault) for storing template and generated resumes
func Init(force bool, vaultFlag string) error {

	// get default paths
	configFilePath, err := ConfigFilePath()
	if err != nil {
		return fmt.Errorf("get config file path: %w", err)
	}

	// get config directory path
	configPath := filepath.Dir(configFilePath)
	vaultPath, err := DefaultVaultPath()
	if err != nil {
		return fmt.Errorf("get vault path: %w", err)
	}

	// conditionally set vault path if user set --vault location
	if vaultFlag != "" {
		vaultPath = vaultFlag
	}

	// use conditional vault path and store resume.yaml there
	templatePath := filepath.Join(vaultPath, "resume.yaml")

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
		VaultPath: vaultPath,
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

	fmt.Printf("Creating resumes folder at %s\n", vaultPath)
	// create resumes folder (vault)
	if err := os.MkdirAll(vaultPath, 0755); err != nil {
		return fmt.Errorf("create vault path: %w", err)
	}

	fmt.Printf("Copying default resume template to %s\n", templatePath)
	// write default resume template in vault
	if err := os.WriteFile(templatePath, defaultTemplate, 0755); err != nil {
		return fmt.Errorf("write template file: %w", err)
	}

	// create output folder
	// TODO: add config setting for output folder name
	outputPath := filepath.Join(vaultPath, "output")
	fmt.Printf("Creating output folder at %s\n", outputPath)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("create output path: %w", err)
	}

	return nil
}
