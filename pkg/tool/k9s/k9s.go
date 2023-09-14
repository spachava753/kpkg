package k9s

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
	"strings"
)

type k9sTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l k9sTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l k9sTool) Name() string {
	return "k9s"
}

func (l k9sTool) ShortDesc() string {
	return "🐶 Kubernetes CLI To Manage Your Clusters In Style!"
}

func (l k9sTool) LongDesc() string {
	return `K9s provides a terminal UI to interact with your Kubernetes clusters. 
The aim of this project is to make it easier to navigate, observe and manage your applications in the wild. 
K9s continually watches Kubernetes for changes and offers subsequent commands to interact with your observed resources`
}

func (l k9sTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf(
		"%sv%s/k9s_%s_", l.MakeReleaseUrl(), version, strings.Title(l.os),
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		url += "x86_64"
	case l.os == "linux" && l.arch == "arm64":
		url += "arm64"
	case l.os == "linux" && l.arch == "arm":
		url += "arm"
	case l.os == "linux" && l.arch == "ppc64le":
		url += "ppc64le"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k9sTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("derailed", "k9s"),
	}
}
