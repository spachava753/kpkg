package tool

import (
	"fmt"
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

	// Versions lists the possible installation candidates for a source like Github releases.
	// If the installedOnly flag is provided, only the installed version are shown
	Versions(installedOnly bool) ([]string, error)
}

// Remove will remove the binary version at the provided path
func Remove(binaryVersionPath string) error {
	// check whether path exists
	// delete version
	return nil
}

// Purge will remove all binary versions at the provided path
func Purge(binaryPath string) error {
	// check whether path exists
	// delete all versions

	return nil
}
