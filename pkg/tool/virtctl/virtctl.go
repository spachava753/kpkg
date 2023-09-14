package virtctl

import (
	"fmt"

	"github.com/Masterminds/semver"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type virtctlTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l virtctlTool) Name() string {
	return "virtctl"
}

func (l virtctlTool) ShortDesc() string {
	return "Kubernetes Virtualization API and runtime in order to define and manage virtual machines"
}

func (l virtctlTool) LongDesc() string {
	return `KubeVirt is a virtual machine management add-on for Kubernetes. The aim is to provide a common ground for 
virtualization solutions on top of Kubernetes`
}

func (l virtctlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf(
		"%sv%s/virtctl-v%s-%s-%s", l.MakeReleaseUrl(), version, version, l.os,
		l.arch,
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return virtctlTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("kubevirt", "kubevirt"),
	}
}
