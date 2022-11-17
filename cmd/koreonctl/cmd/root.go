/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"kore-on/pkg/config"
	"kore-on/pkg/logger"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var KoreOnCtlCmd = &cobra.Command{
	Use:   "koreonctl",
	Short: "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:  `cube, It install kubernetes cluster.`,
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

	KoreOnCtlCmd.AddCommand(
		initCmd(),
		createCmd(),
		destroyCmd(),
		airGapCmd(),
	)

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
