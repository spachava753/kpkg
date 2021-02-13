package dockercompose

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

type composeTool struct {
	arch,
	os string
}

func (l composeTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l composeTool) Name() string {
	return "docker-compose"
}

func (l composeTool) ShortDesc() string {
	return "Define and run multi-container applications with Docker"
}

func (l composeTool) LongDesc() string {
	return `Define and run multi-container applications with DockerDocker Compose is a tool 
for running multi-container applications on Docker defined using the Compose file format. 
A Compose file is used to define how the one or more containers that make up your application 
are configured. Once you have a Compose file, you can create and start your application with 
a single command: docker-compose up.`
}

func (l composeTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/docker-compose-%s-x86_64", version, strings.Title(l.os))
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func (l composeTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "docker", "compose", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "docker", "compose", &github.ListOptions{
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
	return composeTool{
		arch: arch,
		os:   os,
	}
}
