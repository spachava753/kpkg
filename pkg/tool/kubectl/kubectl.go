package kubectl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

// doesn't need GithubReleaseTool
type kubectlTool struct {
	arch,
	os string
}

func (l kubectlTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l kubectlTool) Name() string {
	return "kubectl"
}

func (l kubectlTool) ShortDesc() string {
	return "kubectl is a cli to communicate k8s clusters"
}

func (l kubectlTool) LongDesc() string {
	return "kubectl is a cli to communicate k8s clusters"
}

func (l kubectlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl", version, l.os, l.arch)
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func (l kubectlTool) Versions() ([]string, error) {
	return []string{"1.20.2", "1.20.1", "1.20.0", "1.19.0", "1.18.0", "1.17.0", "1.16.0", "1.15.0",
		"1.14.0", "1.13.0", "1.12.0", "1.11.0", "1.10.0", "1.9.0", "1.8.0", "1.7.0", "1.6.0", "1.5.0"}, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubectlTool{
		arch: arch,
		os:   os,
	}
}
