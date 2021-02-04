package linkerd2

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type linkerd2Tool struct {
	basePath,
	arch,
	os string
}

func (l linkerd2Tool) Install(version string, force bool) (s string, err error) {
	// check that the version exists
	versions, err := l.Versions(false)
	if err != nil {
		return "", err
	}

	if version != "latest" {
		if !util.ContainsString(versions, version) {
			return "", fmt.Errorf("version %s is not valid for binary linkerd2", version)
		}
	}

	// assign latest "stable" version for latest
	if version == "latest" {
		version, err = func(versions []string) (string, error) {
			for _, v := range versions {
				if strings.Contains(v, "stable") {
					return v, nil
				}
			}
			return "", fmt.Errorf("could not find latest stable version")
		}(versions)
		if err != nil {
			return "", err
		}
	}

	ld2Path := filepath.Join(l.basePath, "linkerd2")
	ld2VersionPath := filepath.Join(ld2Path, version)
	ld2BinaryPath := filepath.Join(ld2VersionPath, "linkerd2")
	ld2BinPath := filepath.Join(l.basePath, "bin", "linkerd2")

	// check if installed already
	if ld2Info, err := os.Stat(ld2Path); err == nil {
		if !ld2Info.IsDir() {
			return "", fmt.Errorf("cannot install, path %s contains a file", ld2Path)
		}
		if ld2VersionInfo, err := os.Stat(ld2VersionPath); err == nil {
			if !ld2VersionInfo.IsDir() {
				return "", fmt.Errorf("cannot install, path %s contains a file", ld2VersionPath)
			}
			if ld2BinaryInfo, err := os.Stat(ld2BinaryPath); err == nil {
				if ld2BinaryInfo.IsDir() {
					return "", fmt.Errorf("cannot install, path %s contains a dir", ld2BinaryPath)
				}
				// since we already have it installed, set the symlink to this
				if !force {
					if err = os.Remove(ld2BinPath); err != nil {
						return "", fmt.Errorf("could not remove symlink to path %s: %w", ld2BinPath, err)
					}
					return "", os.Symlink(ld2BinaryPath, filepath.Join(l.basePath, "bin", "linkerd2"))
				}
				// since force is enabled, remove the file and continue
				if err := os.Remove(ld2BinaryInfo.Name()); err != nil {
					return "", err
				}
			}
		}
	}

	// construct the url to fetch the release
	url, err := l.makeUrl(version)
	if err != nil {
		return "", err
	}

	// create a temp file to download the CLI to
	tmpF, err := ioutil.TempFile(os.TempDir(), "linkerd2")
	if err != nil {
		return "", err
	}
	defer func() {
		if e := tmpF.Close(); e != nil {
			err = e
		}
	}()
	defer func() {
		if e := os.Remove(tmpF.Name()); e != nil {
			err = e
		}
	}()

	// download CLI
	err = download.FetchFile(url, tmpF)
	if err != nil {
		return "", err
	}

	// copy to our bin path

	// create binary file
	if _, err := os.Stat(ld2VersionPath); os.IsNotExist(err) {
		if err := os.MkdirAll(ld2VersionPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	// copy the downloaded binary to path
	contents, err := ioutil.ReadFile(tmpF.Name())
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(ld2BinaryPath, contents, os.ModePerm); err != nil {
		return "", err
	}

	// create symlink to bin path
	if info, err := os.Stat(ld2BinPath); err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("could not remove symlink, path %s is a dir", ld2BinPath)
		}
		if err = os.Remove(ld2BinPath); err != nil {
			return "", fmt.Errorf("could not remove symlink to path %s: %w", ld2BinPath, err)
		}
	}
	return ld2BinaryPath, os.Symlink(ld2BinaryPath, filepath.Join(l.basePath, "bin", "linkerd2"))
}

func (l linkerd2Tool) Versions(installedOnly bool) ([]string, error) {
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
		releases = append(releases, r...)
	}
	versions := make([]string, 0, len(releases))
	for _, r := range releases {
		versions = append(versions, *r.Name)
	}
	return versions, nil
}

func (l linkerd2Tool) makeUrl(version string) (string, error) {
	// install the latest stable binary
	if version == "latest" {
		versions, err := l.Versions(false)
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

func MakeBinary(basePath, os, arch string) tool.Binary {
	return linkerd2Tool{
		basePath: basePath,
		arch:     arch,
		os:       os,
	}
}
