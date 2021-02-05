package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
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
	listCmd := cmd.MakeList()
	rmCmd := cmd.MakeRm(root)

	// create a file fetcher for binaries to fetch file
	fileFetcher, err := download.MakeFileFetcherTempDir(&http.Client{
		Timeout: time.Second * 5,
	})
	if err != nil {
		return err
	}

	ld2 := linkerd2.MakeBinary(root, runtime.GOOS, runtime.GOARCH, fileFetcher)

	getCmd.AddCommand(
		cmd.MakeGetBinaryCmd(
			cmd.Linkerd2Usage,
			cmd.Linkerd2Short,
			cmd.Linkerd2Long,
			ld2,
		),
	)

	listCmd.AddCommand(
		cmd.MakeListBinaryCmd(
			cmd.Linkerd2Usage,
			cmd.Linkerd2Short,
			cmd.Linkerd2Long,
			ld2,
		),
	)

	rootCmd.AddCommand(getCmd, listCmd, rmCmd)

	// set outputs
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
