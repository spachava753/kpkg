package kubebench

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

type kubeBenchTool struct {
	arch,
	os string
}

func (l kubeBenchTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l kubeBenchTool) Name() string {
	return "kube-bench"
}

func (l kubeBenchTool) ShortDesc() string {
	return "Checks whether Kubernetes is deployed according to security best practices as defined in the CIS Kubernetes Benchmark"
}

func (l kubeBenchTool) LongDesc() string {
	return `kube-bench is a Go application that checks whether Kubernetes is deployed securely by running 
the checks documented in the CIS Kubernetes Benchmark`
}

func (l kubeBenchTool) MakeUrl(version string) (string, error) {
	if l.os != "linux" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/aquasecurity/kube-bench/releases/download/v%s/kube-bench_%s_%s_%s.tar.gz", version, version, l.os, l.arch)

	switch {
	case l.arch == "amd64", l.arch == "arm64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l kubeBenchTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "aquasecurity", "kube-bench", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "aquasecurity", "kube-bench", &github.ListOptions{
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
	return kubeBenchTool{
		arch: arch,
		os:   os,
	}
}
