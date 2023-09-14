package goreleaser

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"path/filepath"
	"strings"
)

type goreleaserTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l goreleaserTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l goreleaserTool) Name() string {
	return "goreleaser"
}

func (l goreleaserTool) ShortDesc() string {
	return "Deliver Go binaries as fast and easily as possible"
}

func (l goreleaserTool) LongDesc() string {
	return `GoReleaser builds Go binaries for several platforms, 
creates a GitHub release and then pushes a Homebrew formula to a tap repository. 
All that wrapped in your favorite CI`
}

func (l goreleaserTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf(
		"%sv%s/goreleaser_%s_", l.MakeReleaseUrl(), version, strings.Title(l.os),
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		url += "x86_64"
	case l.os == "linux" && l.arch == "arm64":
		url += "arm64"
	case l.os == "windows" && l.arch == "386",
		l.os == "linux" && l.arch == "386":
		url += "i386"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return goreleaserTool{
		arch: arch,
		os:   os,
		GithubReleaseTool: tool.MakeGithubReleaseTool(
			"goreleaser", "goreleaser",
		),
	}
}
