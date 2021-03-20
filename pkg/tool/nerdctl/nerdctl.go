package nerdctl

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type nerdCtlTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l nerdCtlTool) Extract(artifactPath, _ string) (string, error) {
	binPath := filepath.Join(artifactPath, l.Name())
	if _, err := os.Stat(binPath); err != nil {
		return "", err
	}
	return binPath, nil
}

func (l nerdCtlTool) Name() string {
	return "nerdctl"
}

func (l nerdCtlTool) ShortDesc() string {
	return "Docker-compatible CLI for containerd"
}

func (l nerdCtlTool) LongDesc() string {
	return `nerdctl is a Docker-compatible CLI for containerd.

- Same UI/UX as docker

- Supports rootless mode

- Supports lazy-pulling (Stargz)

- Supports encrypted images (ocicrypt)

nerdctl is a non-core sub-project of containerd.`
}

func (l nerdCtlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// only supports linux
	if l.os != "linux" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	// add support for detecting armv7
	switch {
	case l.arch == "amd64",
		l.arch == "arm64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/nerdctl-%s-linux-%s.tar.gz", l.MakeReleaseUrl(), version, version, l.arch)
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return nerdCtlTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("containerd", "nerdctl", 20),
	}
}
