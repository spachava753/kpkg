package polaris

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type polarisTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l polarisTool) Extract(artifactPath, _ string) (string, error) {
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

func (l polarisTool) Name() string {
	return "polaris"
}

func (l polarisTool) ShortDesc() string {
	return "Validation of best practices in your Kubernetes clusters"
}

func (l polarisTool) LongDesc() string {
	return "Fairwinds' Polaris keeps your clusters sailing smoothly. It runs a variety of checks to ensure that Kubernetes pods and controllers are configured using best practices, helping you avoid problems in the future"
}

func (l polarisTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/FairwindsOps/polaris/releases/download/4.0.4/polaris_4.0.4_darwin_amd64.tar.gz
	url := fmt.Sprintf("%s%s/polaris_%s_", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += l.os + "_" + l.arch
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64":
		url += l.os + "_" + l.arch
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return polarisTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("FairwindsOps", "polaris", 20),
	}
}
