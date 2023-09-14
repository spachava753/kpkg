package flux

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	"github.com/thoas/go-funk"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type fluxTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l fluxTool) Name() string {
	return "flux"
}

func (l fluxTool) ShortDesc() string {
	return "The GitOps Kubernetes operator"
}

func (l fluxTool) LongDesc() string {
	return "Flux is a tool that automatically ensures that the state of a cluster matches the config in git. It uses an operator in the cluster to trigger deployments inside Kubernetes, which means you don't need a separate CD tool. It monitors all relevant image repositories, detects new images, triggers deployments and updates the desired running configuration based on that (and a configurable policy).\n\n"
}

func (l fluxTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/fluxctl_", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += l.os + "_" + l.arch
	case l.os == "windows" && l.arch == "amd64":
		url += l.os + "_" + l.arch
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "arm64":
		url += l.os + "_" + l.arch
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l fluxTool) Versions(max uint) ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(
		context.Background(), l.Owner, l.Repo, nil,
	)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && uint(len(releases)) < max {
		r, resp, err = client.Repositories.ListReleases(
			context.Background(), l.Owner, l.Repo, &github.ListOptions{
				Page:    resp.NextPage,
				PerPage: int(max) - len(releases),
			},
		)
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(
		releases, func(release *github.RepositoryRelease) bool {
			return !release.GetPrerelease() &&
				!strings.Contains(release.GetTagName(), "rc") &&
				!strings.Contains(release.GetName(), "rc") &&
				!strings.Contains(release.GetTagName(), "helm")
		},
	).([]*github.RepositoryRelease)

	vs := make([]*semver.Version, len(releases))
	for i, release := range releases {
		v, err := semver.NewVersion(release.GetTagName())
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %w", err)
		}

		vs[i] = v
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	// dont need too many releases
	if uint(len(vs)) > max {
		vs = vs[:max]
	}

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return fluxTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("fluxcd", "flux"),
	}
}
