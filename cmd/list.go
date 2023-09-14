package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spachava753/kpkg/pkg/tool"
)

const CliInstalledVersionsFlag = "installed"
const CliInstalledVersionsShorthandFlag = "i"

func MakeList(basePath string) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List versions of a specific binary",
		Long:  `List different version candidates for installation of a specific binary`,
		Example: `
Show installed tools:
kpkg list -i

Show versions of a specific binary:
kpkg list eksctl

Show installed versions of a specific binary:
kpkg list -i eksctl
`,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			locallyOnly, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
			if err != nil {
				return err
			}

			if locallyOnly {
				binaries, err := tool.ListInstalled(basePath)
				if err != nil {
					return err
				}
				fmt.Println(strings.Join(binaries, "\n"))
				return nil
			}

			cmd.Help()
			return nil
		},
	}

	listCmd.PersistentFlags().BoolP(
		CliInstalledVersionsFlag,
		CliInstalledVersionsShorthandFlag,
		false,
		"show only installed versions",
	)
	InstallMaxVersionsFlag(listCmd)
	return listCmd
}
