package cmd

import (
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spf13/cobra"
	"path/filepath"
)

const CliPurgeFlag = "purge"

func MakeRm(basePath string) *cobra.Command {
	var rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove or purge a binary",
		Long:  `Remove or purge a binary. You must specify a version`,
		Args: func(cmd *cobra.Command, args []string) error {
			// if purge flag is specified
			if cmd.Flags().Changed(CliPurgeFlag) {
				return cobra.ExactArgs(1)(cmd, args)
			}

			// if purge flag is not specified
			if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			purge, err := cmd.Flags().GetBool(CliPurgeFlag)
			if err != nil {
				return err
			}
			if purge {
				return tool.Purge(basePath, args[0])
			}
			return tool.RemoveVersions(basePath, args[0], args[1:])
		},
	}

	rmCmd.PersistentFlags().Bool(CliPurgeFlag, false, "purge all versions of a binary")
	return rmCmd
}
