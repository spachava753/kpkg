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
	_, err = os.Stat(rootDirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		binPath := path.Join(rootDirPath, "bin")
		return os.MkdirAll(binPath, os.ModePerm)
	}

	return nil
}
