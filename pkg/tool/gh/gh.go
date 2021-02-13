package gh

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ghTool struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l ghTool) Extract(artifactPath, _ string) (string, error) {
	// gh releases contain the binary and some examples. Pick only the binary

	// expect given path to be to contain bin folder
	dirsInfo, err := ioutil.ReadDir(artifactPath)
	if err != nil {
		return "", err
	}
	// expect to find exactly on dir
	if len(dirsInfo) != 1 {
		return "", fmt.Errorf("found too many subdirectories at %s", artifactPath)
	}
	binDirPath := filepath.Join(artifactPath, dirsInfo[0].Name(), "bin")
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

func (l ghTool) Name() string {
	return "gh"
}

func (l ghTool) ShortDesc() string {
	return "GitHub’s official command line tool"
}

func (l ghTool) LongDesc() string {
	return `gh is GitHub on the command line. It brings pull requests, issues, 
and other GitHub concepts to the terminal next to where you are already working with git and your code.`
}

func (l ghTool) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	c, err := semver.NewConstraint("< 0.5.6")
	if err != nil {
		return "", err
	}
	// linux/arm64 is not supported until 0.5.6
	if c.Check(v) && l.os == "linux" && l.arch == "arm64" {
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	version = v.String()

	switch {
	case l.os == "darwin" && l.arch == "amd64":
		l.os = "macOS"
		break
	case l.os == "windows" && l.arch == "amd64":
		fallthrough
	case l.os == "windows" && l.arch == "386":
		fallthrough
	case l.os == "linux" && l.arch == "amd64":
		fallthrough
	case l.os == "linux" && l.arch == "arm64":
		fallthrough
	case l.os == "linux" && l.arch == "386":
		break
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	url := fmt.Sprintf("%sv%s/gh_%s_%s_%s", l.MakeReleaseUrl(), version, version, l.os, l.arch)
	if l.os == "windows" {
		return url + ".zip", nil
	}
	return url + ".tar.gz", nil
}

func MakeBinary(os, arch string) tool.Binary {
	return ghTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("cli", "cli", 20),
	}
}
