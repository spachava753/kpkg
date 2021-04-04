package cmd

import (
	"fmt"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spf13/cobra"
	"strings"
)

func MakeGetBinarySubCmds(basePath string, parent *cobra.Command, tools []tool.Binary, f download.FileFetcher, windows bool) {
	for _, t := range tools {
		func(t tool.Binary) {
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
						p, e := tool.Install(
							basePath,
							v,
							force,
							windows,
							t,
							f,
						)
						if e != nil {
							return e
						}
						cmd.Printf("binary installed at path %s\n", p)
						return nil
					},
				},
			)
		}(t)
	}
}

func MakeListBinarySubCmds(parent *cobra.Command, tools []tool.Binary) {
	for _, t := range tools {
		func(t tool.Binary) {
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
						fmt.Println(strings.Join(versions, "\n"))
						return nil
					},
				},
			)
		}(t)
	}
}
