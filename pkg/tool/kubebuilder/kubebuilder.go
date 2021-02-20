package kubebuilder

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"os"
	"path/filepath"
	"strings"
)

type kubeBuilderTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l kubeBuilderTool) Extract(artifactPath, version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	dirPath := filepath.Join(artifactPath, fmt.Sprintf("kubebuilder_%s_%s_%s", version, l.os, l.arch))
	dirPathInfo, err := os.Stat(dirPath)
	if err != nil {
		return "", err
	}
	if !dirPathInfo.IsDir() {
		return "", fmt.Errorf("expected path %s to be a dir", dirPath)
	}

	binDirPath := filepath.Join(dirPath, "bin")
	binDirPathInfo, err := os.Stat(binDirPath)
	if err != nil {
		return "", err
	}
	if !binDirPathInfo.IsDir() {
		return "", fmt.Errorf("expected path %s to be a dir", binDirPath)
	}

	binPath := filepath.Join(binDirPath, l.Name())
	binPathInfo, err := os.Stat(binPath)
	if err != nil {
		return "", err
	}
	if binPathInfo.IsDir() {
		return "", fmt.Errorf("expected path %s to be a file", binPath)
	}

	return binPath, nil
}

func (l kubeBuilderTool) Name() string {
	return "kubebuilder"
}

func (l kubeBuilderTool) ShortDesc() string {
	return "SDK for building Kubernetes APIs using CRDs"
}

func (l kubeBuilderTool) LongDesc() string {
	return `Similar to web development frameworks such as Ruby on Rails and SpringBoot, 
Kubebuilder increases velocity and reduces the complexity managed by developers for rapidly building 
and publishing Kubernetes APIs in Go. It builds on top of the canonical techniques used to build 
the core Kubernetes APIs to provide simple abstractions that reduce boilerplate and toil.`
}

func (l kubeBuilderTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "ppc64le":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/kubebuilder_%s_%s_%s.tar.gz", l.MakeReleaseUrl(), version, version, l.os, l.arch)
	return url, nil
}

func (l kubeBuilderTool) Versions() ([]string, error) {
	versions, err := l.GithubReleaseTool.Versions()
	if err != nil {
		return nil, err
	}

	versions = funk.Filter(versions, func(v string) bool {
		return !strings.Contains(v, "alpha")
	}).([]string)

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubeBuilderTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("kubernetes-sigs", "kubebuilder", 20),
	}
}
