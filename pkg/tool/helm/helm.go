package helm

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

type helmTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l helmTool) Extract(artifactPath, _ string) (string, error) {
	// helm releases contain the binary, LICENSE and a README. Pick only the binary
	var binaryPath string
	err := filepath.Walk(
		artifactPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(filepath.Base(path), "helm") &&
				info != nil &&
				!info.IsDir() &&
				binaryPath == "" {
				binaryPath, err = filepath.Abs(path)
				return err
			}
			return nil
		},
	)

	return binaryPath, err
}

func (l helmTool) Name() string {
	return "helm"
}

func (l helmTool) ShortDesc() string {
	return "The Kubernetes Package Manager"
}

func (l helmTool) LongDesc() string {
	return "Helm is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources"
}

func (l helmTool) MakeUrl(version string) (string, error) {
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "i386",
		l.os == "linux" && l.arch == "s390x",
		l.os == "linux" && l.arch == "ppc64le":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	return fmt.Sprintf(
		"https://get.helm.sh/helm-v%s-%s-%s.tar.gz", version, l.os, l.arch,
	), nil
}

func (l helmTool) Versions(max uint) ([]string, error) {
	versions, err := l.GithubReleaseTool.Versions(max)
	if err != nil {
		return nil, err
	}

	// helm supports more architectures after this release
	c, err := semver.NewConstraint(">= 2.10.0")
	if err != nil {
		return nil, err
	}

	versions = funk.Filter(
		versions, func(v string) bool {
			sv := semver.MustParse(v)
			return c.Check(sv)
		},
	).([]string)

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return helmTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("helm", "helm"),
	}
}
