// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "kpkg", ".")
	return cmd.Run()
}

// A install step
func Install() error {
	mg.Deps(InstallDeps)
	fmt.Println("Installing...")
	cmd := exec.Command("go", "install", "./...")
	return cmd.Run()
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "list", "./...")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

// Clean up after yourself
func Clean() error {
	fmt.Println("Cleaning...")
	cmd := exec.Command("rm", "-rf", "./kpkg")
	return cmd.Run()
}

func GenGithubTool(toolName, org, repoName string) error {
	templateCode := `package {{ .pkgName }}

import (
	"fmt"
	"github.com/Masterminds/semver"
	kpkgerr "github.com/spachava753/kpkg/pkg/error"
	"github.com/spachava753/kpkg/pkg/tool"
)

type {{ .structName }} struct {
	arch,
	os string
	tool.GithubReleaseTool
}

func (l {{ .structName }}) Name() string {
	return "{{ .toolName }}"
}

func (l {{ .structName }}) ShortDesc() string {
	return "REPLACE ME"
}

func (l {{ .structName }}) LongDesc() string {
	return "REPLACE ME"
}

func (l {{ .structName }}) MakeUrl(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}
	version = v.String()

	url := fmt.Sprintf("%s%s/{{ .repoName }}", l.MakeReleaseUrl(), version)
	switch {
	case l.os == "darwin" && l.arch == "amd64":
		url += "darwin"
	case l.os == "windows" && l.arch == "amd64":
		url += ".exe"
	case l.os == "linux" && l.arch == "amd64":
	case l.os == "linux" && l.arch == "arm64":
		url += "arm64"
	default:
		return "", &kpkgerr.UnsupportedRuntimeErr{Binary: l.Name()}
	}
	return url, nil
}

func MakeBinary(os, arch string) tool.Binary {
	return k3supTool{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("{{ .org }}", "{{ .repoName }}", 20),
	}
}
`
	// package names should be lowercase
	pkgName := strings.ToLower(toolName)
	// and shouldn't contain hyphens
	pkgName = strings.ReplaceAll(pkgName, "-", "")
	structName := toolName + "Tool"
	tmpl, err := template.New("template code").Parse(templateCode)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, map[string]string{
		"pkgName":    pkgName,
		"structName": structName,
		"toolName":   toolName,
		"org":        org,
		"repoName":   repoName,
	})
	if err != nil {
		return err
	}

	return nil
}
