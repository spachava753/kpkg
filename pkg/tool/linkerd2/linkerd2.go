package linkerd2

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"sort"
	"strings"
)

type linkerd2Tool struct {
	arch,
	os string
}

func (l linkerd2Tool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l linkerd2Tool) Name() string {
	return "linkerd2"
}

func (l linkerd2Tool) ShortDesc() string {
	return "linkerd2 is a cli to install linkerd2 service mesh"
}

func (l linkerd2Tool) LongDesc() string {
	return "linkerd2 is a cli to install linkerd2 service mesh"
}

func (l linkerd2Tool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	// install the latest stable binary
	switch l.os {
	case "darwin":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/stable-%s/linkerd2-cli-stable-%s-darwin", version, version), nil
	case "windows":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/stable-%s/linkerd2-cli-stable-%s-windows.exe", version, version), nil
	case "linux":
		switch l.arch {
		case "amd64":
			fallthrough
		case "arm":
			fallthrough
		case "arm64":
			return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/stable-%s/linkerd2-cli-stable-%s-linux-%s", version, version, l.arch), nil
		default:
			return "", fmt.Errorf("unsupported architecture: %s", l.arch)
		}
	}
	return "", fmt.Errorf("unsupported os: %s", l.os)
}

func (l linkerd2Tool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "linkerd", "linkerd2", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "linkerd", "linkerd2", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(releases, func(release *github.RepositoryRelease) bool {
		return !release.GetPrerelease() && strings.Contains(release.GetTagName(), "stable")
	}).([]*github.RepositoryRelease)

	vs := make([]*semver.Version, len(releases))
	for i, release := range releases {
		a := release.GetTagName()[7:]
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
	return linkerd2Tool{
		arch: arch,
		os:   os,
	}
}
