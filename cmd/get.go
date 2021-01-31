package cmd

import (
	"github.com/spf13/cobra"
)

func MakeGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
	}
	return cmd
}
