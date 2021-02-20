package buildx

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type buildxTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l buildxTool) Name() string {
	return "buildx"
}

func (l buildxTool) ShortDesc() string {
	return "Docker CLI plugin for extended build capabilities with BuildKit"
}

func (l buildxTool) LongDesc() string {
	return `buildx is a Docker CLI plugin for extended build capabilities with BuildKit.`
}

func (l buildxTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	c, err := semver.NewConstraint("< 0.5.0")
	if err != nil {
		return "", err
	}
	// darwin/arm64 is not supported until 0.5.0
	if c.Check(v) && l.os == "darwin" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "arm64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "s390x":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/buildx-v%s.%s-%s", l.MakeReleaseUrl(), version, version, l.os, l.arch)
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return buildxTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("docker", "buildx", 20),
	}
}
