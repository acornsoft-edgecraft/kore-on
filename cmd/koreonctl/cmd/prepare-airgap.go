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

type strAirGapCmd struct {
	dryRun         bool
	verbose        bool
	privateKey     string
	image          string
	user           string
	command        string
	osRelease      string
	osArchitecture string
	osCurrentUser  string
}

func airGapCmd() *cobra.Command {
	prepareAirgap := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "Preparing a kubernetes cluster and registry for Air gap network",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepareAirgap.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(
		downLoadArchiveCmd(),
		imageUploadCmd(),
	)

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.StringVarP(&prepareAirgap.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&prepareAirgap.user, "user", "u", "", "login user")
	f.BoolVarP(&prepareAirgap.dryRun, "dry-run", "d", false, "dryRun")
	f.SortFlags = false

	return cmd
}

func downLoadArchiveCmd() *cobra.Command {
	downLoadArchive := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "download-archive [flags]",
		Short:        "Download archive files to localhost",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return downLoadArchive.run()
		},
	}

	downLoadArchive.command = "download-archive"

	f := cmd.Flags()
	f.BoolVarP(&downLoadArchive.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&downLoadArchive.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&downLoadArchive.user, "user", "u", "", "login user")

	return cmd
}

func imageUploadCmd() *cobra.Command {
	imageUpload := &strAirGapCmd{}

	cmd := &cobra.Command{
		Use:          "image-upload [flags]",
		Short:        "Images Pull and Push to private registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return imageUpload.run()
		},
	}

	imageUpload.command = "image-upload"

	f := cmd.Flags()
	f.BoolVarP(&imageUpload.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&imageUpload.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&imageUpload.user, "user", "u", "", "login user")

	return cmd
}

func (c *strAirGapCmd) run() error {
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

	logger.Infof("Start provisioning for preparing a kubernetes cluster and registry")

	if err := c.airgap(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAirGapCmd) airgap(workDir string) error {

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
		"prepare-airgap",
	}

	//- koreonctl commands
	if c.command == "image-upload" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "image-upload")
	}

	if c.command == "download-archive" {
		commandArgsKoreonctl = append(commandArgsKoreonctl, "download-archive")
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
