/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"kore-on/cmd/koreonctl/conf"
	baremetal "kore-on/cmd/provider-baremetal"
	common "kore-on/cmd/provider-common"

	"kore-on/pkg/config"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	version bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "kore-on",
	Short:        "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:         `cube, It install kubernetes cluster.`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// KoreOn Version
		if version {
			logger.Info("KoreOn Version: ", conf.KoreOnVersion)
			os.Exit(1) // 다음 명령어 비활성화 후 프로그램 종료
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.CompletionOptions.HiddenDefaultCmd = true

	RootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		c.Println("Error: ", err)
		c.Println(c.UsageString())
		os.Exit(1)
		return nil
	})

	// 공용 플래그 설정
	RootCmd.Flags().BoolVar(&version, "version", false, "Show KoreOn version")

	// 하위 명령 추가
	RootCmd.AddCommand(
		common.InitCmd(),
		baremetal.CreateCmd(),
		baremetal.AddonCmd(),
		baremetal.DestroyCmd(),
		baremetal.AirGapCmd(),
		baremetal.ClusterUpdateCmd(),
		baremetal.TestCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(RootCmd)
}

func initConfig() {
	// create default logger
	err := logger.New()
	if err != nil {
		logger.Fatalf("Could not instantiate log %ss", err.Error())
	}

	// load config file
	err = config.Load()
	if err != nil {
		logger.Fatalf("Could not load configuration: %s", err.Error())
		os.Exit(1)
	}
}
