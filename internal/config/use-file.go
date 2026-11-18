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
