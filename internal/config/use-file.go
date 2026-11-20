package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// verify resolvedPath to file (resume) exists
func UseFile(file string) (string, error) {
	resolvedPath, err := resolvePath(file)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolvedPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {

			// --file was set but it missing
			if file != "" {
				return "", fmt.Errorf("resume file not found at %s: %w", resolvedPath, err)
			}

			// no --file set so vault resume.yaml is missing
			return "", fmt.Errorf("%w: %s: %w", ErrNotInitialized, resolvedPath, err)
		}

		return "", fmt.Errorf("checking resume file: %w", err)

	}

	if info.IsDir() {
		return "", fmt.Errorf("expected a file, got a directory: %s", resolvedPath)
	}

	return resolvedPath, nil
}
