package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
)

func MakeVersion(version, commit, goVersion string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version of this release",
		Long:  `Version of this release`,
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			enc := json.NewEncoder(cmd.OutOrStdout())
			return enc.Encode(map[string]string{
				"Version":   version,
				"GitCommit": commit,
				"GoVersion": goVersion,
			})
		},
	}

	return versionCmd
}
