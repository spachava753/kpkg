package civo

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/thoas/go-funk"

	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type civoTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l civoTool) Name() string {
	return "civo"
}

func (l civoTool) ShortDesc() string {
	return "Civo CLI is a tool to manage your Civo.com account from the terminal"
}

func (l civoTool) LongDesc() string {
	return `Civo CLI is a tool to manage your Civo.com account from the terminal. 
The Civo web control panel has a user-friendly interface for managing your account, 
but in case you want to automate or run scripts on your account, 
or have multiple complex services, the command-line interface outlined here will be useful. `
}

func (l civoTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	c1, err := semver.NewConstraint("<= 0.6.36")
	if err != nil {
		return "", err
	}

	if c1.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c2, err := semver.NewConstraint("<= 0.6.30")
	if err != nil {
		return "", err
	}

	if c2.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c3, err := semver.NewConstraint(">= 0.6.38")
	if err != nil {
		return "", err
	}

	if c3.Check(v) && l.os == "windows" && l.arch == "386" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	c4, err := semver.NewConstraint("> 0.6.18")
	if err != nil {
		return "", err
	}

	if c4.Check(v) {
		switch {
		case l.os == "darwin",
			l.os == "linux",
			l.os == "windows":
		default:
			return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
		}
	}

	switch {
	case l.os == "darwin" && l.arch == "amd64",
		l.os == "windows" && l.arch == "amd64",
		l.os == "windows" && l.arch == "386",
		l.os == "freebsd" && l.arch == "amd64",
		l.os == "freebsd" && l.arch == "386",
		l.os == "freebsd" && l.arch == "arm",
		l.os == "openbsd" && l.arch == "386",
		l.os == "openbsd" && l.arch == "amd64",
		l.os == "solaris" && l.arch == "amd64",
		l.os == "linux" && l.arch == "amd64",
		l.os == "linux" && l.arch == "arm64",
		l.os == "linux" && l.arch == "arm",
		l.os == "linux" && l.arch == "386":
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf(
		"%sv%s/civo-%s-%s-%s", l.MakeReleaseUrl(), version, version, l.os,
		l.arch,
	)
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func (l civoTool) Versions(max uint) ([]string, error) {
	versions, err := l.GithubReleaseTool.Versions(max)
	if err != nil {
		return nil, err
	}

	// civo supports more architectures after this release
	c, err := semver.NewConstraint(">= 0.6.0")
	if err != nil {
		return nil, err
	}

	versions = funk.Filter(
		versions, func(v string) bool {
			sv := semver.MustParse(v)
			return c.Check(sv)
		},
	).([]string)

	return versions, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return civoTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("civo", "cli"),
	}
}
