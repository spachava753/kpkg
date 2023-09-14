package vagrant

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v33/github"
	"github.com/thoas/go-funk"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type vagrantTool struct {
	arch,
	os string
}

func (l vagrantTool) Extract(artifactPath, _ string) (string, error) {
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf(
			"could not extract binary: %w",
			fmt.Errorf("path %s is not a directory", binaryPathInfo),
		)
	}

	return binaryPath, err
}

func (l vagrantTool) Name() string {
	return "vagrant"
}

func (l vagrantTool) ShortDesc() string {
	return "Vagrant is a tool for building and distributing development environments"
}

func (l vagrantTool) LongDesc() string {
	return `Development environments managed by Vagrant can run on local virtualized platforms such as 
VirtualBox or VMware, in the cloud via AWS or OpenStack, or in containers such as with Docker or raw LXC`
}

func (l vagrantTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	// only supports linux/amd64
	if l.os != "linux" || l.arch != "amd64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	url := fmt.Sprintf(
		"https://releases.hashicorp.com/vagrant/%s/vagrant_%s_%s_%s.zip",
		version, version, l.os, l.arch,
	)
	return url, nil
}

func (l vagrantTool) Versions(max uint) ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListTags(
		context.Background(), "hashicorp", "vagrant", nil,
	)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryTag
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < int(max) {
		r, resp, err = client.Repositories.ListTags(
			context.Background(), "hashicorp", "vagrant", &github.ListOptions{
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
		releases, func(release *github.RepositoryTag) bool {
			return !strings.Contains(release.GetName(), "rc")
		},
	).([]*github.RepositoryTag)

	vs := make([]*semver.Version, len(releases))
	for i, release := range releases {
		v, err := semver.NewVersion(release.GetName())
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %w", err)
		}

		vs[i] = v
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	// dont need too many releases
	if len(vs) > int(max) {
		vs = vs[:max]
	}

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return vagrantTool{
		arch: arch,
		os:   os,
	}
}
