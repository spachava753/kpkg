package k3sup

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

type k3supTool struct {
	arch,
	os string
}

func (l k3supTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l k3supTool) Name() string {
	return "k3sup"
}

func (l k3supTool) ShortDesc() string {
	return "bootstrap Kubernetes with k3s over SSH < 1 min 🚀"
}

func (l k3supTool) LongDesc() string {
	return `k3sup is a light-weight utility to get from zero to KUBECONFIG with k3s on any local or remote VM. 
All you need is ssh access and the k3sup binary to get kubectl access immediately`
}

func (l k3supTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/alexellis/k3sup/releases/download/%s/k3sup", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url = url + "darwin"
		break
	case l.os == "windows" && l.arch == "amd64":
		url = url + ".exe"
		break
	case l.os == "linux" && l.arch == "amd64":
		// do nothing
		break
	case l.os == "linux" && l.arch == "arm64":
		url = url + "arm64"
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l k3supTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "alexellis", "k3sup", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "alexellis", "k3sup", &github.ListOptions{
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
	return k3supTool{
		arch: arch,
		os:   os,
	}
}