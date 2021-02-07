package helm

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
	"strings"
)

type helmTool struct {
	arch,
	os string
}

func (l helmTool) Extract(artifactPath, version string) (string, error) {
	// helm releases contain the binary, LICENSE and a README. Pick only the binary
	var binaryPath string
	err := filepath.Walk(artifactPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(filepath.Base(path), "helm") &&
			info != nil &&
			!info.IsDir() &&
			binaryPath == "" {
			binaryPath, err = filepath.Abs(path)
			return err
		}
		return nil
	})

	return binaryPath, err
}

func (l helmTool) Name() string {
	return "helm"
}

func (l helmTool) ShortDesc() string {
	return "The Kubernetes Package Manager"
}

func (l helmTool) LongDesc() string {
	return "Helm is a tool for managing Charts. Charts are packages of pre-configured Kubernetes resources"
}

func (l helmTool) MakeUrl(version string) (string, error) {
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "arm":
		fallthrough
	case l.os == "linux" && l.arch == "i386":
		fallthrough
	case l.os == "linux" && l.arch == "s390x":
		fallthrough
	case l.os == "linux" && l.arch == "ppc64le":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz", version, l.os, l.arch), nil
}

func (l helmTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "helm", "helm", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "helm", "helm", &github.ListOptions{
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

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return helmTool{
		arch: arch,
		os:   os,
	}
}
