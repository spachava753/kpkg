package pack

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type packTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l packTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l packTool) Name() string {
	return "pack"
}

func (l packTool) ShortDesc() string {
	return "CLI for building apps using Cloud Native Buildpacks"
}

func (l packTool) LongDesc() string {
	return `pack makes it easy for app developers to use buildpacks to convert code into runnable images,
buildpack authors to develop and package buildpacks for distribution, and
operators to package buildpacks for distribution and maintain applications.`
}

func (l packTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	c, err := semver.NewConstraint("< 0.17.0")
	if err != nil {
		return "", err
	}

	// darwin/arm64 is not supported until 0.17.0
	if c.Check(v) && l.os == "darwin" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	url := fmt.Sprintf("%sv%s/pack-v%s-", l.MakeReleaseUrl(), version, version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "macos.tgz"
	case l.os == "darwin" && l.arch == "arm64":
		url += "macos-arm64.tgz"
	case l.os == "windows" && l.arch == "amd64":
		url += "windows.zip"
	case l.os == "linux" && l.arch == "amd64":
		url += "linux.tgz"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return packTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("buildpacks", "pack", 20),
	}
}
