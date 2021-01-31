package main

import (
	"fmt"
	"github.com/spachava753/kpkg/cmd"
	"os"
)

func main() {
	fmt.Println("Hello")

	rootCmd := cmd.MakeRoot()
	getCmd := cmd.MakeGet()

	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
