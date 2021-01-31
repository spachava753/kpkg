package cmd

import "github.com/spf13/cobra"

func MakeRoot() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "kp",
		Short: "kp is a CLI to help install Kubernetes-related CLI's",
		Long:  `kp is a CLI to help install Kubernetes-related CLI's`,
	}
	return rootCmd
}
