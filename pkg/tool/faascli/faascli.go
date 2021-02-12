package faascli

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

type faasCliTool struct {
	arch,
	os string
}

func (l faasCliTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l faasCliTool) Name() string {
	return "faas-cli"
}

func (l faasCliTool) ShortDesc() string {
	return "openfaas CLI plugin for extended build capabilities with BuildKit"
}

func (l faasCliTool) LongDesc() string {
	return `faas-cli is a openfaas CLI plugin for extended build capabilities with BuildKit.`
}

func (l faasCliTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/openfaas/faas-cli/releases/download/%s/faas-cli", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url = url + "-darwin"
		break
	case l.os == "windows" && l.arch == "amd64":
		url = url + ".exe"
		break
	case l.os == "linux" && l.arch == "arm":
		url = url + "-armhf"
		break
	case l.os == "linux" && l.arch == "amd64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l faasCliTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "openfaas", "faas-cli", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "openfaas", "faas-cli", &github.ListOptions{
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
	return faasCliTool{
		arch: arch,
		os:   os,
	}
}
