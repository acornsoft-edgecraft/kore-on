package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"kore-on/cmd/koreonctl/conf"

	"github.com/spf13/cobra"
)

type strClusterNodeCmd struct {
	dryRun     bool
	verbose    bool
	privateKey string
	user       string
	command    string
}

func clusterNodeCmd() *cobra.Command {
	clusterNode := &strClusterNodeCmd{}

	cmd := &cobra.Command{
		Use:          "node [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "This command installs the Kubernetes cluster and registry.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clusterNode.run()
		},
	}

	cmd.AddCommand(
		nodeAdd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&clusterNode.verbose, "vvv", false, "verbose")
	f.BoolVarP(&clusterNode.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&clusterNode.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&clusterNode.user, "user", "u", "", "login user")

	return cmd
}

func nodeAdd() *cobra.Command {
	nodeAdd := &strClusterNodeCmd{}

	cmd := &cobra.Command{
		Use:          "add [flags]",
		Short:        "Kubernetes update for node sale in",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodeAdd.run()
		},
	}

	nodeAdd.command = "image-upload"

	f := cmd.Flags()
	f.BoolVarP(&nodeAdd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&nodeAdd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&nodeAdd.user, "user", "u", "", "login user")

	return cmd
}

func (c *strClusterNodeCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.node(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strClusterNodeCmd) node(workDir string) error {
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
		"create",
	}

	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
	}

	if c.verbose {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--verbose")
	}

	if c.dryRun {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--dry-run")
	}

	if c.privateKey != "" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--private-key")
		key := filepath.Base(c.privateKey)
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified"))
	}

	if c.user != "" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--user")
		commandArgsKoreonctl = append(commandArgsKoreonctl, c.user)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an ssh login user must be specified"))
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
