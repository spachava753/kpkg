package dive

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type diveTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l diveTool) Extract(artifactPath, _ string) (string, error) {
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binaryPathInfo))
	}

	return binaryPath, err
}

func (l diveTool) Name() string {
	return "dive"
}

func (l diveTool) ShortDesc() string {
	return "A tool for exploring each layer in a docker image"
}

func (l diveTool) LongDesc() string {
	return "A tool for exploring a docker image, layer contents, and discovering ways to shrink the size of your Docker/OCI image."
}

func (l diveTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_darwin_amd64.tar.gz
	url := fmt.Sprintf("%sv%s/dive_%s_", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "darwin_amd64.tar.gz"
	case l.os == "windows" && l.arch == "amd64":
		url += "windows_amd64.zip"
	case l.os == "linux" && l.arch == "amd64":
		url += "linux_amd64.tar.gz"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return diveTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("wagoodman", "dive", 20),
	}
}
