package k3d

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

type k3dTool struct {
	arch,
	os string
}

func (l k3dTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l k3dTool) Name() string {
	return "k3d"
}

func (l k3dTool) ShortDesc() string {
	return "Little helper to run Rancher Lab's k3s in Docker"
}

func (l k3dTool) LongDesc() string {
	return "k3d creates containerized k3s clusters. This means, that you can spin up a multi-node k3s cluster on a single machine using docker"
}

func (l k3dTool) MakeUrl(version string) (string, error) {
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
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
	url := fmt.Sprintf("https://github.com/rancher/k3d/releases/download/v%s/k3d-%s-%s", version, l.os, l.arch)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func (l k3dTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "rancher", "k3d", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "rancher", "k3d", &github.ListOptions{
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

	// k3d supports more architectures after this release
	c, err := semver.NewConstraint(">= 3.0.0")
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
	return k3dTool{
		arch: arch,
		os:   os,
	}
}
