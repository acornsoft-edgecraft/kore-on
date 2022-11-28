package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"os/exec"
	"syscall"

	"kore-on/cmd/koreonctl/conf"

	"github.com/spf13/cobra"
)

type strBstionCmd struct {
	verbose bool
}

func bastionCmd() *cobra.Command {
	bastionCmd := &strBstionCmd{}
	cmd := &cobra.Command{
		Use:          "bastion [flags]",
		Short:        "Install docker in bastion host",
		Long:         "This command a installation docker on bastion host.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bastionCmd.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&bastionCmd.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strBstionCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.bastion(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strBstionCmd) bastion(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFileName := conf.KoreOnConfigFile
	koreOnConfigFilePath := conf.KoreOnConfigFileSubDir

	koreonToml, err := utils.GetKoreonTomlConfig(workDir + "/" + koreOnConfigFileName)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	commandArgs := []string{
		"docker",
		"run",
		"--rm",
		"--privileged",
		"-it",
	}

	if !koreonToml.KoreOn.ClosedNetwork {
		commandArgs = append(commandArgs, "--pull")
		commandArgs = append(commandArgs, "always")
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir, "/"+koreOnConfigFilePath),
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"./" + koreonImageName,
		"init",
	}

	if c.verbose {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--vvv")
	}

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsKoreonctl...)

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		logger.Fatal(lookErr)
	}

	err = syscall.Exec(binary, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
