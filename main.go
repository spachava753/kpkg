package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/get/linkerd2"
	"os"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	hDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	root, err := config.CreatePath(hDir)
	if err != nil {
		return err
	}

	rootCmd := cmd.MakeRoot()
	getCmd := cmd.MakeGet()
	listCmd := cmd.MakeList()

	ld2 := linkerd2.MakeTool(root, runtime.GOOS, runtime.GOARCH)

	getCmd.AddCommand(
		cmd.MakeGetLinkerd2(ld2),
	)

	listCmd.AddCommand(
		cmd.MakeListLinkerd2(ld2),
	)

	rootCmd.AddCommand(getCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
