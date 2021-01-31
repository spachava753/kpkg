package linkerd2

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/get"
)

type linkerd2UrlConstructor struct {
	version,
	arch,
	os string
}

func (l linkerd2UrlConstructor) Construct() (string, error) {
	if l.version == "latest" {
		l.version = "21.1.4"
	}
	switch l.os {
	case "darwin":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/edge-%s/linkerd2-cli-edge-%s-darwin", l.version, l.version), nil
	case "windows":
		return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/edge-%s/linkerd2-cli-edge-%s-windows.exe", l.version, l.version), nil
	case "linux":
		switch l.arch {
		case "amd64":
			fallthrough
		case "arm":
			fallthrough
		case "arm64":
			return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/edge-%s/linkerd2-cli-edge-%s-linux-%s", l.version, l.version, l.arch), nil
		default:
			return "", fmt.Errorf("unsupported architecture: %s", l.arch)
		}
	}
	return "", fmt.Errorf("unsupported os: %s", l.os)
}

func MakeUrlConstructor(version, os, arch string) get.UrlConstructor {
	return linkerd2UrlConstructor{
		version: version,
		os:      os,
		arch:    arch,
	}
}
