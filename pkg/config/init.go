package config

import (
	"os"
	"path"
)

func CreateBinPath(basePath string) error {
	// check that base path is valid
	_, err := os.Stat(basePath)
	if err != nil {
		return err
	}

	rootDirPath := path.Join(basePath, ".kpkg")
	binPath := path.Join(rootDirPath, "bin")
	return os.MkdirAll(binPath, os.ModePerm)
}
