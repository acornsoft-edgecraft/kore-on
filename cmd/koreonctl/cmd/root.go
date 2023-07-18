/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"kore-on/cmd/koreonctl/conf"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	version bool
)

// RootCmd represents the base command when called without any subcommands
var KoreOnCtlCmd = &cobra.Command{
	Use:          "koreonctl",
	Short:        "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:         `This command proceeds to automate the k8s installation task for on-premise.`,
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
	if err := KoreOnCtlCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	KoreOnCtlCmd.CompletionOptions.HiddenDefaultCmd = true

	KoreOnCtlCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		c.Println("Error: ", err)
		c.Println(c.UsageString())
		os.Exit(1)
		return nil
	})

	// 공용 플래그 설정
	KoreOnCtlCmd.Flags().BoolVar(&version, "version", false, "Show KoreOn version")

	// 하위 명령 추가
	KoreOnCtlCmd.AddCommand(
		initCmd(),
		createCmd(),
		clusterUpdateCmd(),
		destroyCmd(),
		airGapCmd(),
		bastionCmd(),
		addonCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(KoreOnCtlCmd)
}

func initConfig() {
	// create default logger
	err := logger.New()
	if err != nil {
		logger.Fatalf("Could not instantiate log %ss", err.Error())
	}

	// // load config file
	// err = config.Load()
	// if err != nil {
	// 	logger.Fatalf("Could not load configuration: %s", err.Error())
	// 	os.Exit(1)
	// }
}
