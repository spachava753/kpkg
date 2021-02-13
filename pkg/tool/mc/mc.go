package mc

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/thoas/go-funk"
	"strings"
)

type mcTool struct {
	arch,
	os string
}

func (l mcTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l mcTool) Name() string {
	return "mc"
}

func (l mcTool) ShortDesc() string {
	return "bootstrap Kubernetes with k3s over SSH < 1 min 🚀"
}

func (l mcTool) LongDesc() string {
	return `k3sup is a light-weight utility to get from zero to KUBECONFIG with k3s on any local or remote VM. 
All you need is ssh access and the k3sup binary to get kubectl access immediately`
}

func (l mcTool) MakeUrl(version string) (string, error) {
	// minio client doesn't use semantic versioning
	url := fmt.Sprintf("https://dl.min.io/client/mc/release/%s-%s/archive/mc.%s", l.os, l.arch, version)
	versions, err := l.Versions()
	if err != nil {
		return "", err
	}
	if version == versions[0] {
		url = fmt.Sprintf("https://dl.min.io/client/mc/release/%s-%s/mc", l.os, l.arch)
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "mips64",
		l.os == "linux" && l.arch == "ppc64le",
		l.os == "linux" && l.arch == "s390x":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l mcTool) Versions() ([]string, error) {
	// minio client doesn't use semantic versioning
	max := 20
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "minio", "minio", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < max {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "minio", "minio", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: max - len(releases),
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}

	releases = funk.Filter(releases, func(release *github.RepositoryRelease) bool {
		return !release.GetPrerelease() && strings.Contains(release.GetTagName(), "RELEASE")
	}).([]*github.RepositoryRelease)

	versions := make([]string, 0, len(releases))
	for _, v := range releases {
		versions = append(versions, v.GetTagName())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return mcTool{
		arch: arch,
		os:   os,
	}
}
