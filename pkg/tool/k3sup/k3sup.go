package k3sup

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type k3supTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k3supTool) Name() string {
	return "k3sup"
}

func (l k3supTool) ShortDesc() string {
	return "bootstrap Kubernetes with k3s over SSH < 1 min 🚀"
}

func (l k3supTool) LongDesc() string {
	return `k3sup is a light-weight utility to get from zero to KUBECONFIG with k3s on any local or remote VM. 
All you need is ssh access and the k3sup binary to get kubectl access immediately`
}

func (l k3supTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/k3sup", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "darwin"
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
	case l.os == "linux" && l.arch == "amd64":
	case l.os == "linux" && l.arch == "arm64":
		url += "arm64"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3supTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("alexellis", "k3sup"),
	}
}
