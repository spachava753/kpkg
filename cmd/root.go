package cmd

import (
	"github.com/spf13/cobra"
)

func MakeRoot() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "kp",
		Short: "kp is your goto tool for managing binaries in the Kubernetes ecosystem",
		Long:  `kp is your goto tool for managing binaries in the Kubernetes ecosystem`,
	}
	return rootCmd
}
