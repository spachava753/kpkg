package eksctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type eksctlTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l eksctlTool) Extract(artifactPath, _ string) (string, error) {
	binPath := filepath.Join(artifactPath, l.Name())
	info, err := os.Stat(binPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("expected a binary at path %s, found a directory", binPath)
	}
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l eksctlTool) Name() string {
	return "eksctl"
}

func (l eksctlTool) ShortDesc() string {
	return "The official CLI for Amazon EKS"
}

func (l eksctlTool) LongDesc() string {
	return `eksctl is a simple CLI tool for creating clusters on EKS - Amazon's new managed Kubernetes service for EC2. 
It is written in Go, and uses CloudFormation`
}

func (l eksctlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// add support for armv7 and armv6

	url := fmt.Sprintf("%s%s/eksctl_", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "Darwin_amd64"
	case l.os == "darwin" && l.arch == "arm64":
		url += "Darwin_arm64"
	case l.os == "windows" && l.arch == "amd64":
		url += "Windows_amd64"
	case l.os == "linux" && l.arch == "amd64":
		url += "Linux_amd64"
	case l.os == "linux" && l.arch == "arm64":
		url += "Linux_arm64"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return eksctlTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("weaveworks", "eksctl", 20),
	}
}
