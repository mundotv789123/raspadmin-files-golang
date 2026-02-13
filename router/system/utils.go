package system

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mundotv789123/raspadmin/internal/config"
)

func SafeJoinPath(path string) (string, error) {
	fullPath := filepath.Join(config.AbsRootDir, path)
	cleanedPath := filepath.Clean(fullPath)
	absFullPath, err := filepath.Abs(cleanedPath)

	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, config.AbsRootDir) {
		return "", fmt.Errorf("path traversal attempt detected: %s", path)
	}

	return absFullPath, nil
}
