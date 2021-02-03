package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
	"os"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
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
	rmCmd := cmd.MakeRm()

	ld2 := linkerd2.MakeBinary(root, runtime.GOOS, runtime.GOARCH)

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

	rmCmd.AddCommand(
		cmd.MakeRmBinaryCmd(
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
