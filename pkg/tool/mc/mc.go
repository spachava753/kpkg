package mc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/thoas/go-funk"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
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
	return "MinIO Client (mc) provides a modern alternative to UNIX commands like ls, cat, cp, mirror, diff, find etc"
}

func (l mcTool) LongDesc() string {
	return `MinIO Client (mc) provides a modern alternative to UNIX commands like ls, cat, cp, mirror, diff, find etc. 
It supports filesystems and Amazon S3 compatible cloud storage service (AWS Signature v2 and v4)`
}

func (l mcTool) MakeUrl(version string) (string, error) {
	// minio client doesn't use semantic versioning
	url := fmt.Sprintf(
		"https://dl.min.io/client/mc/release/%s-%s/archive/mc.%s", l.os, l.arch,
		version,
	)
	versions, err := l.Versions(1)
	if err != nil {
		return "", err
	}
	if version == versions[0] {
		url = fmt.Sprintf(
			"https://dl.min.io/client/mc/release/%s-%s/mc", l.os, l.arch,
		)
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
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l mcTool) Versions(max uint) ([]string, error) {
	// minio client doesn't use semantic versioning
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(
		context.Background(), "minio", "minio", nil,
	)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) < int(max) {
		r, resp, err = client.Repositories.ListReleases(
			context.Background(), "minio", "minio", &github.ListOptions{
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
		releases, func(release *github.RepositoryRelease) bool {
			return !release.GetPrerelease() && strings.Contains(
				release.GetTagName(), "RELEASE",
			)
		},
	).([]*github.RepositoryRelease)

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
