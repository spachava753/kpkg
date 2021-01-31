package main

import (
	"fmt"
	"github.com/spachava753/kpkg/cmd"
)

func main() {
	fmt.Println("Hello")

	// wire up CLI commands

	rootCmd := cmd.MakeRoot()
	getCmd := cmd.MakeGet()

	rootCmd.AddCommand(getCmd)

}
