package tkn

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type tknTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l tknTool) Extract(artifactPath, _ string) (string, error) {
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

func (l tknTool) Name() string {
	return "tkn"
}

func (l tknTool) ShortDesc() string {
	return "A CLI for interacting with Tekton!"
}

func (l tknTool) LongDesc() string {
	return "The Tekton Pipelines CLI project provides a command-line interface (CLI) for interacting with Tekton, an open-source framework for Continuous Integration and Delivery (CI/CD) systems"
}

func (l tknTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/tektoncd/cli/releases/download/v0.19.1/tkn_0.19.1_Darwin_x86_64.tar.gz
	url := fmt.Sprintf("%sv%s/tkn_%s_", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "Darwin" + "_" + "x86_64"
	case l.os == "windows" && l.arch == "amd64":
		url += "Windows" + "_" + "x86_64"
	case l.os == "linux" && l.arch == "amd64":
		url += "Linux" + "_" + "x86_64"
	case l.os == "linux" && l.arch == "ppc64le":
		url += "Linux" + "_" + "ppc64le"
	case l.os == "linux" && l.arch == "s390x":
		url += "Linux" + "_" + "s390x"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	if l.os == "windows" {
		url += ".zip"
	} else {
		url += ".tar.gz"
	}

	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return tknTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("tektoncd", "cli", 20),
	}
}
