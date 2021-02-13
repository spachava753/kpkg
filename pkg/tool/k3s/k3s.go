package k3s

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type k3sTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k3sTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l k3sTool) Name() string {
	return "k3s"
}

func (l k3sTool) ShortDesc() string {
	return "Lightweight Kubernetes"
}

func (l k3sTool) LongDesc() string {
	return "Lightweight Kubernetes. Production ready, easy to install, half the memory, all in a binary less than 100 MB"
}

func (l k3sTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	if l.os != "linux" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/k3s", l.MakeReleaseUrl(), version)
	switch {
	case l.arch == "amd64":
		break
	case l.arch == "arm64":
		url = url + "-arm64"
	case l.arch == "arm":
		url = url + "-armhf"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3sTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("k3s-io", "k3s", 20),
	}
}
