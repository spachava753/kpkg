package yq

import (
	"fmt"

	"github.com/Masterminds/semver"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type yqTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l yqTool) Name() string {
	return "yq"
}

func (l yqTool) ShortDesc() string {
	return "yq is a portable command-line YAML processor"
}

func (l yqTool) LongDesc() string {
	return `a lightweight and portable command-line YAML processor. 
yq uses jq like syntax but works with yaml files as well as json. 
It doesn't yet support everything jq does - but it does support the most common operations and functions, 
and more is being added continuously.`
}

func (l yqTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf(
		"%sv%s/yq_%s_%s", l.MakeReleaseUrl(), version, l.os, l.arch,
	)

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "386",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "arm",
		l.os == "linux" && l.arch == "386",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "mips",
		l.os == "linux" && l.arch == "mips64",
		l.os == "linux" && l.arch == "mips64le",
		l.os == "linux" && l.arch == "mipsle",
		l.os == "linux" && l.arch == "ppc64",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "s390x",
		l.os == "netbsd" && l.arch == "386",
		l.os == "netbsd" && l.arch == "amd64",
		l.os == "netbsd" && l.arch == "arm",
		l.os == "openbsd" && l.arch == "386",
		l.os == "openbsd" && l.arch == "amd64":
	case l.os == "windows" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64":
		url += ".exe"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return yqTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("mikefarah", "yq"),
	}
}
