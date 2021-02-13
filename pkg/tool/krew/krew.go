package krew

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

type krewTool struct {
	arch,
	os string
}

func (l krewTool) Extract(artifactPath, _ string) (string, error) {
	if l.os == "windows" {
		return filepath.Join(artifactPath, l.Name()), nil
	}
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return filepath.Join(artifactPath, l.Name()+fmt.Sprintf("-%s_%s", l.os, l.arch)), nil
}

func (l krewTool) Name() string {
	return "krew"
}

func (l krewTool) ShortDesc() string {
	return "📦 Find and install kubectl plugins"
}

func (l krewTool) LongDesc() string {
	return `Krew is a tool that makes it easy to use kubectl plugins. 
Krew helps you discover plugins, install and manage them on your machine. 
It is similar to tools like apt, dnf or brew. Today, over 100 kubectl plugins are available on Krew`
}

func (l krewTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/kubernetes-sigs/krew/releases/download/v%s/krew", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".exe", err
	}
	return url + ".tar.gz", nil
}

func (l krewTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "krew", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "krew", &github.ListOptions{
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
	return krewTool{
		arch: arch,
		os:   os,
	}
}
