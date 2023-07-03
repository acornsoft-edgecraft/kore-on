package cmd

import (
	"fmt"
	"io/ioutil"
	"kore-on/pkg/config"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"kore-on/cmd/koreonctl/conf"

	"github.com/elastic/go-sysinfo"
	"github.com/spf13/cobra"
)

type strAddonInitCmd struct {
	verbose        bool
	osRelease      string
	osArchitecture string
	osCurrentUser  string
}

func addonInitCmd() *cobra.Command {
	addonInit := &strAddonInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get Addon configuration file",
		Long:         "This command downloads a sample file that can set addon applications.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addonInit.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&addonInit.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strAddonInitCmd) run() error {
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

	logger.Infof("Star Deployment in k8s cluster")

	if err := c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAddonInitCmd) init(workDir string) error {

	currTime := time.Now()

	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"
	addOnConfigFile := conf.AddOnConfigFile

	if !utils.CheckUserInput("Do you really want to init? \nIs this ok [y/n]: ", "y") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}

	addOnConfigFilePath, err := filepath.Abs(addOnConfigFile)
	if err != nil {
		ioutil.WriteFile(workDir+"/"+addOnConfigFile, []byte(config.AddonTemplate), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", addOnConfigFile))
	} else {
		fmt.Println("Previous " + addOnConfigFile + " file exist and it will be backup")
		os.Rename(addOnConfigFilePath, addOnConfigFilePath+"_"+currTime.Format("20060102150405"))
		ioutil.WriteFile(workDir+"/"+addOnConfigFile, []byte(config.AddonTemplate), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", addOnConfigFile))
	}
	return nil
}
