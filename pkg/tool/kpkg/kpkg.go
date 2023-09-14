package kpkg

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type kpkgTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kpkgTool) Extract(artifactPath, _ string) (string, error) {
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf(
			"could not extract binary: %w",
			fmt.Errorf("path %s is not a directory", binaryPathInfo),
		)
	}

	return binaryPath, err
}

func (l kpkgTool) Name() string {
	return "kpkg"
}

func (l kpkgTool) ShortDesc() string {
	return "A binary to install various K8s ecosystem related binaries"
}

func (l kpkgTool) LongDesc() string {
	return "A tool to install various K8s ecosystem related binaries.\n\n"
}

func (l kpkgTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/kpkg_", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "dragonfly",
		l.os == "freebsd",
		l.os == "linux",
		l.os == "netbsd",
		l.os == "openbsd",
		l.os == "solaris",
		l.os == "windows",
		l.os == "darwin":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	switch {
	case l.arch == "amd64":
	case l.arch == "386":
	case l.arch == "arm":
	case l.arch == "arm64":
	case l.arch == "mips":
	case l.arch == "mips64":
	case l.arch == "mips64le":
	case l.arch == "mipsle":
	case l.arch == "ppc64":
	case l.arch == "ppc64le":
	case l.arch == "riscv64":
	case l.arch == "s390x":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url += l.os + "_" + l.arch + ".zip"
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kpkgTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("spachava753", "kpkg"),
	}
}
