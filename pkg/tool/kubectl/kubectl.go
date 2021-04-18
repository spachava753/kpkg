package kubectl

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

const (
	owner = "kubernetes"
	repo  = "kubectl"
	max   = 100
)

// doesn't need GithubReleaseTool
type kubectlTool struct {
	arch,
	os string
}

func (l kubectlTool) Extract(artifactPath, _ string) (string, error) {
	return artifactPath, nil
}

func (l kubectlTool) Name() string {
	return "kubectl"
}

func (l kubectlTool) ShortDesc() string {
	return "kubectl is a cli to communicate k8s clusters"
}

func (l kubectlTool) LongDesc() string {
	return "kubectl is a cli to communicate k8s clusters"
}

func (l kubectlTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl", version, l.os, l.arch)
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func (l kubectlTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	tags, resp, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryTag
	for resp != nil && resp.NextPage != resp.LastPage && uint(len(tags)) < max {
		r, resp, err = client.Repositories.ListTags(context.Background(), owner, repo, &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: max - len(tags),
		})
		if err != nil {
			return nil, err
		}
		tags = append(tags, r...)
	}

	tags = funk.Filter(tags, func(release *github.RepositoryTag) bool {
		return !strings.Contains(release.GetName(), "rc") &&
			!strings.Contains(release.GetName(), "alpha") &&
			!strings.Contains(release.GetName(), "beta")
	}).([]*github.RepositoryTag)

	vs := make([]*semver.Version, len(tags))
	for i, release := range tags {
		tagName := release.GetName()
		tagName = tagName[2:]
		tagName = "v1" + tagName
		v, err := semver.NewVersion(tagName)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %w", err)
		}
		vs[i] = v
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))

	// dont need too many tags
	if uint(len(vs)) > max {
		vs = vs[:max]
	}

	versions := make([]string, 0, len(vs))
	for _, v := range vs {
		versions = append(versions, v.String())
	}

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return kubectlTool{
		arch: arch,
		os:   os,
	}
}
