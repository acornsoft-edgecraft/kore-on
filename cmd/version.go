package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kore-on/pkg/conf"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of koreonctl",
	Long:  `All software has versions. This is koreonctl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("koreonctl v%s GitCommit:%s BuildDate:%s Platform:%s/%s\n", conf.Version, conf.CommitId, conf.BuildDate, runtime.GOOS, runtime.GOARCH)
	},
}
