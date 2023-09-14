package krew

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
)

type krewTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l krewTool) Extract(artifactPath, _ string) (string, error) {
	if l.os == "windows" {
		return filepath.Join(artifactPath, l.Name()), nil
	}
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return filepath.Join(
		artifactPath, l.Name()+fmt.Sprintf("-%s_%s", l.os, l.arch),
	), nil
}

func (l krewTool) Name() string {
	return "krew"
}

func (l krewTool) ShortDesc() string {
	return "📦 Find and install kubectl plugins"
}

func (l krewTool) LongDesc() string {
	return `Krew is a tool that makes it easy to use kubectl plugins. 
Krew helps you discover plugins, install and manage them on your machine. 
It is similar to tools like apt, dnf or brew. Today, over 100 kubectl plugins are available on Krew`
}

func (l krewTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/krew", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".exe", err
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return krewTool{
		arch: arch,
		os:   os,
		GithubReleaseTool: tool.MakeGithubReleaseTool(
			"kubernetes-sigs", "krew",
		),
	}
}
