package cmd

import (
	"fmt"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"

	"kore-on/cmd/koreonctl/conf"

	"github.com/elastic/go-sysinfo"
	"github.com/spf13/cobra"
)

type strClusterUpdateCmd struct {
	dryRun         bool
	verbose        bool
	privateKey     string
	user           string
	command        string
	kubeconfig     string
	osRelease      string
	osArchitecture string
	osCurrentUser  string
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
	// 설치 directory tree check
	workDir, err := checkDirTree()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// Check installed Podman
	if err := installPodman(workDir); err != nil {
		logger.Fatal(err)
	}

	// system info
	host, err := sysinfo.Host()
	if err != nil {
		logger.Fatal(err)
	}
	currentUser, err := user.Current()
	if err != nil {
		logger.Fatal(err)
	}

	c.osCurrentUser = currentUser.Username
	c.osArchitecture = host.Info().Architecture
	c.osRelease = host.Info().OS.Platform

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

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFileName := conf.KoreOnConfigFile

	koreonToml, err := utils.GetKoreonTomlConfig(workDir + "/" + koreOnConfigFileName)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	commandArgs := []string{}

	cmdDefault := []string{
		"podman",
		"run",
		"--rm",
		"--privileged",
		"-it",
	}

	if c.osRelease == "ubuntu" && c.osCurrentUser != "root" {
		commandArgs = append(commandArgs, "sudo")
	}

	commandArgs = append(commandArgs, cmdDefault...)

	if !koreonToml.KoreOn.ClosedNetwork {
		commandArgs = append(commandArgs, "--pull")
		commandArgs = append(commandArgs, "always")
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/config", "/"+conf.KoreOnConfigDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/logs", "/"+conf.KoreOnLogsDir),
		"-mount",
		fmt.Sprintf("type=bind,source=%s,target=%s,readonly", workDir+"/archive", "/"+conf.KoreOnArchiveFileDir),
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"./" + koreonImageName,
		"update",
	}

	//- koreonctl commands
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
	//-end koreonctl commands

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsKoreonctl...)

	binary := ""
	if c.osRelease == "ubuntu" && c.osCurrentUser != "root" {
		binary, err = exec.LookPath("sudo")
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		binary, err = exec.LookPath("podman")
		if err != nil {
			logger.Fatal(err)
		}
	}

	logger.Info(commandArgs)
	err = syscall.Exec(binary, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
