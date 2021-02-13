package kops

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type kopsTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kopsTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l kopsTool) Name() string {
	return "kops"
}

func (l kopsTool) ShortDesc() string {
	return "Kubernetes Operations (kops) - Production Grade K8s Installation, Upgrades, and Management"
}

func (l kopsTool) LongDesc() string {
	return `kops will not only help you create, destroy, upgrade and maintain production-grade, 
highly available, Kubernetes cluster, but it will also provision the necessary cloud infrastructure`
}

func (l kopsTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/kops-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kopsTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("kubernetes", "kops", 20),
	}
}
