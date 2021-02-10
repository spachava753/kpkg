package istioctl

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
	"sort"
)

type istioctlTool struct {
	arch,
	os string
}

func (l istioctlTool) Extract(artifactPath, version string) (string, error) {
	// istioctl releases contain the binary and some examples. Pick only the binary

	// expect given path to be to contain bin folder
	binDirPath := filepath.Join(artifactPath, fmt.Sprintf("istio-%s", version), "bin")
	binDirPathInfo, err := os.Stat(binDirPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if !binDirPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}
	binaryPath := filepath.Join(binDirPath, "istioctl")
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}

	return binaryPath, err
}

func (l istioctlTool) Name() string {
	return "istioctl"
}

func (l istioctlTool) ShortDesc() string {
	return "Connect, secure, control, and observe services"
}

func (l istioctlTool) LongDesc() string {
	return `Istio is an open platform for providing a uniform way to integrate microservices, 
manage traffic flow across microservices, enforce policies and aggregate telemetry data. 
Istio's control plane provides an abstraction layer over the underlying cluster management platform, 
such as Kubernetes.`
}

func (l istioctlTool) MakeUrl(version string) (string, error) {
	var url string

	// based of the script "curl -L https://istio.io/downloadIstio"
	archSupported := "1.5"
	// compare major and minor version
	versionCmp := []string{archSupported, version[:3]}
	sort.Sort(sort.Reverse(sort.StringSlice(versionCmp)))
	if versionCmp[0] == archSupported {
		// version is less than 1.5, have to use different url for downloading
		if l.arch != "amd64" {
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
		switch {
		case l.os == "darwin":
			url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-osx.tar.gz", version, version)
		case l.os == "windows":
			url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-win.zip", version, version)
		case l.os == "linux":
			url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-linux.tar.gz", version, version)
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
		return url, nil
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-osx.tar.gz", version, version)
	case l.os == "windows" && l.arch == "amd64":
		url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-win.zip", version, version)
	case l.os == "linux" && l.arch == "armv7":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		url = fmt.Sprintf("https://github.com/istio/istio/releases/download/%s/istio-%s-linux-%s.tar.gz", version, version, l.arch)
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func (l istioctlTool) Versions() ([]string, error) {
	client := github.NewClient(nil)
	var resp *github.Response
	releases, resp, err := client.Repositories.ListReleases(context.Background(), "istio", "istio", nil)
	if err != nil {
		return nil, err
	}
	var r []*github.RepositoryRelease
	for resp != nil && resp.NextPage != resp.LastPage && len(releases) <= 30 {
		r, resp, err = client.Repositories.ListReleases(context.Background(), "istio", "istio", &github.ListOptions{
			Page:    resp.NextPage,
			PerPage: 20,
		})
		if err != nil {
			return nil, err
		}
		releases = append(releases, r...)
	}
	versions := make([]string, 0, len(releases))
	for _, r := range releases {
		if !r.GetPrerelease() {
			versions = append(versions, r.GetTagName())
		}
	}
	// keep only the first 30 releases
	if len(versions) > 30 {
		versions = versions[:31]
	}

	// sort results
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return istioctlTool{
		arch: arch,
		os:   os,
	}
}
