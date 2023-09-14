package kubebuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/thoas/go-funk"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
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

	downloadConstraint, err := semver.NewConstraint(">= 3.0.0")
	if err != nil {
		return "", err
	}
	if downloadConstraint.Check(v) {
		return artifactPath, nil
	}

	version = v.String()
	dirPath := filepath.Join(
		artifactPath,
		fmt.Sprintf("kubebuilder_%s_%s_%s", version, l.os, l.arch),
	)
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
	downloadConstraint, err := semver.NewConstraint("< 3.0.0")
	if err != nil {
		return "", err
	}

	archConstraint, err := semver.NewConstraint("<= 2.0.0")
	if err != nil {
		return "", err
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
	default:
		if archConstraint.Check(v) &&
			l.os == "linux" &&
			(l.arch == "arm64" || l.arch == "ppc64le") {
			break
		}
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	if downloadConstraint.Check(v) {
		url := fmt.Sprintf(
			"%sv%s/kubebuilder_%s_%s_%s.tar.gz", l.MakeReleaseUrl(), version,
			version, l.os, l.arch,
		)
		return url, nil
	}

	// https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.1.0/kubebuilder_darwin_amd64
	url := fmt.Sprintf(
		"%sv%s/kubebuilder_%s_%s", l.MakeReleaseUrl(), version, l.os, l.arch,
	)
	return url, nil
}

func (l kubeBuilderTool) Versions(max uint) ([]string, error) {
	versions, err := l.GithubReleaseTool.Versions(max)
	if err != nil {
		return nil, err
	}

	versions = funk.Filter(
		versions, func(v string) bool {
			return !strings.Contains(v, "alpha")
		},
	).([]string)

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubeBuilderTool{
		arch: arch,
		os:   os,
		GithubReleaseTool: tool.MakeGithubReleaseTool(
			"kubernetes-sigs", "kubebuilder",
		),
	}
}
