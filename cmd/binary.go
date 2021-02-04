package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spf13/cobra"
)

func MakeGetBinaryCmd(usage, short, long string, binary tool.Binary) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   usage,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			v := "latest"
			if len(args) != 0 {
				v = args[0]
			}
			f, err := cmd.Flags().GetBool(CliForceInstallFlag)
			if err != nil {
				return err
			}
			p, e := binary.Install(v, f)
			if e != nil {
				return e
			}
			cmd.Printf("binary installed at path %s\n", p)
			return nil
		},
	}
	return cmd
}

func MakeListBinaryCmd(usage, short, long string, binary tool.Binary) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   usage,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
			if err != nil {
				return err
			}
			versions, err := binary.Versions()
			fmt.Println(versions)
			return err
		},
	}
	return cmd
}
