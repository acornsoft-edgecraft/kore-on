package cmd

import (
	"fmt"
	"kore-on/pkg/conf"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of koreonctl",
	Long:  `All software has versions. This is koreonctl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("koreonctl %s GitCommit:%s BuildDate:%s Platform:%s/%s\n", conf.Version, conf.CommitId, conf.BuildDate, runtime.GOOS, runtime.GOARCH)
	},
}
