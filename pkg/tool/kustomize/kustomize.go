package kustomize

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

type kustomizeTool struct {
	arch,
	os string
}

func (l kustomizeTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l kustomizeTool) Name() string {
	return "kustomize"
}

func (l kustomizeTool) ShortDesc() string {
	return "Customization of kubernetes YAML configurations"
}

func (l kustomizeTool) LongDesc() string {
	return `kustomize lets you customize raw, template-free YAML files for multiple purposes, 
leaving the original YAML untouched and usable as is`
}

func (l kustomizeTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v%s/kustomize_v%s_%s_%s.tar.gz", version, version, l.os, l.arch)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l kustomizeTool) Versions() ([]string, error) {
	max := 50
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kustomize", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < max {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kustomize", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: max - len(releases),
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(releases, func(release *github.RepositoryRelease) bool {
		return !release.GetPrerelease() && strings.Contains(release.GetTagName(), "kustomize")
	}).([]*github.RepositoryRelease)

	vs := make([]*semver.Version, len(releases))
	for i, release := range releases {
		a := release.GetTagName()[10:]
		v, err := semver.NewVersion(a)
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
	return kustomizeTool{
		arch: arch,
		os:   os,
	}
}
