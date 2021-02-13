package doctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type buildxTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l buildxTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l buildxTool) Name() string {
	return "doctl"
}

func (l buildxTool) ShortDesc() string {
	return "The official command line interface for the DigitalOcean API"
}

func (l buildxTool) LongDesc() string {
	return `The official command line interface for the DigitalOcean API`
}

func (l buildxTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	c, err := semver.NewConstraint("< 1.54.1")
	if err != nil {
		return "", err
	}
	// linux/arm64 is not supported until 1.54.1
	if c.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "386":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "386":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/doctl-%s-%s-%s", l.MakeReleaseUrl(), version, version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return buildxTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("digitalocean", "doctl", 20),
	}
}
