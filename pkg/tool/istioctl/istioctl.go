package istioctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type istioctlTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l istioctlTool) Extract(artifactPath, version string) (string, error) {
	// istioctl releases contain the binary and some examples. Pick only the binary

	// expect given path to be to contain bin folder
	binDirPath := filepath.Join(artifactPath, fmt.Sprintf("istio-%s", version), "bin")
	binDirPathInfo, err := os.Stat(binDirPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if !binDirPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}
	binaryPath := filepath.Join(binDirPath, "istioctl")
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}

	return binaryPath, err
}

func (l istioctlTool) Name() string {
	return "istioctl"
}

func (l istioctlTool) ShortDesc() string {
	return "Connect, secure, control, and observe services"
}

func (l istioctlTool) LongDesc() string {
	return `Istio is an open platform for providing a uniform way to integrate microservices, 
manage traffic flow across microservices, enforce policies and aggregate telemetry data. 
Istio's control plane provides an abstraction layer over the underlying cluster management platform, 
such as Kubernetes.`
}

func (l istioctlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	c, err := semver.NewConstraint("< 1.6.0")
	if err != nil {
		return "", err
	}
	url := l.MakeReleaseUrl()

	// based of the script "curl -L https://istio.io/downloadIstio"
	if c.Check(v) && l.os == "linux" {
		// version is less than 1.5, have to use different url for downloading
		if l.arch != "amd64" {
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
		url += fmt.Sprintf("%s/istio-%s-linux.tar.gz", version, version)
		return url, nil
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += fmt.Sprintf("%s/istio-%s-osx.tar.gz", version, version)
	case l.os == "windows" && l.arch == "amd64":
		url += fmt.Sprintf("%s/istio-%s-win.zip", version, version)
	case l.os == "linux" && l.arch == "armv7":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		url += fmt.Sprintf("%s/istio-%s-linux-%s.tar.gz", version, version, l.arch)
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return istioctlTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("istio", "istio", 20),
	}
}
