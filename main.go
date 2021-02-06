package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/download"
	"net/http"
	"os"
	"runtime"
	"time"
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

	// create a file fetcher for binaries to fetch file
	fileFetcher, err := download.MakeFileFetcherTempDir(&http.Client{
		Timeout: time.Second * 10,
	})
	if err != nil {
		return err
	}
	fileFetcher, err = download.MakeRetryFileFetcher(3, os.Stdout, fileFetcher)
	if err != nil {
		return err
	}

	tools := cmd.GetTools(root, runtime.GOOS, runtime.GOARCH, fileFetcher)

	cmd.MakeGetBinarySubCmds(root, getCmd, tools, fileFetcher)

	cmd.MakeListBinarySubCmds(root, listCmd, tools, fileFetcher)

	rootCmd.AddCommand(getCmd, listCmd, rmCmd)

	// set outputs
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
