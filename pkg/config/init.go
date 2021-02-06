package config

import (
	"os"
	"path/filepath"
)

func CreatePath(basePath string) (string, error) {
	// check that base path is valid
	_, err := os.Stat(basePath)
	if err != nil {
		return "", err
	}

	rootDirPath := filepath.Join(basePath, ".kpkg")
	binPath := filepath.Join(rootDirPath, "bin")
	if err = os.MkdirAll(binPath, os.ModePerm); err != nil {
		return "", err
	}

	return rootDirPath, nil
}
