package cmd

import "github.com/spf13/cobra"

const CliForceInstallFlag = "force"

func MakeGet() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get or install a binary",
		Long:  `Get or install a binary. By default, the latest version of the binary will be downloaded`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	getCmd.PersistentFlags().Bool(CliForceInstallFlag, false, "force a re-install if already installed")
	return getCmd
}
