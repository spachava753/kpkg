package argocd

import (
	"fmt"

	"github.com/Masterminds/semver"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type argocdTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l argocdTool) Name() string {
	return "argocd"
}

func (l argocdTool) ShortDesc() string {
	return "Declarative continuous deployment for Kubernetes"
}

func (l argocdTool) LongDesc() string {
	return `argocd is a cli to complement the ArgoCD project`
}

func (l argocdTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// only supports amd64
	if l.arch != "amd64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	url := fmt.Sprintf(
		"%sv%s/argocd-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch,
	)
	switch {
	case l.os == "darwin",
		l.os == "linux":
	case l.os == "windows":
		url += ".exe"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return argocdTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("argoproj", "argo-cd"),
	}
}
