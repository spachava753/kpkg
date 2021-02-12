package helmfile

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

type helmFileTool struct {
	arch,
	os string
}

func (l helmFileTool) Extract(artifactPath, version string) (string, error) {
	return artifactPath, nil
}

func (l helmFileTool) Name() string {
	return "helmfile"
}

func (l helmFileTool) ShortDesc() string {
	return "Deploy Kubernetes Helm Charts"
}

func (l helmFileTool) LongDesc() string {
	return `Helmfile is a declarative spec for deploying helm charts`
}

func (l helmFileTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "darwin" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "386":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/roboll/helmfile/releases/download/v%s/helmfile_%s_%s", version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".exe", nil
	}
	return url, nil
}

func (l helmFileTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "roboll", "helmfile", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "roboll", "helmfile", &github.ListOptions{
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

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return helmFileTool{
		arch: arch,
		os:   os,
	}
}
