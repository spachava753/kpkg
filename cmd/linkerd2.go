package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/get"
	"github.com/spf13/cobra"
)

func MakeGetLinkerd2(installer get.Installer) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh. The latest stable version is installed by default`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v := "latest"
			if len(args) != 0 {
				v = args[0]
			}
			f, err := cmd.Flags().GetBool(CliForceInstallFlag)
			if err != nil {
				return err
			}

			return installer.Install(v, f)
		},
	}
	return cmd
}

func MakeListLinkerd2(lister get.VersionLister) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			i, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
			if err != nil {
				return err
			}
			versions, err := lister.Versions(i)
			fmt.Println(versions)
			return err
		},
	}
	return cmd
}
