package cmd

import "github.com/spf13/cobra"

const CliInstalledVersionsFlag = "installed"

func MakeList() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List versions of a specific binary",
		Long:  `List different version candidates for installation of a specific binary`,
	}

	listCmd.PersistentFlags().Bool(CliInstalledVersionsFlag, false, "show only installed versions")
	return listCmd
}
