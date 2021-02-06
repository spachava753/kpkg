package linkerd2

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool"
	"sort"
	"strings"
)

type linkerd2Tool struct {
	basePath,
	arch,
	os string
	fileFetcher download.FileFetcher
}

func (l linkerd2Tool) Name() string {
	return "linkerd2"
}

func (l linkerd2Tool) ShortDesc() string {
	return "linkerd2 is a cli to install linkerd2 service mesh"
}

func (l linkerd2Tool) LongDesc() string {
	return "linkerd2 is a cli to install linkerd2 service mesh"
}

func (l linkerd2Tool) MakeUrl(version string) (string, error) {
	// install the latest stable binary
	if version == "latest" {
		versions, err := l.Versions()
		if err != nil {
			return "", err
		}
		for _, v := range versions {
			if strings.Contains(v, "stable") {
				version = v
				break
			}
		}
	}
	switch l.os {
	case "darwin":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-darwin", version, version), nil
	case "windows":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-windows.exe", version, version), nil
	case "linux":
		switch l.arch {
		case "amd64":
			fallthrough
		case "arm":
			fallthrough
		case "arm64":
			return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-linux-%s", version, version, l.arch), nil
		default:
			return "", fmt.Errorf("unsupported architecture: %s", l.arch)
		}
	}
	return "", fmt.Errorf("unsupported os: %s", l.os)
}

func (l linkerd2Tool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "linkerd", "linkerd2", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "linkerd", "linkerd2", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}
	versions := make([]string, 0, len(releases))
	for _, r := range releases {
		versions = append(versions, *r.Name)
	}

	// sort results
	return sortVersions(versions), nil
}

func sortVersions(versions []string) []string {
	if len(versions) < 2 {
		return versions
	}
	// first split versions into "stable" and "edge"
	var stable, edge []string
	for _, v := range versions {
		if strings.Contains(v, "stable") {
			stable = append(stable, v)
			continue
		}
		if strings.Contains(v, "edge") {
			edge = append(edge, v)
			continue
		}
	}
	stableSort := sort.StringSlice(stable)
	sort.Sort(sort.Reverse(stableSort))
	edgeSort := sort.StringSlice(edge)
	sort.Sort(sort.Reverse(edgeSort))
	return append(stableSort, []string(edgeSort)...)
}

func (l linkerd2Tool) makeUrl(version string) (string, error) {
	// install the latest stable binary
	if version == "latest" {
		versions, err := l.Versions()
		if err != nil {
			return "", err
		}
		for _, v := range versions {
			if strings.Contains(v, "stable") {
				version = v
				break
			}
		}
	}
	switch l.os {
	case "darwin":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-darwin", version, version), nil
	case "windows":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-windows.exe", version, version), nil
	case "linux":
		switch l.arch {
		case "amd64":
			fallthrough
		case "arm":
			fallthrough
		case "arm64":
			return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-linux-%s", version, version, l.arch), nil
		default:
			return "", fmt.Errorf("unsupported architecture: %s", l.arch)
		}
	}
	return "", fmt.Errorf("unsupported os: %s", l.os)
}

func MakeBinary(basePath, os, arch string, fetcher download.FileFetcher) tool.Binary {
	return linkerd2Tool{
		basePath:    basePath,
		arch:        arch,
		os:          os,
		fileFetcher: fetcher,
	}
}
