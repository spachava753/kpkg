package natscli

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type natscliTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l natscliTool) Extract(artifactPath, version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	binaryPath := filepath.Join(artifactPath, l.Name()+"-"+v.String()+"-"+l.os+"-"+l.arch)
	binaryPath = filepath.Join(binaryPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binaryPath))
	}

	return binaryPath, err
}

func (l natscliTool) Name() string {
	return "nats"
}

func (l natscliTool) ShortDesc() string {
	return "The NATS Command Line Interface"
}

func (l natscliTool) LongDesc() string {
	return "A command line utility to interact with and manage NATS."
}

func (l natscliTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/nats-io/natscli/releases/download/0.0.24/nats-0.0.24-darwin-amd64.zip
	url := fmt.Sprintf("%s%s/nats-%s-", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += l.os + "-" + l.arch
	case l.os == "windows" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386":
		url += l.os + "-" + l.arch
	case l.os == "freebsd" && l.arch == "amd64":
		url += l.os + "-" + l.arch
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "386",
		l.os == "linux" && l.arch == "arm64":
		url += l.os + "-" + l.arch
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".zip", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return natscliTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("nats-io", "natscli", 20),
	}
}
