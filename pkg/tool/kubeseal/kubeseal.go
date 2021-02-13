package kubeseal

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

type kubesealTool struct {
	arch,
	os string
}

func (l kubesealTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l kubesealTool) Name() string {
	return "kubeseal"
}

func (l kubesealTool) ShortDesc() string {
	return "bootstrap Kubernetes with k3s over SSH < 1 min 🚀"
}

func (l kubesealTool) LongDesc() string {
	return `k3sup is a light-weight utility to get from zero to KUBECONFIG with k3s on any local or remote VM. 
All you need is ssh access and the k3sup binary to get kubectl access immediately`
}

func (l kubesealTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("https://github.com/bitnami-labs/sealed-secrets/releases/download/v%s/kubeseal", version)
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
		url += fmt.Sprintf("-%s-%s", l.os, l.arch)
		break
	case l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm":
		url += fmt.Sprintf("-%s", l.arch)
		break
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l kubesealTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "bitnami-labs", "sealed-secrets", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < 20 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "bitnami-labs", "sealed-secrets", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 20 - len(releases),
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(releases, func(release *github.RepositoryRelease) bool {
		return !release.GetPrerelease() && !strings.Contains(release.GetTagName(), "helm")
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
	return kubesealTool{
		arch: arch,
		os:   os,
	}
}
