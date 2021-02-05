package tool

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InstalledErr struct {
	Version    string
	BinaryName string
}

func (i InstalledErr) Error() string {
	return fmt.Sprintf("binary %s with %s version is already installed", i.BinaryName, i.Version)
}

// Binary is an interface that all binaries must implement
type Binary interface {
	// Install will first check to see if the version is already downloaded.
	// If it is, it will only download again if the force flag is enabled.
	// If it is not downloaded, it will download the binary.
	// After checking/downloaded, it will form the symlink.
	// It returns the path to symlink and an error if applicable
	Install(version string, force bool) (string, error)

	// Versions lists the possible installation candidates for a source like Github releases
	Versions() ([]string, error)
}

// RemoveVersions will remove the binary version at the provided path
// basePath is the path where the .kpkg folder is located
// binary is the binary name
// versions is a list of versions to remove
func RemoveVersions(basePath string, binary string, versions []string) error {
	// check that supplied versions are valid
	if versions == nil || len(versions) == 0 {
		return fmt.Errorf("not enough versions were passed in")
	}

	installedVersion, err := LinkedVersion(basePath, binary)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if installedVersion == v {
			return fmt.Errorf("cannot uninstalled version %s, currently in use. Please install another version first", v)
		}
		if err := os.RemoveAll(filepath.Join(basePath, v)); err != nil {
			return err
		}
	}
	return nil
}

// Purge will remove all binary versions at the provided path
func Purge(path string) error {
	return os.RemoveAll(path)
}

// LinkedVersion checks if a binary is symlinked, and returns the version symlinked
// It returns an error if the binary is not symlinked
// The symlink we are checking could exist in three states, "existing", "broken", and "not existing"
// "existing" is a happy symlink. Symlink exists, and the link is not broken
// "broken" is a sad symlink. Symlink exists, and the link not broken
// "not existing" is depressed symlink. Symlink does not exist, but it's ok
func LinkedVersion(basePath, binary string) (string, error) {
	// check that a valid binary was passed in
	if binary == "" {
		return "", errors.New("binary name cannot be empty")
	}

	// check that given base path exists
	info, err := os.Stat(basePath)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("given path is not a dir: %s", basePath)
	}

	symPath := filepath.Join(basePath, "bin", binary)

	// symlink doesn't exist
	if _, err := os.Readlink(symPath); err != nil {
		return "", nil
	}
	// returns an err for broken symlink
	linkPath, err := filepath.EvalSymlinks(filepath.Join(basePath, "bin", binary))
	if err != nil {
		return "", err
	}
	// returns version for happy symlink ;)
	dirs := strings.Split(linkPath, string(os.PathSeparator))
	return dirs[len(dirs)-2], nil
}

// Installed checks if a binary version is already downloaded.
// While walking through the paths, it also checks to make sure that
// the path is valid. If a binary is found at the end of the path,
// it returns true
func Installed(basePath, binary, version string) (bool, error) {
	binaryPath := filepath.Join(basePath, binary)
	binaryVersionPath := filepath.Join(binaryPath, version)
	binaryFilePath := filepath.Join(binaryVersionPath, binary)

	if ld2Info, err := os.Stat(binaryPath); err == nil {
		if !ld2Info.IsDir() {
			return false, fmt.Errorf("path %s contains a file", binaryPath)
		}
		if ld2VersionInfo, err := os.Stat(binaryVersionPath); err == nil {
			if !ld2VersionInfo.IsDir() {
				return false, fmt.Errorf("path %s contains a file", binaryVersionPath)
			}
			if ld2BinaryInfo, err := os.Stat(binaryFilePath); err == nil {
				if ld2BinaryInfo.IsDir() {
					return false, fmt.Errorf("path %s contains a dir", binaryFilePath)
				}
				return true, nil
			}
		}
	}
	return false, nil
}
