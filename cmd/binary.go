package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spf13/cobra"
)

func MakeGetBinarySubCmds(basePath string, parent *cobra.Command, tools []tool.Binary, f download.FileFetcher) {
	for _, t := range tools {
		parent.AddCommand(
			&cobra.Command{
				Use:   t.Name(),
				Short: t.ShortDesc(),
				Long:  t.LongDesc(),
				RunE: func(cmd *cobra.Command, args []string) error {
					v := "latest"
					if len(args) != 0 {
						v = args[0]
					}
					force, err := cmd.Flags().GetBool(CliForceInstallFlag)
					if err != nil {
						return err
					}
					p, e := tool.Install(basePath, v, force, t, f)
					if e != nil {
						return e
					}
					cmd.Printf("binary installed at path %s\n", p)
					return nil
				},
			},
		)
	}
}

func MakeListBinarySubCmds(basePath string, parent *cobra.Command, tools []tool.Binary, f download.FileFetcher) {
	for _, t := range tools {
		parent.AddCommand(
			&cobra.Command{
				Use:   t.Name(),
				Short: t.ShortDesc(),
				Long:  t.LongDesc(),
				RunE: func(cmd *cobra.Command, args []string) error {
					locallyOnly, err := cmd.Flags().GetBool(CliInstalledVersionsFlag)
					if err != nil {
						return err
					}
					if locallyOnly {
						return nil
					}
					versions, err := t.Versions()
					if err != nil {
						return err
					}
					fmt.Println(versions)
					return nil
				},
			},
		)
	}
}
