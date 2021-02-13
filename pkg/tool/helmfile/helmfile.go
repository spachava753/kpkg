package helmfile

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type helmFileTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l helmFileTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l helmFileTool) Name() string {
	return "helmfile"
}

func (l helmFileTool) ShortDesc() string {
	return "Deploy Kubernetes Helm Charts"
}

func (l helmFileTool) LongDesc() string {
	return `Helmfile is a declarative spec for deploying helm charts`
}

func (l helmFileTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "386":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/helmfile_%s_%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".exe", nil
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return helmFileTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("roboll", "helmfile", 20),
	}
}
