package kail

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type kailTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kailTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l kailTool) Name() string {
	return "kail"
}

func (l kailTool) ShortDesc() string {
	return "kubernetes log viewer"
}

func (l kailTool) LongDesc() string {
	return `Kubernetes tail. Streams logs from all containers of all matched pods. 
Match pods by service, replicaset, deployment, and others. 
Adjusts to a changing cluster - pods are added and removed from 
logging as they fall in or out of the selection.`
}

func (l kailTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return fmt.Sprintf("%sv%s/kail_%s_%s_%s.tar.gz", l.MakeReleaseUrl(), version, version, l.os, l.arch), nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kailTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("boz", "kail", 20),
	}
}
