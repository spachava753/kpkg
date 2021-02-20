package inletsctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type inletsctlTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l inletsctlTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l inletsctlTool) Name() string {
	return "inletsctl"
}

func (l inletsctlTool) ShortDesc() string {
	return "The fastest way to create self-hosted exit-servers"
}

func (l inletsctlTool) LongDesc() string {
	return `inletsctl automates the task of creating an exit-node on cloud infrastructure. 
Once provisioned, you'll receive a command to connect with. 
You can use this tool whether you want to use inlets or inlets-pro for L4 TCP`
}

func (l inletsctlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/inletsctl", l.MakeReleaseUrl(), version)
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

	return url + ".tgz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return inletsctlTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("inlets", "inletsctl", 20),
	}
}
