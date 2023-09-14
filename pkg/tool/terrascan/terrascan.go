package terrascan

import (
	"fmt"
	"path/filepath"

	"github.com/Masterminds/semver"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type terrascanTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l terrascanTool) Extract(artifactPath, _ string) (string, error) {
	return filepath.Join(artifactPath, l.Name()), nil
}

func (l terrascanTool) Name() string {
	return "terrascan"
}

func (l terrascanTool) ShortDesc() string {
	return "Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure"
}

func (l terrascanTool) LongDesc() string {
	return `Terrascan detects security vulnerabilities and compliance violations across your Infrastructure as Code. 
Mitigate risks before provisioning cloud native infrastructure. Run locally or integrate with your CI\CD.`
}

func (l terrascanTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf(
		"%sv%s/terrascan_%s_", l.MakeReleaseUrl(), version, version,
	)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "Darwin_x86_64"
	case l.os == "windows" && l.arch == "386":
		url += "Windows_i386"
	case l.os == "windows" && l.arch == "amd64":
		url += "Windows_x86_64"
	case l.os == "linux" && l.arch == "386":
		url += "Linux_i386"
	case l.os == "linux" && l.arch == "amd64":
		url += "Linux_x86_64"
	case l.os == "linux" && l.arch == "arm64":
		url += "Linux_arm64"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return terrascanTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("accurics", "terrascan"),
	}
}
