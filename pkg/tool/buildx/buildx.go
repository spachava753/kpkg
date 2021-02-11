package buildx

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

type buildxTool struct {
	arch,
	os string
}

func (l buildxTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l buildxTool) Name() string {
	return "buildx"
}

func (l buildxTool) ShortDesc() string {
	return "Docker CLI plugin for extended build capabilities with BuildKit"
}

func (l buildxTool) LongDesc() string {
	return `buildx is a Docker CLI plugin for extended build capabilities with BuildKit.`
}

func (l buildxTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	c, err := semver.NewConstraint("< 0.5.0")
	if err != nil {
		return "", err
	}
	// darwin/arm64 is not supported until 0.5.0
	if c.Check(v) && l.os == "darwin" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "darwin" && l.arch == "arm64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "ppc64le":
		fallthrough
	case l.os == "linux" && l.arch == "s390x":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/docker/buildx/releases/download/v%s/buildx-v%s.%s-%s", version, version, l.os, l.arch)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func (l buildxTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "docker", "buildx", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "docker", "buildx", &github.ListOptions{
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

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return buildxTool{
		arch: arch,
		os:   os,
	}
}
