package minikube

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type minikubeTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l minikubeTool) Name() string {
	return "minikube"
}

func (l minikubeTool) ShortDesc() string {
	return "bootstrap Kubernetes with k3s over SSH < 1 min 🚀"
}

func (l minikubeTool) LongDesc() string {
	return `k3sup is a light-weight utility to get from zero to KUBECONFIG with k3s on any local or remote VM. 
All you need is ssh access and the k3sup binary to get kubectl access immediately`
}

func (l minikubeTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// darwin/arm64 is not supported until 1.17.1
	c1, err := semver.NewConstraint("< 1.17.1")
	if err != nil {
		return "", err
	}

	if !c1.Check(v) && l.os == "darwin" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	// other archs for linux is not supported until 1.10.0
	c2, err := semver.NewConstraint("< 1.10.0")
	if err != nil {
		return "", err
	}

	if c2.Check(v) && l.os == "linux" && l.arch != "amd64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	url := fmt.Sprintf("%sv%s/minikube-%s-%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "arm64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "s390x":
		break
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return minikubeTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("kubernetes", "minikube", 20),
	}
}
