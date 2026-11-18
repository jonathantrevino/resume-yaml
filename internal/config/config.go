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

//go:embed resume.yaml
var defaultTemplate []byte

type Config struct {
	VaultPath string `yaml:"vault_path"`
}

type Paths struct {
	ConfigFilePath  string
	VaultFolderPath string
	ResumeFilePath  string
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

// unmarshal config.yaml
func (c *Config) Load(p string) error {

	data, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("reading config file %s: %w", p, err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return err
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

// verify resolvedPath to file (resume) exists
func UseFile(file string) (string, error) {
	resolvedPath, err := resolvePath(file)
	info, err := os.Stat(resolvedPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("%w: %s", ErrNotInitialized, resolvedPath)
		}

		return "", fmt.Errorf("checking resume file: %w", err)

	}

	if info.IsDir() {
		return "", fmt.Errorf("expected a file, got a directory: %s", resolvedPath)
	}

	return resolvedPath, nil
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

func FindPaths(file string) (Paths, error) {

	var p Paths
	var c Config

	// get config file location
	configFilePath, err := ConfigFilePath()
	if err != nil {
		return p, fmt.Errorf("get config file path: %w", err)
	}

	p.ConfigFilePath = configFilePath

	// load config
	// to get vault folder location
	if err := c.Load(configFilePath); err != nil {

		return p, fmt.Errorf("get vault folder %w", err)
	}

	vaultFolderPath := c.VaultPath
	p.VaultFolderPath = vaultFolderPath

	// get resume.yaml file location
	resumeFilePath, err := UseFile(file)
	if err != nil {
		return p, fmt.Errorf("use file: %w", err)
	}
	p.ResumeFilePath = resumeFilePath

	return p, nil
}
