package get

import "fmt"

type InstalledErr struct {
	Version    string
	BinaryName string
}

func (i InstalledErr) Error() string {
	return fmt.Sprintf("binary %s with %s version is already installed", i.BinaryName, i.Version)
}

type Installer interface {
	// Install will check to see if the version is already installed.
	// If it is, it will only install if the force flag is enabled.
	// If it is not installed, it will download the binary and install it
	Install(version string, force bool) error
}

type VersionLister interface {
	// Versions lists the possible installation candidates for a source like Github releases.
	// If the installedOnly flag is provided, only the installed version are shown
	Versions(installedOnly bool) ([]string, error)
}

type Tool interface {
	Installer
	VersionLister
}
