package terraform

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

type terraformTool struct {
	arch,
	os string
}

func (l terraformTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l terraformTool) Name() string {
	return "terraform"
}

func (l terraformTool) ShortDesc() string {
	return "hashicorp CLI plugin for extended build capabilities with BuildKit"
}

func (l terraformTool) LongDesc() string {
	return `terraform is a hashicorp CLI plugin for extended build capabilities with BuildKit.`
}

func (l terraformTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "386":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		fallthrough
	case l.os == "linux" && l.arch == "386":
		fallthrough
	case l.os == "freebsd" && l.arch == "amd64":
		fallthrough
	case l.os == "freebsd" && l.arch == "arm":
		fallthrough
	case l.os == "freebsd" && l.arch == "386":
		fallthrough
	case l.os == "openbsd" && l.arch == "amd64":
		fallthrough
	case l.os == "openbsd" && l.arch == "386":
		fallthrough
	case l.os == "solaris" && l.arch == "amd64":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip", version, version, l.os, l.arch)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func (l terraformTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "hashicorp", "terraform", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "hashicorp", "terraform", &github.ListOptions{
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

	// dont need too many releases
	vs = vs[:20]

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return terraformTool{
		arch: arch,
		os:   os,
	}
}
