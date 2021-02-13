package tool

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	"github.com/thoas/go-funk"
	"sort"
	"strings"
)

type GithubReleaseTool struct {
	owner, repo string
	max         uint
}

func (l GithubReleaseTool) MakeReleaseUrl() string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/", l.owner, l.repo)
}

func (l GithubReleaseTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l GithubReleaseTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), l.owner, l.repo, nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && uint(len(releases)) < l.max {
		r, resp, err = client.Repositories.ListReleases(context.Background(), l.owner, l.repo, &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: int(l.max) - len(releases),
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

	// dont need too many releases
	if uint(len(vs)) > l.max {
		vs = vs[:l.max]
	}

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeGithubReleaseTool(org, repo string, max uint) GithubReleaseTool {
	return GithubReleaseTool{
		org, repo, max,
	}
}
