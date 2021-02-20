package kubebench

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type kubeBenchTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kubeBenchTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l kubeBenchTool) Name() string {
	return "kube-bench"
}

func (l kubeBenchTool) ShortDesc() string {
	return "Checks whether Kubernetes is deployed according to security best practices as defined in the CIS Kubernetes Benchmark"
}

func (l kubeBenchTool) LongDesc() string {
	return `kube-bench is a Go application that checks whether Kubernetes is deployed securely by running 
the checks documented in the CIS Kubernetes Benchmark`
}

func (l kubeBenchTool) MakeUrl(version string) (string, error) {
	if l.os != "linux" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/kube-bench_%s_%s_%s.tar.gz", l.MakeReleaseUrl(), version, version, l.os, l.arch)

	switch {
	case l.arch == "amd64", l.arch == "arm64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubeBenchTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("aquasecurity", "kube-bench", 20),
	}
}
