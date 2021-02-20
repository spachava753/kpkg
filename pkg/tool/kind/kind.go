package kind

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type kindTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kindTool) Name() string {
	return "kind"
}

func (l kindTool) ShortDesc() string {
	return "Kubernetes IN Docker - local clusters for testing Kubernetes"
}

func (l kindTool) LongDesc() string {
	return `kind is a tool for running local Kubernetes clusters using Docker container "nodes". kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.`
}

func (l kindTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/kind-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kindTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("kubernetes-sigs", "kind", 20),
	}
}
