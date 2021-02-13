package stern

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type sternTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l sternTool) Name() string {
	return "stern"
}

func (l sternTool) ShortDesc() string {
	return "⎈ Multi pod and container log tailing for Kubernetes"
}

func (l sternTool) LongDesc() string {
	return `Stern allows you to tail multiple pods on Kubernetes and multiple containers within the pod. 
Each result is color coded for quicker debugging`
}

func (l sternTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/stern_%s_%s", l.MakeReleaseUrl(), version, l.os, l.arch)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".exe", nil
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return sternTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("wercker", "stern", 20),
	}
}
