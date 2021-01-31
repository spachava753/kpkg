package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/cmd"
	"github.com/spachava753/kpkg/pkg/config"
	"os"
)

func main() {
	hDir, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
	config.CreateBinPath(hDir)

	rootCmd := cmd.MakeRoot()
	getCmd := cmd.MakeGet()
	listCmd := cmd.MakeList()

	getCmd.AddCommand(
		cmd.MakeGetLinkerd2(),
	)

	listCmd.AddCommand(
		cmd.MakeListLinkerd2(),
	)

	rootCmd.AddCommand(getCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
