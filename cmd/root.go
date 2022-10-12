/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	baremetal "cube/cmd/provider-baremetal"
	common "cube/cmd/provider-common"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cube",
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
	RootCmd.AddCommand(
		common.InitCmd(),
		baremetal.CreateCmd(),
		baremetal.ApplyCmd(),
		baremetal.DestroyCmd(),
		baremetal.TestCmd(),
	)
}
