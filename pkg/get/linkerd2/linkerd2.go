package linkerd2

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/get"
	"io/ioutil"
	"os"
	"path"
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

func DownloadLinkerd2(version, opsys, arch string) error {
	urlConstructor := MakeUrlConstructor(version, opsys, arch)
	url, err := urlConstructor.Construct()
	if err != nil {
		return err
	}

	// create a temp file to download the CLI to
	tmpF, err := ioutil.TempFile(os.TempDir(), "linkerd2")
	if err != nil {
		return err
	}
	defer tmpF.Close()
	defer os.Remove(tmpF.Name())

	// download CLI
	err = download.FetchFile(url, tmpF)
	if err != nil {
		return err
	}

	// copy to our bin path
	hDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	// create binary file
	basePath := path.Join(hDir, ".kpkg")
	linkerd2Path := path.Join(basePath, "linkerd2", version)
	if _, err := os.Stat(linkerd2Path); os.IsNotExist(err) {
		if err := os.MkdirAll(linkerd2Path, os.ModePerm); err != nil {
			return err
		}
	}

	// copy the downloaded binary to path
	contents, err := ioutil.ReadFile(tmpF.Name())
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(linkerd2Path, "linkerd2"), contents, os.ModePerm); err != nil {
		return err
	}

	// create symlink to bin path
	return os.Symlink(path.Join(linkerd2Path, "linkerd2"), path.Join(basePath, "bin", "linkerd2"))
}
