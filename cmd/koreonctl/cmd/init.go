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

type strInitCmd struct {
	verbose bool
}

func initCmd() *cobra.Command {
	init := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get configuration file",
		Long:         "This command downloads a sample file that can set machine information and installation information.",
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

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFilePath := conf.KoreOnConfigFileSubDir

	commandArgs := []string{
		"docker",
		"run",
		"--pull",
		"always",
		"--rm",
		"--privileged",
		"-it",
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

	err := syscall.Exec(binary, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
