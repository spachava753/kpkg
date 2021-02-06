package kind

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type kubectlTool struct {
	basePath,
	arch,
	os string
}

func (l kubectlTool) Name() string {
	return "kind"
}

func (l kubectlTool) ShortDesc() string {
	return "Kubernetes IN Docker - local clusters for testing Kubernetes"
}

func (l kubectlTool) LongDesc() string {
	return `kind is a tool for running local Kubernetes clusters using Docker container "nodes". kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.`
}

func (l kubectlTool) MakeUrl(version string) (string, error) {
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "ppc64le":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://github.com/kubernetes-sigs/kind/releases/download/%s/kind-%s-%s", version, l.os, l.arch)
	return url, nil
}

func (l kubectlTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kind", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "kubernetes-sigs", "kind", &github.ListOptions{
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
		if !r.GetPrerelease() {
			versions = append(versions, *r.Name)
		}
	}

	// sort results
	return versions, nil
}

func MakeBinary(basePath, os, arch string) tool.Binary {
	return kubectlTool{
		basePath: basePath,
		arch:     arch,
		os:       os,
	}
}
