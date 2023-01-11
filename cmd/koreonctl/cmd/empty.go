package cmd

import "github.com/spf13/cobra"

func emptyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
