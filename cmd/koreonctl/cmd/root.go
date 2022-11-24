/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"kore-on/pkg/logger"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var KoreOnCtlCmd = &cobra.Command{
	Use:   "koreonctl",
	Short: "Install kubernetes cluster to on-premise system with registry and storage system",
	Long:  `This command proceeds to automate the k8s installation task for on-premise.`,
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
}
