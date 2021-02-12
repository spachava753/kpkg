package opa

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

type opaTool struct {
	arch,
	os string
}

func (l opaTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l opaTool) Name() string {
	return "opa"
}

func (l opaTool) ShortDesc() string {
	return "An open source, general-purpose policy engine"
}

func (l opaTool) LongDesc() string {
	return `The Open Policy Agent (OPA) is an open source, general-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack`
}

func (l opaTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	// only amd64 supported
	if l.arch != "amd64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	switch {
	case l.os == "darwin":
		fallthrough
	case l.os == "windows":
		fallthrough
	case l.os == "linux":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/open-policy-agent/opa/releases/download/v%s/opa_%s_amd64", version, l.os)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func (l opaTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "open-policy-agent", "opa", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "open-policy-agent", "opa", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 15 - len(releases),
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

	// dont need to many releases
	vs = vs[:15]

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return opaTool{
		arch: arch,
		os:   os,
	}
}
