package dockercompose

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"strings"
)

type composeTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l composeTool) Name() string {
	return "docker-compose"
}

func (l composeTool) ShortDesc() string {
	return "Define and run multi-container applications with Docker"
}

func (l composeTool) LongDesc() string {
	return `Define and run multi-container applications with DockerDocker Compose is a tool 
for running multi-container applications on Docker defined using the Compose file format. 
A Compose file is used to define how the one or more containers that make up your application 
are configured. Once you have a Compose file, you can create and start your application with 
a single command: docker-compose up.`
}

func (l composeTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf(
		"%s%s/docker-compose-%s-x86_64", l.MakeReleaseUrl(), version,
		strings.Title(l.os),
	)
	if l.os == "windows" {
		url += ".exe"
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return composeTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("docker", "compose"),
	}
}
