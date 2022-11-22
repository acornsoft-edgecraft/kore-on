/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	baremetal "kore-on/cmd/provider-baremetal"
	common "kore-on/cmd/provider-common"

	"kore-on/pkg/config"
	"kore-on/pkg/logger"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kore-on",
	Short: "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:  `cube, It install kubernetes cluster.`,
}

// RootCmd represents the base command when called without any subcommands
var KoreOnCtlCmd = &cobra.Command{
	Use:   "koreonctl",
	Short: "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:  `cube, It install kubernetes cluster.`,
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

	RootCmd.AddCommand(
		common.InitCmd(),
		baremetal.CreateCmd(),
		baremetal.ApplyCmd(),
		baremetal.DestroyCmd(),
		baremetal.AirGapCmd(),
		baremetal.TestCmd(),
	)
	KoreOnCtlCmd.AddCommand()

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
