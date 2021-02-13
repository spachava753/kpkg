package popeye

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
	"strings"
)

type k3supTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k3supTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l k3supTool) Name() string {
	return "popeye"
}

func (l k3supTool) ShortDesc() string {
	return "👀 A Kubernetes cluster resource sanitizer"
}

func (l k3supTool) LongDesc() string {
	return `Popeye is a utility that scans live Kubernetes cluster and reports potential issues with deployed resources and configurations`
}

func (l k3supTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/popeye_%s_", l.MakeReleaseUrl(), version, strings.Title(l.os))
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64":
		url += "x86_64"
		break
	case l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm":
		url += l.arch
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3supTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("derailed", "popeye", 20),
	}
}
