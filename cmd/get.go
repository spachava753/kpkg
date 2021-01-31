package cmd

import "github.com/spf13/cobra"

const CliVersionFlag = "version"

func MakeGet() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get or install a binary",
		Long:  `Get or install a binary. By default, the latest version of the binary will be downloaded`,
	}

	getCmd.PersistentFlags().String(CliVersionFlag, "latest", "specify a version when downloading. Default value is latest")
	return getCmd
}
