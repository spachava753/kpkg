package kail

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

type kailTool struct {
	arch,
	os string
}

func (l kailTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l kailTool) Name() string {
	return "kail"
}

func (l kailTool) ShortDesc() string {
	return "kubernetes log viewer"
}

func (l kailTool) LongDesc() string {
	return `Kubernetes tail. Streams logs from all containers of all matched pods. 
Match pods by service, replicaset, deployment, and others. 
Adjusts to a changing cluster - pods are added and removed from 
logging as they fall in or out of the selection.`
}

func (l kailTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return fmt.Sprintf("https://github.com/boz/kail/releases/download/v%s/kail_%s_%s_%s.tar.gz", version, version, l.os, l.arch), nil
}

func (l kailTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "boz", "kail", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "boz", "kail", &github.ListOptions{
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

	// sort results
	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kailTool{
		arch: arch,
		os:   os,
	}
}
