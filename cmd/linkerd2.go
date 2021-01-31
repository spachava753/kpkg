package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/get/linkerd2"
	"github.com/spf13/cobra"
	"runtime"
)

func MakeGetLinkerd2() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh. The latest stable version is installed by default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := cmd.Flags().GetString(CliVersionFlag)
			if err != nil {
				return err
			}
			return linkerd2.Download(v, runtime.GOOS, runtime.GOARCH)
		},
	}
	return cmd
}

func MakeListLinkerd2() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
			if err != nil {
				return err
			}
			versions, err := linkerd2.Versions(v)
			fmt.Println(versions)
			return err
		},
	}
	return cmd
}
