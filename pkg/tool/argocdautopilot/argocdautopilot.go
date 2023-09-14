package argocdautopilot

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"os"
	"path/filepath"
)

type argocdAutopilotTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l argocdAutopilotTool) Extract(artifactPath, _ string) (string, error) {
	binaryPath := filepath.Join(artifactPath, l.Name()+"-"+l.os+"-"+l.arch)
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

func (l argocdAutopilotTool) Name() string {
	return "argocd-autopilot"
}

func (l argocdAutopilotTool) ShortDesc() string {
	return "The Argo-CD Autopilot is a tool which offers an opinionated way of installing Argo-CD and managing GitOps repositories"
}

func (l argocdAutopilotTool) LongDesc() string {
	return "New users to GitOps and Argo CD are not often sure how they should structure their repos, add applications, promote apps across environments, and manage the Argo CD installation itself using GitOps. The Argo-CD Autopilot is a tool which offers an opinionated way of installing Argo-CD and managing GitOps repositories"
}

func (l argocdAutopilotTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	// https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.8/argocd-autopilot-darwin-amd64.tar.gz
	url := fmt.Sprintf("%sv%s/argocd-autopilot-", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += l.os + "-" + l.arch
	case l.os == "windows" && l.arch == "amd64":
		url += l.os + "-" + l.arch + ".exe"
	case l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64":
		url += l.os + "-" + l.arch
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return argocdAutopilotTool{
		arch: arch,
		os:   os,
		GithubReleaseTool: tool.MakeGithubReleaseTool(
			"argoproj-labs", "argocd-autopilot",
		),
	}
}
