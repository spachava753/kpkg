package k3s

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"sort"
	"strings"
)

type k3sTool struct {
	arch,
	os string
}

func (l k3sTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l k3sTool) Name() string {
	return "k3s"
}

func (l k3sTool) ShortDesc() string {
	return "Lightweight Kubernetes"
}

func (l k3sTool) LongDesc() string {
	return "Lightweight Kubernetes. Production ready, easy to install, half the memory, all in a binary less than 100 MB"
}

func (l k3sTool) MakeUrl(version string) (string, error) {
	var url string
	switch {
	case l.arch == "amd64":
		url = fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/v%s/k3s", version)
	case l.arch == "arm64":
		url = fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/v%s/k3s-arm64", version)
	case l.arch == "arm":
		url = fmt.Sprintf("https://github.com/k3s-io/k3s/releases/download/v%s/k3s-armhf", version)
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l k3sTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "k3s-io", "k3s", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "k3s-io", "k3s", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 100,
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
	return k3sTool{
		arch: arch,
		os:   os,
	}
}
