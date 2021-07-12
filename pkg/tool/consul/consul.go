package consul

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type consulTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l consulTool) Extract(artifactPath, _ string) (string, error) {
	binaryPath := filepath.Join(artifactPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binaryPathInfo))
	}

	return binaryPath, err
}

func (l consulTool) Name() string {
	return "consul"
}

func (l consulTool) ShortDesc() string {
	return "Consul is a distributed, highly available, and data center aware solution to connect and configure applications across dynamic, distributed infrastructure"
}

func (l consulTool) LongDesc() string {
	return "Consul is a distributed, highly available, and data center aware solution to connect and configure applications across dynamic, distributed infrastructure"
}

func (l consulTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "386",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "arm",
		l.os == "freebsd" && l.arch == "386",
		l.os == "openbsd" && l.arch == "amd64",
		l.os == "openbsd" && l.arch == "386",
		l.os == "solaris" && l.arch == "amd64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("https://releases.hashicorp.com/consul/%s/consul_%s_%s_%s.zip", version, version, l.os, l.arch)
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return consulTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("hashicorp", "consul", 20),
	}
}
