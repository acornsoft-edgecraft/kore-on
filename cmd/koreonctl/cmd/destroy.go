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

type strDestroyCmd struct {
	verbose        bool
	dryRun         bool
	privateKey     string
	user           string
	command        string
	osRelease      string
	osArchitecture string
	osCurrentUser  string
}

func destroyCmd() *cobra.Command {
	destroy := &strDestroyCmd{}
	cmd := &cobra.Command{
		Use:   "destroy [flags]",
		Short: "Delete kubernetes cluster, registry, prepare-airgap",
		Long: "This command can delete [Kubernetes cluster / registry / storage].\n" +
			"* If you do not use the sub command, it is all deleted.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy.run()
		},
	}

	destroy.command = "reset-all"

	// SubCommand add
	cmd.AddCommand(
		destroyPrepareAirGapCmd(),
		destroyClusterCmd(),
		destroyRegistryCmd(),
		destroyStorageCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&destroy.verbose, "verbose", false, "verbose")
	f.BoolVarP(&destroy.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroy.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroy.user, "user", "u", "", "login user")

	return cmd
}

func destroyPrepareAirGapCmd() *cobra.Command {
	destroyPrepareAirGapCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "Destroy prepare-airgap",
		Long:         "This command deletes the registry of the prepare-airgap host and deletes related directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyPrepareAirGapCmd.run()
		},
	}

	destroyPrepareAirGapCmd.command = "reset-prepare-airgap"

	f := cmd.Flags()
	f.BoolVarP(&destroyPrepareAirGapCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyPrepareAirGapCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroyPrepareAirGapCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyPrepareAirGapCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyClusterCmd() *cobra.Command {
	destroyClusterCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "cluster [flags]",
		Short:        "Destroy cluster",
		Long:         "This command only deletes the Kubernetes cluster.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyClusterCmd.run()
		},
	}

	destroyClusterCmd.command = "reset-cluster"

	f := cmd.Flags()
	f.BoolVarP(&destroyClusterCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyClusterCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroyClusterCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyClusterCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyRegistryCmd() *cobra.Command {
	destroyRegistryCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "registry [flags]",
		Short:        "Destroy registry",
		Long:         "This command deletes the installed registry(harbor) and deletes related services and directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyRegistryCmd.run()
		},
	}

	destroyRegistryCmd.command = "reset-registry"

	f := cmd.Flags()
	f.BoolVarP(&destroyRegistryCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyRegistryCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroyRegistryCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyRegistryCmd.user, "user", "u", "", "login user")

	return cmd
}

func destroyStorageCmd() *cobra.Command {
	destroyStorageCmd := &strDestroyCmd{}

	cmd := &cobra.Command{
		Use:          "storage [flags]",
		Short:        "Destroy storage",
		Long:         "This command deletes the installed storage(NFS) and deletes related services and directories.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroyStorageCmd.run()
		},
	}

	destroyStorageCmd.command = "reset-storage"

	f := cmd.Flags()
	f.BoolVarP(&destroyStorageCmd.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroyStorageCmd.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&destroyStorageCmd.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&destroyStorageCmd.user, "user", "u", "", "login user")

	return cmd
}

func (c *strDestroyCmd) run() error {
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

	logger.Infof("Start destroy cloud infrastructure")

	if err := c.destroy(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strDestroyCmd) destroy(workDir string) error {

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
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/archive/koreon", "/"+conf.KoreOnArchiveFileDir+"/koreon"),
	}

	// podman commands
	if c.privateKey != "" {
		key := filepath.Base(c.privateKey)
		keyPath, _ := filepath.Abs(c.privateKey)
		commandArgsVol = append(commandArgsVol, "--mount")
		commandArgsVol = append(commandArgsVol, fmt.Sprintf("type=bind,source=%s,target=/home/%s,readonly", keyPath, key))
	}

	commandArgsKoreonctl := []string{
		koreOnImage,
		"./" + koreonImageName,
		"destroy",
	}

	//- koreonctl commands
	if c.command != "" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "--tags")
		commandArgsKoreonctl = append(commandArgsKoreonctl, c.command)
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
