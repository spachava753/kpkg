package copilot

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type copilotTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l copilotTool) Name() string {
	return "copilot"
}

func (l copilotTool) ShortDesc() string {
	return "The AWS Copilot CLI is a tool for developers to build, release and operate production ready containerized applications on Amazon ECS and AWS Fargate"
}

func (l copilotTool) LongDesc() string {
	return `The AWS Copilot CLI is a tool for developers to build, release and operate production ready containerized applications on Amazon ECS and AWS Fargate. 
From getting started, pushing to a test environment, and releasing to production, Copilot helps you through the entire life of your app development`
}

func (l copilotTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%sv%s/copilot-", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "darwin"
	case l.os == "windows" && l.arch == "amd64":
		url += "windows.exe"
	case l.os == "linux" && l.arch == "amd64":
		url += "linux"
	case l.os == "linux" && l.arch == "arm64":
		url += "linux-arm64"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return copilotTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("aws", "copilot-cli"),
	}
}
