package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type strInitCmd struct {
	verbose bool
}

func initCmd() *cobra.Command {
	init := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "get config file",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&init.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strInitCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strInitCmd) init(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := viper.GetString("KoreOn.KoreonImageName")
	koreOnImage := viper.GetString("KoreOn.KoreOnImage")
	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)

	commandArgs := []string{
		"docker",
		"run",
		"--name",
		koreonImageName,
		"--rm",
		"--privileged",
		"-it",
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir, koreOnConfigFilePath),
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"destroy",
	}

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsKoreonctl...)

	if c.verbose {
		commandArgs = append(commandArgs, "--vvv")
	}

	fmt.Println(commandArgs)

	// err := syscall.Exec("/usr/local/bin/docker", commandArgs, os.Environ())
	// if err != nil {
	// 	log.Printf("Command finished with error: %v", err)
	// }

	return nil
}
