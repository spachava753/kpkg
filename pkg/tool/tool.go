package tool

import (
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

	installedVersion, err := InstalledVersion(basePath, binary)
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

// InstalledVersion checks if a binary is installed, and returns the version installed
// It returns an error if the binary is not installed
// The symlink we are checking could exist in three states, "existing", "broken", and "not existing"
// "existing" is a happy symlink. Symlink exists, and the link is not broken
// "broken" is a sad symlink. Symlink exists, and the link not broken
// "not existing" is depressed symlink. Symlink does not exist, but it's ok
func InstalledVersion(basePath, binary string) (string, error) {
	symPath := filepath.Join(basePath, "bin", binary)

	// symlink doesn't exist
	if _, err := os.Readlink(symPath); err != nil {
		return "", nil
	}
	// returns an err for broken symlink
	link, err := filepath.EvalSymlinks(filepath.Join(basePath, "bin", binary))
	if err != nil {
		return "", err
	}
	// returns version for happy symlink ;)
	dirs := strings.Split(link, string(os.PathSeparator))
	return dirs[len(dirs)-2], nil
}
