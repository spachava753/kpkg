package config

import (
	"errors"
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
		var pathError *os.PathError
		if !errors.As(err, &pathError) {
			return err
		}

		return os.Mkdir(rootDirPath, os.ModePerm)
	}

	return nil
}
