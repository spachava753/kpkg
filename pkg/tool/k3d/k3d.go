package k3d

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type k3dTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k3dTool) Name() string {
	return "k3d"
}

func (l k3dTool) ShortDesc() string {
	return "Little helper to run Rancher Lab's k3s in Docker"
}

func (l k3dTool) LongDesc() string {
	return "k3d creates containerized k3s clusters. This means, that you can spin up a multi-node k3s cluster on a single machine using docker"
}

func (l k3dTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "386":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf(
		"%sv%s/k3d-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch,
	)
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3dTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("rancher", "k3d"),
	}
}
