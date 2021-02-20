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
	return "Run Kubernetes locally"
}

func (l minikubeTool) LongDesc() string {
	return `minikube implements a local Kubernetes cluster on macOS, Linux, and Windows. 
minikube's primary goals are to be the best tool for local Kubernetes application development and 
to support all Kubernetes features that fit`
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
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
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
