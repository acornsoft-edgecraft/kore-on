package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type strDestroyCmd struct {
	verbose    bool
	dryRun     bool
	privateKey string
	user       string
}

func destroyCmd() *cobra.Command {
	destroy := &strDestroyCmd{}
	cmd := &cobra.Command{
		Use:          "destroy [flags]",
		Short:        "Delete kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&destroy.verbose, "vvv", false, "verbose")
	f.BoolVarP(&destroy.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroy.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")
	f.StringVarP(&destroy.user, "user", "u", "", "SSH login user")

	return cmd
}

func (c *strDestroyCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.destroy(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strDestroyCmd) destroy(workDir string) error {
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

	if c.dryRun {
		commandArgs = append(commandArgs, "--dry-run")
	}

	if c.privateKey != "" {
		commandArgs = append(commandArgs, "--private-key")
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified"))
	}

	if c.user != "" {
		commandArgs = append(commandArgs, "--user")
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an ssh login user must be specified"))
	}

	err := syscall.Exec("/usr/local/bin/docker", commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
