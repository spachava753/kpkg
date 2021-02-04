package cmd

import "github.com/spf13/cobra"

const CliPurgeFlag = "purge"

func MakeRm() *cobra.Command {
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
	}

	rmCmd.PersistentFlags().Bool(CliPurgeFlag, false, "purge all versions of a binary")
	return rmCmd
}
