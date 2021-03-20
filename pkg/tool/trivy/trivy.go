package trivy

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type trivyTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l trivyTool) Extract(artifactPath, _ string) (string, error) {
	binPath := filepath.Join(artifactPath, l.Name())
	if _, err := os.Stat(binPath); err != nil {
		return "", err
	}
	return binPath, nil
}

func (l trivyTool) Name() string {
	return "trivy"
}

func (l trivyTool) ShortDesc() string {
	return "A Simple and Comprehensive Vulnerability Scanner for Container Images, Git Repositories and Filesystems. Suitable for CI"
}

func (l trivyTool) LongDesc() string {
	return `Trivy (tri pronounced like trigger, vy pronounced like envy) is a simple and comprehensive vulnerability 
scanner for containers and other artifacts. A software vulnerability is a glitch, flaw, or weakness present in the 
software or in an Operating System. Trivy detects vulnerabilities of OS packages (Alpine, RHEL, CentOS, etc.) and 
application dependencies (Bundler, Composer, npm, yarn, etc.). Trivy is easy to use. Just install the binary and you're 
ready to scan. All you need to do for scanning is to specify a target such as an image name of the container`
}

func (l trivyTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/trivy_%s_", l.MakeReleaseUrl(), version, version)
	switch l.os {
	case "linux":
		url += "Linux-"
		switch l.arch {
		case "386":
			url += "32bit"
		case "amd64":
			url += "64bit"
		case "arm":
			url += "ARM"
		case "arm64":
			url += "ARM64"
		case "ppc64le":
			url += "PPC64LE"
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	case "freebsd":
		url += "FreeBSD-"
		switch l.arch {
		case "386":
			url += "32bit"
		case "amd64":
			url += "64bit"
		case "arm":
			url += "ARM"
		case "arm64":
			url += "ARM64"
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	case "openbsd":
		url += "FreeBSD-"
		switch l.arch {
		case "386":
			url += "32bit"
		case "amd64":
			url += "64bit"
		case "arm":
			url += "ARM"
		case "arm64":
			url += "ARM64"
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	case "darwin":
		if l.arch != "amd64" {
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
		url += "macOS-64bit"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return trivyTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("aquasecurity", "trivy", 20),
	}
}
