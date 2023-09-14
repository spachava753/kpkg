package kubeprompt

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type kubePromptTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kubePromptTool) Extract(artifactPath, _ string) (string, error) {
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

func (l kubePromptTool) Name() string {
	return "kube-prompt"
}

func (l kubePromptTool) ShortDesc() string {
	return "An interactive kubernetes client featuring auto-complete"
}

func (l kubePromptTool) LongDesc() string {
	return "An interactive kubernetes client featuring auto-complete using go-prompt"
}

func (l kubePromptTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/c-bata/kube-prompt/releases/download/v1.0.11/kube-prompt_v1.0.11_linux_amd64.zip
	url := fmt.Sprintf(
		"%sv%s/kube-prompt_v%s_", l.MakeReleaseUrl(), version, version,
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "386":
		url += l.os + "_" + l.arch + ".zip"
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "386",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "arm64":
		url += l.os + "_" + l.arch + ".zip"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubePromptTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("c-bata", "kube-prompt"),
	}
}
