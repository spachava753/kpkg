package kubectl

import (
	"fmt"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type kubectlTool struct {
	basePath,
	arch,
	os string
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
	url := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/%s/%s/kubectl", version, l.os, l.arch)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func (l kubectlTool) Versions() ([]string, error) {
	return []string{"v1.20.2", "v1.20.1", "v1.20.0", "v1.19.0", "v1.18.0", "v1.17.0", "v1.16.0", "v1.15.0",
		"v1.14.0", "v1.13.0", "v1.12.0", "v1.11.0", "v1.10.0", "v1.9.0", "v1.8.0", "v1.7.0", "v1.6.0", "v1.5.0"}, nil
}

func MakeBinary(basePath, os, arch string) tool.Binary {
	return kubectlTool{
		basePath: basePath,
		arch:     arch,
		os:       os,
	}
}
