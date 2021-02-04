package cmd

import "github.com/spf13/cobra"

const CliInstalledVersionsFlag = "installed"
const CliInstalledVersionsShorthandFlag = "i"

func MakeList() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List versions of a specific binary",
		Long:  `List different version candidates for installation of a specific binary`,
		Args: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed(CliInstalledVersionsFlag) {
				return cobra.ExactArgs(1)(cmd, args)
			}
			return nil
		},
	}

	listCmd.PersistentFlags().BoolP(
		CliInstalledVersionsFlag,
		CliInstalledVersionsShorthandFlag,
		false,
		"show only installed versions",
	)
	return listCmd
}
