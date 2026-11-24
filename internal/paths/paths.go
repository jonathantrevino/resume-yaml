package paths

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"resume-yaml/internal/config"
)

type Paths struct {
	ConfigFile  Location
	VaultFolder Location
	ResumeFile  Location
}

type Location struct {
	Path      string
	Exists    bool
	IsDefault bool
	Err       error
}

func FindPaths(file string) (Paths, error) {

	var p Paths
	var c config.Config

	// get config file location
	configFilePath, err := config.ConfigFilePath()

	switch err {
	case nil:
		p.ConfigFile.Exists = true
	case fs.ErrNotExist:
		p.ConfigFile.Exists = false
	default:
		p.ConfigFile.Err = err

	}
	// set config file location
	p.ConfigFile.Path = configFilePath

	// check if config file exists
	if _, err := os.Stat(configFilePath); err != nil {

		// file does not exist
		if errors.Is(err, fs.ErrNotExist) {

			p.ConfigFile.Exists = false
			p.ConfigFile.Err = err
		}

	}

	// load config
	// to get vault folder location
	if err := c.Load(configFilePath); err != nil {

		if errors.Is(err, fs.ErrNotExist) {
			p.ConfigFile.Exists = false
		}

	}

	// use vault path from config
	p.VaultFolder.Path = c.VaultPath

	// if config does not exist, then use default vault path
	if !p.ConfigFile.Exists {

		p.VaultFolder.Path, _ = config.DefaultVaultPath()
		p.VaultFolder.IsDefault = true

	}

	p.VaultFolder.Exists = true

	// check if vault path exists
	if _, err := os.Stat(p.VaultFolder.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			p.VaultFolder.Exists = false
		}
	}

	// work out the resume path
	resumePath := file
	if resumePath == "" {
		resumePath = filepath.Join(p.VaultFolder.Path, "resume.yaml")
	}
	p.ResumeFile.Path = resumePath

	// does it exist
	info, err := os.Stat(resumePath)
	p.ResumeFile.Exists = err == nil && !info.IsDir()

	// err is nil if usable
	_, p.ResumeFile.Err = config.UseFile(file)

	return p, nil
}
