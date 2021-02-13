package opa

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type opaTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l opaTool) Name() string {
	return "opa"
}

func (l opaTool) ShortDesc() string {
	return "An open source, general-purpose policy engine"
}

func (l opaTool) LongDesc() string {
	return `The Open Policy Agent (OPA) is an open source, general-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack`
}

func (l opaTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()
	// only amd64 supported
	if l.arch != "amd64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}

	switch {
	case l.os == "darwin":
		fallthrough
	case l.os == "windows":
		fallthrough
	case l.os == "linux":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/opa_%s_amd64", l.GithubReleaseTool.MakeReleaseUrl(), version, l.os)
	if l.os == "windows" {
		url = url + ".exe"
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return opaTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("open-policy-agent", "opa", 15),
	}
}
