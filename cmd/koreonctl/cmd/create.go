package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"strings"
	"syscall"

	"kore-on/cmd/koreonctl/conf"

	"github.com/spf13/cobra"
)

type strCreateCmd struct {
	dryRun     bool
	verbose    bool
	privateKey string
	user       string
}

func createCmd() *cobra.Command {
	create := &strCreateCmd{}

	cmd := &cobra.Command{
		Use:          "create [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&create.verbose, "vvv", false, "verbose")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&create.privateKey, "private-key", "p", "", "Specify ansible playbook privateKey")
	f.StringVarP(&create.user, "user", "u", "", "SSH login user")

	return cmd
}

func (c *strCreateCmd) run() error {

	//if !utils.CheckUserInput("Do you really want to create? Only 'yes' will be accepted to confirm: ", "yes") {
	//	fmt.Println("nothing to changed. exit")
	//	os.Exit(1)
	//}

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.create(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strCreateCmd) create(workDir string) error {
	// Doker check
	utils.CheckDocker()

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFilePath := conf.KoreOnConfigFileSubDir

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
		fmt.Sprintf("%s:%s", workDir, "/"+koreOnConfigFilePath),
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"./" + koreonImageName,
		"create",
	}

	if c.privateKey != "" {
		key := strings.Split(c.privateKey, "/")
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", c.privateKey, key[len(key)-1]))
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
		key := strings.Split(c.privateKey, "/")
		commandArgs = append(commandArgs, "/home/"+key[len(key)-1])
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified"))
	}

	if c.user != "" {
		commandArgs = append(commandArgs, "--user")
		commandArgs = append(commandArgs, c.user)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an ssh login user must be specified"))
	}

	err := syscall.Exec("/usr/local/bin/docker", commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
