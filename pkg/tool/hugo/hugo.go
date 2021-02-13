package hugo

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type hugoTool struct {
	arch,
	os string
}

func (l hugoTool) Extract(artifactPath, version string) (string, error) {
	// hugo releases contain the binary and some examples. Pick only the binary
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is a dir", binaryPath))
	}

	return binaryPath, err
}

func (l hugoTool) Name() string {
	return "hugo"
}

func (l hugoTool) ShortDesc() string {
	return "The world’s fastest framework for building websites"
}

func (l hugoTool) LongDesc() string {
	return `Hugo is a static HTML and CSS website generator written in Go. It is optimized for speed, ease of use, and configurability. Hugo takes a directory with content and templates and renders them into a full HTML website`
}

func (l hugoTool) MakeUrl(version string) (string, error) {
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

	url := fmt.Sprintf("https://github.com/gohugoio/hugo/releases/download/v%s/", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_macOS-64bit", version)
		break
	case l.os == "windows" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_Windows-64bit", version)
		break
	case l.os == "windows" && l.arch == "386":
		url = url + fmt.Sprintf("hugo_%s_Windows-32bit", version)
		break
	case l.os == "linux" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_Linux-64bit", version)
		break
	case l.os == "linux" && l.arch == "arm64":
		url = url + fmt.Sprintf("hugo_%s_Linux-ARM64", version)
		break
	case l.os == "linux" && l.arch == "arm":
		url = url + fmt.Sprintf("hugo_%s_Linux-ARM", version)
		break
	case l.os == "linux" && l.arch == "386":
		url = url + fmt.Sprintf("hugo_%s_Linux-32bit", version)
		break
	case l.os == "freebsd" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_FreeBSD-64bit", version)
		break
	case l.os == "freebsd" && l.arch == "arm":
		url = url + fmt.Sprintf("hugo_%s_FreeBSD-ARM", version)
		break
	case l.os == "freebsd" && l.arch == "386":
		url = url + fmt.Sprintf("hugo_%s_FreeBSD-64bit", version)
		break
	case l.os == "dragonfly" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_DragonFlyBSD-64bit", version)
		break
	case l.os == "openbsd" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_OpenBSD-64bit", version)
		break
	case l.os == "openbsd" && l.arch == "arm":
		url = url + fmt.Sprintf("hugo_%s_OpenBSD-ARM", version)
		break
	case l.os == "openbsd" && l.arch == "386":
		url = url + fmt.Sprintf("hugo_%s_OpenBSD-32bit", version)
		break
	case l.os == "netbsd" && l.arch == "amd64":
		url = url + fmt.Sprintf("hugo_%s_NetBSD-64bit", version)
		break
	case l.os == "netbsd" && l.arch == "arm":
		url = url + fmt.Sprintf("hugo_%s_NetBSD-ARM", version)
		break
	case l.os == "netbsd" && l.arch == "386":
		url = url + fmt.Sprintf("hugo_%s_NetBSD-32bit", version)
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func (l hugoTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "gohugoio", "hugo", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "gohugoio", "hugo", &github.ListOptions{
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
	return hugoTool{
		arch: arch,
		os:   os,
	}
}
