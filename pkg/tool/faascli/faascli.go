package faascli

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type faasCliTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l faasCliTool) Name() string {
	return "faas-cli"
}

func (l faasCliTool) ShortDesc() string {
	return "openfaas CLI plugin for extended build capabilities with BuildKit"
}

func (l faasCliTool) LongDesc() string {
	return `faas-cli is a openfaas CLI plugin for extended build capabilities with BuildKit.`
}

func (l faasCliTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/faas-cli", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "-darwin"
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
	case l.os == "linux" && l.arch == "arm":
		url += "-armhf"
	case l.os == "linux" && l.arch == "amd64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return faasCliTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("openfaas", "faas-cli", 20),
	}
}
