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

type strAddonCmd struct {
	dryRun         bool
	verbose        bool
	privateKey     string
	user           string
	installHelm    bool
	helmBinaryFile string
	osRelease      string
	osArchitecture string
	osCurrentUser  string
}

func addonCmd() *cobra.Command {
	addon := &strAddonCmd{}

	cmd := &cobra.Command{
		Use:   "addon [flags]",
		Short: "Deployment Applications in kubernetes cluster",
		Long: "This command deploys the application to Kubernetes.\n" +
			"Use helm as the package manager for Kubernetes.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addon.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(addonInitCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&addon.verbose, "vvv", false, "verbose")
	f.BoolVarP(&addon.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&addon.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&addon.user, "user", "u", "", "login user")
	f.BoolVar(&addon.installHelm, "install-helm", false, "Helm installation options")
	f.StringVar(&addon.helmBinaryFile, "helm-binary-file", "", "helm binary file")

	return cmd
}

func (c *strAddonCmd) run() error {
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

	logger.Infof("Start deployment for k8s cluster")

	if err := c.addon(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAddonCmd) addon(workDir string) error {

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage

	addonToml, err := utils.GetAddonTomlConfig(workDir + "/" + conf.AddOnConfigFile)
	if err != nil {
		logger.Fatal(err)
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

	if !addonToml.Addon.ClosedNetwork {
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
		"addon",
	}

	//- koreonctl commands
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
