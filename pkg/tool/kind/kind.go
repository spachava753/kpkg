package kind

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

type kindTool struct {
	arch,
	os string
}

func (l kindTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l kindTool) Name() string {
	return "kind"
}

func (l kindTool) ShortDesc() string {
	return "Kubernetes IN Docker - local clusters for testing Kubernetes"
}

func (l kindTool) LongDesc() string {
	return `kind is a tool for running local Kubernetes clusters using Docker container "nodes". kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.`
}

func (l kindTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "ppc64le":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/kubernetes-sigs/kind/releases/download/v%s/kind-%s-%s", version, l.os, l.arch)
	return url, nil
}

func (l kindTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kind", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kind", &github.ListOptions{
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

	// sort results
	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kindTool{
		arch: arch,
		os:   os,
	}
}
