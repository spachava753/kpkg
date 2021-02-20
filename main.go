package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/download"
	"os"
	"runtime"
)

var (
	version,
	commit,
	goVersion string
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// fetch the users home dir
	hDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	// set up kpkg's necessary folders if not setup already (.kpkg, .kpkg/bin)
	root, err := config.CreatePath(hDir)
	if err != nil {
		return err
	}

	// create instances of top level commands
	rootCmd := cmd.MakeRoot()
	getCmd := cmd.MakeGet()
	listCmd := cmd.MakeList(root)
	rmCmd := cmd.MakeRm(root)
	versionCmd := cmd.MakeVersion(version, commit, goVersion)

	fileFetcher, err := download.InitFileFetcher()
	if err != nil {
		return err
	}

	tools := cmd.GetTools(runtime.GOOS, runtime.GOARCH)

	cmd.MakeGetBinarySubCmds(root, getCmd, tools, fileFetcher)

	cmd.MakeListBinarySubCmds(listCmd, tools)

	rootCmd.AddCommand(getCmd, listCmd, rmCmd, versionCmd)

	// set outputs
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
