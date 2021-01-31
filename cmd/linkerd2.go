package cmd

import (
	"github.com/spachava753/kpkg/pkg/get/linkerd2"
	"github.com/spf13/cobra"
	"runtime"
)

func MakeLinkerd2() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := cmd.Flags().GetString(CliVersionFlag)
			if err != nil {
				return err
			}
			return linkerd2.DownloadLinkerd2(v, runtime.GOOS, runtime.GOARCH)
		},
	}
	return cmd
}
