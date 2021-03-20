// +build mage

package main

import (
	"bytes"
	"fmt"
	"github.com/spachava753/kpkg/pkg/util"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
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
	return {{ .structName }}{
		arch:              arch,
		os:                os,
		GithubReleaseTool: tool.MakeGithubReleaseTool("{{ .org }}", "{{ .repoName }}", 20),
	}
}
`
	// package names should be lowercase
	pkgName := strings.ToLower(toolName)
	// and shouldn't contain hyphens. This will also be the folder name
	pkgName = strings.ReplaceAll(pkgName, "-", "")

	fmt.Printf("pkg name will be %s\n", pkgName)

	// make sure that the folder doesn't exist
	folderPath := filepath.Join("pkg", "tool", pkgName)
	var folderExists bool
	info, err := os.Stat(folderPath)
	if err == nil {
		if !info.IsDir() {
			fmt.Errorf("found a file with the conflicting name %s, delete this file first", pkgName)
		}
		// since the dir already exists, make sure to check that we can create the file
		var files []string
		filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if files != nil && util.ContainsString(files, pkgName) {
			fmt.Errorf("filepath %s already exists", filepath.Join(folderPath, pkgName))
		}
		folderExists = true
	}

	structName := ToLowerCamel(toolName) + "Tool"
	tmpl, err := template.New("template code").Parse(templateCode)
	if err != nil {
		return err
	}
	var bytes bytes.Buffer
	err = tmpl.Execute(&bytes, map[string]string{
		"pkgName":    pkgName,
		"structName": structName,
		"toolName":   toolName,
		"org":        org,
		"repoName":   repoName,
	})
	if err != nil {
		return err
	}

	// create a new folder if it doesn't exist
	if !folderExists {
		fmt.Printf("creating folder %s\n", folderPath)
		if err := os.Mkdir(folderPath, 0775); err != nil {
			return err
		}
	} else {
		fmt.Printf("folder exists %s, skipping creation\n", folderPath)
	}

	// create the file
	codePath := filepath.Join(folderPath, pkgName+".go")
	fmt.Printf("creating file at location %s\n", codePath)
	f, err := os.Create(codePath)
	if err != nil {
		return err
	}

	// and dump the code
	fmt.Printf("dumping template code in %s\n", codePath)
	_, err = io.Copy(f, &bytes)

	return err
}

// COPIED BELOW CODE FROM https://github.com/iancoleman/strcase

var uppercaseAcronym = map[string]string{
	"ID": "id",
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if a, ok := uppercaseAcronym[s]; ok {
		s = a
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	return toCamelInitCase(s, false)
}
