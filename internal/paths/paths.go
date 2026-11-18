package paths

import (
	"fmt"
	"resume-yaml/internal/config"
)

type Paths struct {
	ConfigFilePath  string
	VaultFolderPath string
	ResumeFilePath  string
}

func FindPaths(file string) (Paths, error) {

	var p Paths
	var c config.Config

	// get config file location
	configFilePath, err := config.ConfigFilePath()
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
	resumeFilePath, err := config.UseFile(file)
	if err != nil {
		return p, fmt.Errorf("use file: %w", err)
	}
	p.ResumeFilePath = resumeFilePath

	return p, nil
}
