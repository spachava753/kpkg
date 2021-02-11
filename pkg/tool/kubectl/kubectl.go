package kubectl

import (
	"fmt"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

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
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl", version, l.os, l.arch)
	if l.os == "windows" {
		url = url + ".exe"
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
