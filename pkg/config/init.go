package config

import (
	"os"
	"path"
)

func CreatePath(basePath string) (string, error) {
	// check that base path is valid
	_, err := os.Stat(basePath)
	if err != nil {
		return "", err
	}

	rootDirPath := path.Join(basePath, ".kpkg")
	binPath := path.Join(rootDirPath, "bin")
	if err = os.MkdirAll(binPath, os.ModePerm); err != nil {
		return "", err
	}

	return rootDirPath, nil
}
