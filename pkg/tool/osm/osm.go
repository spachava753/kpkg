package osm

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type osmTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l osmTool) Extract(artifactPath, _ string) (string, error) {
	// osm releases contain the binary and some files. Pick only the binary

	// expect given path to be to contain bin folder
	binDirPath := filepath.Join(artifactPath, fmt.Sprintf("%s-%s", l.os, l.arch))
	binDirPathInfo, err := os.Stat(binDirPath)
	if err != nil {
		return "", fmt.Errorf("could not extract binary: %w", err)
	}
	if !binDirPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}
	binaryPath := filepath.Join(binDirPath, l.Name())
	binaryPathInfo, err := os.Stat(binaryPath)
	if err != nil {
		return "", err
	}
	if binaryPathInfo.IsDir() {
		return "", fmt.Errorf("could not extract binary: %w", fmt.Errorf("path %s is not a directory", binDirPath))
	}

	return binaryPath, err
}

func (l osmTool) Name() string {
	return "osm"
}

func (l osmTool) ShortDesc() string {
	return `Open Service Mesh (OSM) is a lightweight, extensible, cloud native service mesh that allows users to 
uniformly manage, secure, and get out-of-the-box observability features for highly dynamic microservice environments`
}

func (l osmTool) LongDesc() string {
	return `Open Service Mesh (OSM) is a lightweight, extensible, Cloud Native service mesh that allows users to 
uniformly manage, secure, and get out-of-the-box observability features for highly dynamic microservice environments. 
The OSM project builds on the ideas and implementations of many cloud native ecosystem projects including 
Linkerd, Istio, Consul, Envoy, Kuma, Helm, and the SMI specification`
}

func (l osmTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// only supports amd64
	if l.arch != "amd64" {
		return "", err
	}

	url := fmt.Sprintf("%sv%s/osm-v%s-%s-%s.zip", l.MakeReleaseUrl(), version, version, l.os, l.arch)
	switch {
	case l.os == "darwin",
		l.os == "windows",
		l.os == "linux":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return osmTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("openservicemesh", "osm", 20),
	}
}
