package terraform

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type terraformTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l terraformTool) Extract(artifactPath, _ string) (string, error) {
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

func (l terraformTool) Name() string {
	return "terraform"
}

func (l terraformTool) ShortDesc() string {
	return "Write infrastructure as code using declarative configuration files"
}

func (l terraformTool) LongDesc() string {
	return "Terraform is an open-source infrastructure as code software tool that provides a consistent CLI workflow to manage hundreds of cloud services"
}

func (l terraformTool) MakeUrl(version string) (string, error) {
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
	url := fmt.Sprintf(
		"https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip",
		version, version, l.os, l.arch,
	)
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return terraformTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("hashicorp", "terraform"),
	}
}
