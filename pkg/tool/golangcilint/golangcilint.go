package golangcilint

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type golangciLintTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l golangciLintTool) Extract(artifactPath, version string) (
	string, error,
) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	expectedPath := filepath.Join(
		artifactPath,
		fmt.Sprintf("golangci-lint-%s-%s-%s", version, l.os, l.arch), l.Name(),
	)
	info, err := os.Stat(version)
	if err != nil {
		return "", fmt.Errorf(
			"expected path %s to contain binary", expectedPath,
		)
	}
	if info.IsDir() {
		return "", fmt.Errorf(
			"expected path %s to contain binary, found directory", expectedPath,
		)
	}
	return expectedPath, nil
}

func (l golangciLintTool) Name() string {
	return "golangci-lint"
}

func (l golangciLintTool) ShortDesc() string {
	return "Fast linters Runner for Go"
}

func (l golangciLintTool) LongDesc() string {
	return `golangci-lint is a fast Go linters runner. It runs linters in parallel, uses caching, supports yaml config, 
has integrations with all major IDE and has dozens of linters included`
}

func (l golangciLintTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// add support armv6 and armv7
	url := fmt.Sprintf(
		"%sv%s/golangci-lint-%s-%s-%s", l.MakeReleaseUrl(), version, version,
		l.os, l.arch,
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "arm64",
		l.os == "windows" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "386",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "linux" && l.arch == "386",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "mips64",
		l.os == "linux" && l.arch == "mips64le",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "s390x":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return golangciLintTool{
		arch: arch,
		os:   os,
		GithubReleaseTool: tool.MakeGithubReleaseTool(
			"golangci", "golangci-lint",
		),
	}
}
