package kubens

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
	return "kubens"
}

func (l k3supTool) ShortDesc() string {
	return "Faster way to switch between namespaces in kubectl"
}

func (l k3supTool) LongDesc() string {
	return `kubens is a utility to switch between Kubernetes namespaces`
}

func (l k3supTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// NOTE: there was a re-write in go from bash after version v0.9.0, but is unstable.
	// Once stable in the future, add support for downloading go binaries

	url := fmt.Sprintf(
		"https://raw.githubusercontent.com/ahmetb/kubectx/v%s/kubens", version,
	)

	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3supTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("ahmetb", "kubectx"),
	}
}
