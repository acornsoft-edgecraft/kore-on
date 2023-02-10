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

type strClusterUpdateCmd struct {
	dryRun     bool
	verbose    bool
	privateKey string
	user       string
	command    string
	kubeconfig string
}

func clusterUpdateCmd() *cobra.Command {
	clusterUpdate := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "update [flags]",
		Short:        "Update kubernetes cluster(node scale in/out)",
		Long:         "This command update the Kubernetes cluster nodes (node scale in/out)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return clusterUpdate.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(
		getKubeConfigCmd(),
		updateInitCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&clusterUpdate.verbose, "vvv", false, "verbose")
	f.BoolVarP(&clusterUpdate.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&clusterUpdate.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&clusterUpdate.user, "user", "u", "", "login user")
	f.StringVar(&clusterUpdate.kubeconfig, "kubeconfig", "", "get kubeconfig")

	return cmd
}

func (c *strClusterUpdateCmd) run() error {
	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.clusterUpdate(workDir); err != nil {
		return err
	}
	return nil
}

func getKubeConfigCmd() *cobra.Command {
	getKubeConfig := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "get-kubeconfig [flags]",
		Short:        "Get Kubeconfig file",
		Long:         "This command get kubeconfig file in k8s controlplane node.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getKubeConfig.run()
		},
	}

	getKubeConfig.command = "get-kubeconfig"

	f := cmd.Flags()
	f.BoolVarP(&getKubeConfig.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&getKubeConfig.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&getKubeConfig.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&getKubeConfig.user, "user", "u", "", "login user")

	return cmd
}

func updateInitCmd() *cobra.Command {
	updateInit := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get Installed Config file",
		Long:         "This command get installed config file in k8s controlplane node.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateInit.run()
		},
	}

	updateInit.command = "update-init"

	f := cmd.Flags()
	f.BoolVarP(&updateInit.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&updateInit.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&updateInit.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&updateInit.user, "user", "u", "", "login user")
	f.StringVar(&updateInit.kubeconfig, "kubeconfig", "", "get kubeconfig")

	return cmd
}

func (c *strClusterUpdateCmd) clusterUpdate(workDir string) error {
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
		"update",
	}

	// sub command
	if c.command != "" {
		if c.command == "update-init" {
			c.command = "init"
		}
		commandArgsKoreonctl = append(commandArgsKoreonctl, c.command)
	}

	if c.command != "update-init" && c.command != "get-kubeconfig" && c.kubeconfig != "" {
		key := filepath.Base(c.kubeconfig)
		keyPath, _ := filepath.Abs(c.kubeconfig)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--kubeconfig")
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
	}

	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--private-key")
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run this ansible-playbook an kubeconfig option must be specified.\n You can get kubeconfig with 'get-kubeconfig' command"))
	}

	if c.verbose {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--verbose")
	}

	if c.dryRun {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--dry-run")
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
