package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "koreon",
	Short: "koreon - a simple CLI to provider for openstack",
	Long: `koreon is a CLI for provider openstack
   
One can use koreon to openstack modify or openstack installing from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
