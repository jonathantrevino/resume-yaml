package config_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"resume-yaml/internal/config"
	"resume-yaml/internal/testutil"
	"testing"
)

func TestInit(t *testing.T) {

	// fresh init run
	// expected: sets vault file path to tempDir,
	// outputs and resume.yaml should exist at location
	t.Run("fresh init", func(t *testing.T) {
		fmt.Println("inhere")

		// initalize
		var c config.Config
		tempDir := testutil.TempHelper(t)
		force := false
		vaultFlag := "" // empty because we want to use tempDir instead
		vaultPath := filepath.Join(tempDir, "resumes")

		if err := config.Init(force, vaultFlag); err != nil {
			// test depends on Init working
			t.Fatalf("init failed: %v", err)
		}

		configFilePath, err := config.ConfigFilePath()
		if err != nil {
			t.Fatalf("configFilePath failed: %v", err)
		}

		if err := c.Load(configFilePath); err != nil {
			// test depends on Load working
			t.Fatalf("load failed: %v", err)
		}

		// vault path not set
		if c.VaultPath != vaultPath {
			t.Errorf("expected %s got %s", tempDir, c.VaultPath)
		}

		resumeFile := filepath.Join(vaultPath, "resume.yaml")
		tInfo, err := os.Stat(resumeFile)
		if err == nil {
			// resume is a directory
			if tInfo.IsDir() {
				t.Errorf("resume is dir not file")
			}
		} else if errors.Is(err, os.ErrNotExist) {
			// resume does not exist
			t.Errorf("resume does not exist")
		} else {
			t.Errorf("unexpected error: %v", err)
		}

		outputDir := filepath.Join(vaultPath, "output")

		oInfo, err := os.Stat(outputDir)
		if err == nil {
			// output is a file
			if !oInfo.IsDir() {
				t.Errorf("output is file not dir")
			}
		} else if errors.Is(err, os.ErrNotExist) {
			// output folder does not exist
			t.Errorf("output does not exist")
		} else {
			t.Errorf("unexpected error: %v", err)
		}

	})

	// refuses init run
	// expected: alerts --force command must be used because of
	// existing config and or resume.yaml file
	t.Run("refuse init", func(t *testing.T) {
		// initalize
		tempDir := testutil.TempHelper(t)
		force := false
		vaultFlag := "" // empty because we want to use tempDir instead
		vaultPath := filepath.Join(tempDir, "resumes")

		// first t.Run("fresh init") will flag if init is not working, assume it is

		// run Init first time
		if err := config.Init(force, vaultFlag); err != nil {
			// test depends on Init working
			t.Fatalf("init failed 1: %v", err)
		}

		resumeFile := filepath.Join(vaultPath, "resume.yaml")

		// write data to resume to check if it was overwritten
		marker := []byte("editing resume")
		err := os.WriteFile(resumeFile, marker, 0644)
		if err != nil {
			t.Fatalf("write marker: %v", err)
		}

		// run Init second time to force alert
		if err := config.Init(force, vaultFlag); err != nil {
			// check if err is alert
			if !errors.Is(err, config.ErrAlreadyExists) {
				t.Errorf("expected ErrAlreadyExists, got %v", err)
			}
		}

		// check if resume was overwritten
		got, err := os.ReadFile(resumeFile)
		if err != nil {
			t.Fatalf("read resume: %v", err)
		}

		if string(got) != string(marker) {
			t.Errorf("file was overwritten")
		}

	})

	// force init run
	// expected: should overwrite config and or resume.yaml file
	t.Run("force init", func(t *testing.T) {
		// initalize
		tempDir := testutil.TempHelper(t)
		force := false
		vaultFlag := "" // empty because we want to use tempDir instead
		vaultPath := filepath.Join(tempDir, "resumes")

		// first t.Run("fresh init") will flag if init is not working, assume it is

		// run Init first time
		if err := config.Init(force, vaultFlag); err != nil {
			// test depends on Init working
			t.Fatalf("init failed 1: %v", err)
		}

		resumeFile := filepath.Join(vaultPath, "resume.yaml")

		// write data to resume to check if it was overwritten
		marker := []byte("editing resume")
		err := os.WriteFile(resumeFile, marker, 0644)
		if err != nil {
			t.Fatalf("write marker: %v", err)
		}

		// run Init second time with --force flag
		force = true
		if err := config.Init(force, vaultFlag); err != nil {
			// check if err is alert
			if !errors.Is(err, config.ErrAlreadyExists) {
				t.Errorf("expected ErrAlreadyExists, got %v", err)
			}
		}

		// check if resume was overwritten
		got, err := os.ReadFile(resumeFile)
		if err != nil {
			t.Fatalf("read resume: %v", err)
		}

		if string(got) == string(marker) {
			t.Errorf("file was not overwritten")
		}

	})

}
