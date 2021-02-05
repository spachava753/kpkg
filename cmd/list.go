package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spf13/cobra"
)

const CliInstalledVersionsFlag = "installed"
const CliInstalledVersionsShorthandFlag = "i"

func MakeList(basePath string) *cobra.Command {
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			locallyOnly, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
			if err != nil {
				return err
			}
			if locallyOnly {
				versions, err := tool.ListInstalled(basePath, cmd.Name())
				if err != nil {
					return err
				}
				fmt.Println(versions)
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
