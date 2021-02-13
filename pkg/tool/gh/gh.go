package gh

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ghTool struct {
	arch,
	os string
}

func (l ghTool) Extract(artifactPath, version string) (string, error) {
	// gh releases contain the binary and some examples. Pick only the binary

	// expect given path to be to contain bin folder
	dirsInfo, err := ioutil.ReadDir(artifactPath)
	if err != nil {
		return "", err
	}
	// expect to find exactly on dir
	if len(dirsInfo) != 1 {
		return "", fmt.Errorf("found too many subdirectories at %s", artifactPath)
	}
	binDirPath := filepath.Join(artifactPath, dirsInfo[0].Name(), "bin")
	binDirPathInfo, err := os.Stat(binDirPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if !binDirPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}
	binaryPath := filepath.Join(binDirPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}

	return binaryPath, err
}

func (l ghTool) Name() string {
	return "gh"
}

func (l ghTool) ShortDesc() string {
	return "GitHub’s official command line tool"
}

func (l ghTool) LongDesc() string {
	return `gh is GitHub on the command line. It brings pull requests, issues, 
and other GitHub concepts to the terminal next to where you are already working with git and your code.`
}

func (l ghTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	c, err := semver.NewConstraint("< 0.5.6")
	if err != nil {
		return "", err
	}
	// linux/arm64 is not supported until 0.5.6
	if c.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		l.os = "macOS"
		break
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
	url := fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_%s_%s", version, version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func (l ghTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "cli", "cli", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "cli", "cli", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 20 - len(releases),
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(releases, func(release *github.RepositoryRelease) bool {
		return !release.GetPrerelease() && !strings.Contains(release.GetTagName(), "rc")
	}).([]*github.RepositoryRelease)

	vs := make([]*semver.Version, len(releases))
	for i, release := range releases {
		v, err := semver.NewVersion(release.GetTagName())
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %w", err)
		}

		vs[i] = v
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return ghTool{
		arch: arch,
		os:   os,
	}
}
