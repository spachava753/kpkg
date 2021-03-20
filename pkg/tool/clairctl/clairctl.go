package clairctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type clairTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l clairTool) Name() string {
	return "clairctl"
}

func (l clairTool) ShortDesc() string {
	return "Vulnerability Static Analysis for Containers"
}

func (l clairTool) LongDesc() string {
	return `Clair is an open source project for the static analysis of vulnerabilities in application containers 
(currently including OCI and docker).`
}

func (l clairTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/clairctl-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "386",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return clairTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("quay", "clair", 20),
	}
}
