package inletsctl

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"path/filepath"
	"sort"
	"strings"
)

type inletsctlTool struct {
	arch,
	os string
}

func (l inletsctlTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l inletsctlTool) Name() string {
	return "inletsctl"
}

func (l inletsctlTool) ShortDesc() string {
	return "The fastest way to create self-hosted exit-servers"
}

func (l inletsctlTool) LongDesc() string {
	return `inletsctl automates the task of creating an exit-node on cloud infrastructure. 
Once provisioned, you'll receive a command to connect with. 
You can use this tool whether you want to use inlets or inlets-pro for L4 TCP`
}

func (l inletsctlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/inlets/inletsctl/releases/download/%s/inletsctl", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url = url + "darwin"
		break
	case l.os == "windows" && l.arch == "amd64":
		url = url + ".exe"
		break
	case l.os == "linux" && l.arch == "amd64":
		// do nothing
		break
	case l.os == "linux" && l.arch == "arm64":
		url = url + "arm64"
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	return url + ".tgz", nil
}

func (l inletsctlTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "inlets", "inletsctl", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "inlets", "inletsctl", &github.ListOptions{
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
	return inletsctlTool{
		arch: arch,
		os:   os,
	}
}
