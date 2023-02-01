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
	"time"

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
		os_release := runExecCommand([]string{"grep -ri 'ID=' /etc/os-release | awk 'NR==1 {print $0}' | cut -f2 -d= | sed 's/\"//g'"})
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
		var bastionTemp string
		var repoPath string
		if os_release == "ubuntu" {
			var theTime = time.Now().Format("20060102150405")
			//Backup apt repository
			commandArgs := []string{
				"sudo",
				"mv",
				"/etc/apt/sources.list.d",
				"/etc/apt/sources.list.d-" + theTime,
			}
			runExecCommand(commandArgs)

			//Backup apt repository
			commandArgs = []string{
				"sudo",
				"mkdir",
				"/etc/apt/sources.list.d",
			}
			runExecCommand(commandArgs)

			//Backup apt repository
			commandArgs = []string{
				"sudo",
				"cp",
				"/etc/apt/sources.list",
				"/etc/apt/sources.list-" + theTime,
			}
			runExecCommand(commandArgs)

			//Replace apt repository
			commandArgs = []string{
				"sudo",
				"sed 's/^deb/#deb/g' /etc/apt/sources.list",
			}
			runExecCommand(commandArgs)

			bastionTemp = templates.UbuntuBastionLocalRepoText
			repoPath = "/etc/apt/sources.list.d/bastion-local-to-file.list"
		} else if os_release == "centos" || os_release == "rhel" {
			bastionTemp = templates.BastionLocalRepoText
			repoPath = "/etc/yum.repos.d/bastion-local.repo"
		} else {
			logger.Fatal("This command option is only supported on the Linux platform(Centos, RedHat, Ubuntu).")
		}

		temp, err := bastionText.Parse(bastionTemp)
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

		err = ioutil.WriteFile(repoPath, buff.Bytes(), 0644)
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
	commandArgs = []string{
		"sh",
		"-c",
		"sudo",
		"grep -ri 'ID=' /etc/os-release | awk 'NR==1 {print $0}' | cut -f2 -d= | sed 's/\"//g'",
	}
	os_release := runExecCommand(commandArgs)
	if c.archiveFilePath != "" {
		if !utils.CheckUserInput("> Do you want to install docker-ce? [y/n] ", "y") {
			fmt.Println("nothing to changed. exit")
			os.Exit(1)
		}
		if os_release == "ubuntu" {
			//docker install
			commandArgs = []string{
				"sudo",
				"apt-get",
				"install",
				"-y",
				"docker-ce",
			}
			runExecCommand(commandArgs)
		} else if os_release == "centos" || os_release == "rhel" {
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
		} else {
			logger.Fatal("This command option is only supported on the Linux platform(CentOS, RedHat, Ubuntu).")
		}

		// Calling Sleep method
		time.Sleep(5 * time.Second)
		dockerRestart()
	} else {
		if !utils.CheckUserInput("> Is this bastion node online network status?\n Are you sure you want to install docker-ce on this node? [y/n] ", "y") {
			fmt.Println("nothing to changed. exit")
			os.Exit(1)
		}
		if os_release == "ubuntu" {
			commandArgs = []string{
				"sudo",
				"mkdir",
				"-p",
				"/etc/apt/keyrings",
			}
			runExecCommand(commandArgs)
			commandArgs = []string{
				"sudo",
				"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg",
			}
			runExecCommand(commandArgs)
			commandArgs = []string{
				"sudo",
				"echo 'deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu " +
					"$(lsb_release -cs) stable' | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null",
			}
			runExecCommand(commandArgs)
			commandArgs = []string{
				"sudo",
				"apt-get update",
			}
			runExecCommand(commandArgs)
			commandArgs = []string{
				"sudo",
				"apt-get install -y docker-ce",
			}
			runExecCommand(commandArgs)

		} else if os_release == "centos" || os_release == "rhel" {
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
		} else {
			logger.Fatal("This command option is only supported on the Linux platform(CentOS, RedHat, Ubuntu).")
		}

		// Calling Sleep method
		time.Sleep(5 * time.Second)
		dockerRestart()
	}

	return nil
}

func runExecCommand(commandArgs []string) string {
	commandLen := len(commandArgs)
	cmd := utils.ExecCommand(commandArgs[0], commandArgs[1:commandLen])
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			logger.Fatal("ExitError:", string(ee.Stderr))
			// return fmt.Errorf("ExitError: %v", string(ee.Stderr))
		} else {
			logger.Fatal("err: %v", err)
			// return "", fmt.Errorf("err: %v", err)
		}
	}
	return string(out)
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

func dockerRestart() {
	var commandArgs = []string{}
	commandArgs = []string{
		"sudo",
		"systemctl",
		"daemon-reload",
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
