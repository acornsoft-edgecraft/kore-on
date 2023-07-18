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

type strCreateCmd struct {
	dryRun         bool
	verbose        bool
	privateKey     string
	user           string
	osRelease      string
	osArchitecture string
	osCurrentUser  string
}

func createCmd() *cobra.Command {
	create := &strCreateCmd{}

	cmd := &cobra.Command{
		Use:          "create [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "This command installs the Kubernetes cluster and registry.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&create.verbose, "vvv", false, "verbose")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")
	f.StringVarP(&create.privateKey, "private-key", "p", "", "Specify ssh key path")
	f.StringVarP(&create.user, "user", "u", "", "login user")

	return cmd
}

func (c *strCreateCmd) run() error {
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

	if err := c.create(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strCreateCmd) create(workDir string) error {

	koreonImageName := conf.KoreOnImageName
	koreOnImage := conf.KoreOnImage
	koreOnConfigFileName := conf.KoreOnConfigFile

	koreonToml, err := utils.GetKoreonTomlConfig(workDir + "/config/" + koreOnConfigFileName)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	commandArgs := []string{}

	if c.osRelease == "ubuntu" && c.osCurrentUser != "root" {
		commandArgs = append(commandArgs, "sudo")
	}

	if koreonToml.KoreOn.ClosedNetwork {
		podmanLoad(workDir+"/archive/koreon/"+conf.KoreOnImageArchive, commandArgs)
	}

	cmdDefault := []string{
		"podman",
		"run",
		"--rm",
		"--privileged",
		"-it",
	}

	commandArgs = append(commandArgs, cmdDefault...)

	if !koreonToml.KoreOn.ClosedNetwork {
		commandArgs = append(commandArgs, "--pull")
		commandArgs = append(commandArgs, "always")
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/archive", "/"+conf.KoreOnArchiveFileDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/config", "/"+conf.KoreOnConfigDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/extends", "/"+conf.KoreOnExtendsFileDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/logs", "/"+conf.KoreOnLogsDir),
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
		"create",
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

	// logger.Info(commandArgs)
	err = syscall.Exec(binary, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}

func podmanLoad(koreon_img string, commandArgs []string) error {
	logger.Info("The loading of Korean images has begun in a closed network.")
	commandPodman := []string{
		"podman",
		"load",
		"--input",
		koreon_img,
	}

	commandArgs = append(commandArgs, commandPodman...)

	commandLen := len(commandArgs)
	cmd := utils.ExecCommand(commandArgs[0], commandArgs[1:commandLen])
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError:", string(ee.Stderr))
		} else {
			fmt.Println("err:", err)
		}
	}
	return nil
}
