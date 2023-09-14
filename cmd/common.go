package cmd

import "github.com/spf13/cobra"

const CliMaxVersionsInstallFlag = "max"

func InstallMaxVersionsFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Uint(
		CliMaxVersionsInstallFlag, 20,
		"max number of versions to list or search through",
	)
}
