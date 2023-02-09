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
	getKubeConfigCmd := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "get-kubeconfig [flags]",
		Short:        "Get Kubeconfig file",
		Long:         "This command get kubeconfig file in k8s controlplane node.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getKubeConfigCmd.run()
		},
	}

	getKubeConfigCmd.command = "get-kubeconfig"

	f := cmd.Flags()
	f.BoolVarP(&getKubeConfigCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&getKubeConfigCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&getKubeConfigCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&getKubeConfigCmd.user, "user", "u", "", "login user")

	return cmd
}

func updateInitCmd() *cobra.Command {
	updateInitCmd := &strClusterUpdateCmd{}

	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get Installed Config file",
		Long:         "This command get installed config file in k8s controlplane node.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateInitCmd.run()
		},
	}

	updateInitCmd.command = "update-init"

	f := cmd.Flags()
	f.BoolVarP(&updateInitCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&updateInitCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&updateInitCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&updateInitCmd.user, "user", "u", "", "login user")
	f.StringVar(&updateInitCmd.kubeconfig, "kubeconfig", "", "get kubeconfig")

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
		commandArgsKoreonctl = append(commandArgsKoreonctl, c.command)
	}

	if c.kubeconfig != "" {
		key := filepath.Base(c.kubeconfig)
		keyPath, _ := filepath.Abs(c.kubeconfig)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--kubeconfig")
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
	} else {
		logger.Fatal(fmt.Errorf("[ERROR]: %s", "To run ansible-playbook an privateKey must be specified"))
	}

	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
		commandArgsKoreonctl = append(commandArgsKoreonctl, "/home/"+key)
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--private-key")
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
