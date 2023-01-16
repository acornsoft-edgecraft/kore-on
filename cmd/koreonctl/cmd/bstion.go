package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"kore-on/cmd/koreonctl/conf"
	"kore-on/cmd/koreonctl/conf/templates"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

type strBstionCmd struct {
	verbose         bool
	archiveFilePath string
}

func bastionCmd() *cobra.Command {
	bastionCmd := &strBstionCmd{}
	cmd := &cobra.Command{
		Use:          "bastion [flags]",
		Short:        "Install docker in bastion host",
		Long:         "This command a installation docker on bastion host.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bastionCmd.run()
		},
	}

	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&bastionCmd.verbose, "vvv", false, "verbose")
	f.StringVar(&bastionCmd.archiveFilePath, "archive-file-path", "", "archive file path")

	return cmd
}

func (c *strBstionCmd) run() error {
	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.bastion(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strBstionCmd) bastion(workDir string) error {
	if runtime.GOOS != "linux" {
		logger.Fatal("This command option is only supported on the Linux platform.")
	}

	// Doker check
	_, dockerCheck := exec.LookPath("docker")
	if dockerCheck == nil {
		logger.Info("Docker already.")
		dockerLoad()
		os.Exit(1)
	}

	if c.archiveFilePath != "" {
		// mkdir local directory
		path := "local"
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				logger.Fatal(err)
			}
		}

		//untar gzip file
		archiveFilePath, _ := filepath.Abs(c.archiveFilePath)
		err := archiver.Unarchive(archiveFilePath, path)
		if err != nil {
			logger.Fatal(err)
		}

		// Processing template
		bastionText := template.New("bastionLocalRepoText")
		temp, err := bastionText.Parse(templates.BastionLocalRepoText)
		if err != nil {
			logger.Errorf("Template has errors. cause(%s)", err.Error())
			return err
		}

		// TODO: 진행상황을 어떻게 클라이언트에 보여줄 것인가?
		var buff bytes.Buffer
		localPath, _ := filepath.Abs(path)
		err = temp.Execute(&buff, localPath)
		if err != nil {
			logger.Errorf("Template execution failed. cause(%s)", err.Error())
			return err
		}

		repoPath := "/etc/yum.repos.d"
		err = ioutil.WriteFile(repoPath+"/bastion-local.repo", buff.Bytes(), 0644)
		if err != nil {
			logger.Fatal(err)
		}
	}
	c.dockerInstall()
	dockerLoad()

	return nil
}

func (c *strBstionCmd) dockerInstall() error {
	var commandArgs = []string{}
	if c.archiveFilePath != "" {
		if !utils.CheckUserInput("> Do you want to install docker-ce? [y/n]", "y") {
			fmt.Println("nothing to changed. exit")
			os.Exit(1)
		}
		commandArgs = []string{
			"sudo",
			"yum",
			"install",
			"-y",
			"--disablerepo=*",
			"--enablerepo=bastion-local-to-file",
			"docker-ce",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"systemctl",
			"enable",
			"docker",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"systemctl",
			"start",
			"docker",
		}
		runExecCommand(commandArgs)
	} else {
		if !utils.CheckUserInput("> Is this bastion node online network status?\n Are you sure you want to install docker-ce on this node? [y/n] ", "y") {
			fmt.Println("nothing to changed. exit")
			os.Exit(1)
		}
		commandArgs = []string{
			"sudo",
			"yum",
			"install",
			"-y",
			"yum-utils",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"yum-config-manager",
			"--add-repo",
			"https://download.docker.com/linux/centos/docker-ce.repo",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"yum",
			"install",
			"-y",
			"docker-ce",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"systemctl",
			"enable",
			"docker",
		}
		runExecCommand(commandArgs)

		commandArgs = []string{
			"sudo",
			"systemctl",
			"start",
			"docker",
		}
		runExecCommand(commandArgs)
	}

	return nil
}

func runExecCommand(commandArgs []string) error {
	commandLen := len(commandArgs)
	cmd := utils.ExecCommand(commandArgs[0], commandArgs[1:commandLen])
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError:", string(ee.Stderr))
			return fmt.Errorf("ExitError: %v", string(ee.Stderr))
		} else {
			return fmt.Errorf("err: %v", err)
		}
	}
	return nil
}

func dockerLoad() error {
	commandArgs := []string{
		"docker",
		"load",
		"--input",
		conf.KoreOnImageArchive,
	}
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
