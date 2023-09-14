package fzf

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type fzfTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l fzfTool) Extract(artifactPath, _ string) (string, error) {
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

func (l fzfTool) Name() string {
	return "fzf"
}

func (l fzfTool) ShortDesc() string {
	return "🌸 A command-line fuzzy finder"
}

func (l fzfTool) LongDesc() string {
	return "fzf is a general-purpose command-line fuzzy finder"
}

func (l fzfTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	darwinArm64Constraint, err := semver.NewConstraint("> 0.25.1")
	if err != nil {
		return "", err
	}

	githubReleaseConstraint, err := semver.NewConstraint("< 0.24.0")
	if err != nil {
		return "", err
	}

	if githubReleaseConstraint.Check(v) {
		return "", fmt.Errorf("fzf does not have github releases before 0.24.0")
	}

	url := fmt.Sprintf("%s%s/fzf-%s-", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += l.os + "_" + l.arch + ".zip"
	case l.os == "windows" && l.arch == "amd64":
		url += l.os + "_" + l.arch + ".zip"
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "openbsd" && l.arch == "amd64":
		url += l.os + "_" + l.arch + ".tar.gz"
	default:
		if darwinArm64Constraint.Check(v) &&
			l.os == "darwin" && l.arch == "arm64" {
			url += l.os + "_" + l.arch + ".zip"
		}
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return fzfTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("junegunn", "fzf"),
	}
}
