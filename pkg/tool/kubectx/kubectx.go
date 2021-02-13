package kubectx

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/spachava753/kpkg/pkg/tool"
)

type k3supTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k3supTool) Name() string {
	return "kubectx"
}

func (l k3supTool) ShortDesc() string {
	return "Faster way to switch between clusters in kubectl"
}

func (l k3supTool) LongDesc() string {
	return `kubectx is a utility to manage and switch between kubectl(1) contexts`
}

func (l k3supTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// NOTE: there was a re-write in go from bash after version v0.9.0, but is unstable.
	// Once stable in the future, add support for downloading go binaries

	url := fmt.Sprintf("https://raw.githubusercontent.com/ahmetb/kubectx/v%s/kubectx", version)

	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3supTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("ahmetb", "kubectx", 20),
	}
}
