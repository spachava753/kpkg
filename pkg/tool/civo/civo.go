package civo

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

type civoTool struct {
	arch,
	os string
}

func (l civoTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l civoTool) Name() string {
	return "civo"
}

func (l civoTool) ShortDesc() string {
	return "Civo CLI is a tool to manage your Civo.com account from the terminal"
}

func (l civoTool) LongDesc() string {
	return `Civo CLI is a tool to manage your Civo.com account from the terminal. 
The Civo web control panel has a user-friendly interface for managing your account, 
but in case you want to automate or run scripts on your account, 
or have multiple complex services, the command-line interface outlined here will be useful. `
}

func (l civoTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	c1, err := semver.NewConstraint("<= 0.6.36")
	if err != nil {
		return "", err
	}

	if c1.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c2, err := semver.NewConstraint("<= 0.6.30")
	if err != nil {
		return "", err
	}

	if c2.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c3, err := semver.NewConstraint(">= 0.6.38")
	if err != nil {
		return "", err
	}

	if c3.Check(v) && l.os == "windows" && l.arch == "386" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c4, err := semver.NewConstraint("> 0.6.18")
	if err != nil {
		return "", err
	}

	if c4.Check(v) {
		switch {
		case l.os == "darwin":
			fallthrough
		case l.os == "linux":
			fallthrough
		case l.os == "windows":
			break
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "386":
		fallthrough
	case l.os == "freebsd" && l.arch == "amd64":
		fallthrough
	case l.os == "freebsd" && l.arch == "386":
		fallthrough
	case l.os == "freebsd" && l.arch == "arm":
		fallthrough
	case l.os == "openbsd" && l.arch == "386":
		fallthrough
	case l.os == "openbsd" && l.arch == "amd64":
		fallthrough
	case l.os == "solaris" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		fallthrough
	case l.os == "linux" && l.arch == "386":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/civo/cli/releases/download/v%s/civo-%s-%s-%s", version, version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func (l civoTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "civo", "cli", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "civo", "cli", &github.ListOptions{
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

	// civo supports more architectures after this release
	c, err := semver.NewConstraint(">= 0.6.0")
	if err != nil {
		return nil, err
	}

	vs = funk.Filter(vs, func(v *semver.Version) bool {
		return c.Check(v)
	}).([]*semver.Version)

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return civoTool{
		arch: arch,
		os:   os,
	}
}
