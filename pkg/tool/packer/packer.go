package packer

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type packerTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l packerTool) Extract(artifactPath, _ string) (string, error) {
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

func (l packerTool) Name() string {
	return "packer"
}

func (l packerTool) ShortDesc() string {
	return "Packer is a tool for creating identical machine images for multiple platforms from a single source configuration"
}

func (l packerTool) LongDesc() string {
	return `Packer is lightweight, runs on every major operating system, and is highly performant, 
creating machine images for multiple platforms in parallel. Packer comes out of the box with support 
for many platforms, the full list of which can be found at https://www.packer.io/docs/builders.`
}

func (l packerTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	c, err := semver.NewConstraint("< 1.4.1")
	if err != nil {
		return "", err
	}

	if c.Check(v) && l.os == "linux" {
		switch {
		case l.arch == "mips",
			l.arch == "mips64",
			l.arch == "mipsle",
			l.arch == "ppc64le",
			l.arch == "s390x":
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "386",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "arm",
		l.os == "freebsd" && l.arch == "386",
		l.os == "openbsd" && l.arch == "amd64",
		l.os == "openbsd" && l.arch == "386",
		l.os == "solaris" && l.arch == "amd64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_%s_%s.zip", version, version, l.os, l.arch)
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return packerTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("hashicorp", "packer", 20),
	}
}
