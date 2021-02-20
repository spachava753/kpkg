package error

import (
	"fmt"
	"runtime"
)

type InstalledErr struct {
	Version    string
	BinaryName string
}

func (i InstalledErr) Error() string {
	return fmt.Sprintf("binary %s with %s version is already installed", i.BinaryName, i.Version)
}

type UnsupportedRuntimeErr struct {
	Binary string
}

func (u *UnsupportedRuntimeErr) Error() string {
	return fmt.Sprintf("downloading binary %s is not support with runtime %s/%s", u.Binary, runtime.GOOS, runtime.GOARCH)
}
