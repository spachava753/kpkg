package hugo

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type hugoTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l hugoTool) Extract(artifactPath, _ string) (string, error) {
	// hugo releases contain the binary and some examples. Pick only the binary
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is a dir", binaryPath))
	}

	return binaryPath, err
}

func (l hugoTool) Name() string {
	return "hugo"
}

func (l hugoTool) ShortDesc() string {
	return "The world’s fastest framework for building websites"
}

func (l hugoTool) LongDesc() string {
	return `Hugo is a static HTML and CSS website generator written in Go. 
It is optimized for speed, ease of use, and configurability. 
Hugo takes a directory with content and templates and renders them into a full HTML website`
}

func (l hugoTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_macOS-64bit", version)
	case l.os == "windows" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_Windows-64bit", version)
	case l.os == "windows" && l.arch == "386":
		url += fmt.Sprintf("hugo_%s_Windows-32bit", version)
	case l.os == "linux" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_Linux-64bit", version)
	case l.os == "linux" && l.arch == "arm64":
		url += fmt.Sprintf("hugo_%s_Linux-ARM64", version)
	case l.os == "linux" && l.arch == "arm":
		url += fmt.Sprintf("hugo_%s_Linux-ARM", version)
	case l.os == "linux" && l.arch == "386":
		url += fmt.Sprintf("hugo_%s_Linux-32bit", version)
	case l.os == "freebsd" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_FreeBSD-64bit", version)
	case l.os == "freebsd" && l.arch == "arm":
		url += fmt.Sprintf("hugo_%s_FreeBSD-ARM", version)
	case l.os == "freebsd" && l.arch == "386":
		url += fmt.Sprintf("hugo_%s_FreeBSD-64bit", version)
	case l.os == "dragonfly" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_DragonFlyBSD-64bit", version)
	case l.os == "openbsd" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_OpenBSD-64bit", version)
	case l.os == "openbsd" && l.arch == "arm":
		url += fmt.Sprintf("hugo_%s_OpenBSD-ARM", version)
	case l.os == "openbsd" && l.arch == "386":
		url += fmt.Sprintf("hugo_%s_OpenBSD-32bit", version)
	case l.os == "netbsd" && l.arch == "amd64":
		url += fmt.Sprintf("hugo_%s_NetBSD-64bit", version)
	case l.os == "netbsd" && l.arch == "arm":
		url += fmt.Sprintf("hugo_%s_NetBSD-ARM", version)
	case l.os == "netbsd" && l.arch == "386":
		url += fmt.Sprintf("hugo_%s_NetBSD-32bit", version)
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return hugoTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("gohugoio", "hugo", 20),
	}
}
