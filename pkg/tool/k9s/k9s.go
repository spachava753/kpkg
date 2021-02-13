package k9s

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"path/filepath"
	"sort"
	"strings"
)

type k9sTool struct {
	arch,
	os string
}

func (l k9sTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l k9sTool) Name() string {
	return "k9s"
}

func (l k9sTool) ShortDesc() string {
	return "🐶 Kubernetes CLI To Manage Your Clusters In Style!"
}

func (l k9sTool) LongDesc() string {
	return `K9s provides a terminal UI to interact with your Kubernetes clusters. 
The aim of this project is to make it easier to navigate, observe and manage your applications in the wild. 
K9s continually watches Kubernetes for changes and offers subsequent commands to interact with your observed resources`
}

func (l k9sTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/derailed/k9s/releases/download/v%s/k9s_%s_", version, strings.Title(l.os))
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		url = url + "x86_64"
		break
	case l.os == "linux" && l.arch == "arm64":
		url = url + "arm64"
		break
	case l.os == "linux" && l.arch == "arm":
		url = url + "arm"
		break
	case l.os == "linux" && l.arch == "ppc64le":
		url = url + "ppc64le"
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func (l k9sTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "derailed", "k9s", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "derailed", "k9s", &github.ListOptions{
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

	// sort results
	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k9sTool{
		arch: arch,
		os:   os,
	}
}
