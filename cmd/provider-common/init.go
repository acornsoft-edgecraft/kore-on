package common

import (
	"fmt"
	"io/ioutil"
	"kore-on/pkg/config"
	"kore-on/pkg/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type strInitCmd struct {
	verbose       bool
	check         bool
	playbookFiles []string
}

func InitCmd() *cobra.Command {
	init := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "get config file",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run()
		},
	}

	init.playbookFiles = []string{
		"./internal/playbooks/koreon-playbook/init.yaml",
	}

	f := cmd.Flags()
	f.BoolVarP(&init.verbose, "verbose", "v", false, "verbose")
	f.BoolVar(&init.check, "check", false, "check validation in config file")

	return cmd
}

func (c *strInitCmd) run() error {
	// currDir, _ := os.Getwd()
	currTime := time.Now()
	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"

	koreOnConfigFileName := viper.GetString("KoreOn.KoreOnConfigFile")
	koreOnConfigFilePath := utils.IskoreOnConfigFilePath(koreOnConfigFileName)

	// if use check flag then validation for configfile
	// var data map[string]interface{}
	if c.check {
		utils.ValidateKoreonTomlConfig(koreOnConfigFilePath, "init")
		os.Exit(0)
	}

	// create init configration file
	if utils.FileExists(koreOnConfigFilePath) {
		fmt.Println("Previous " + koreOnConfigFileName + " file exist and it will be backup")
		os.Rename(koreOnConfigFilePath, koreOnConfigFilePath+"_"+currTime.Format("20060102150405"))
	}
	ioutil.WriteFile(koreOnConfigFilePath, []byte(config.Template), 0600)
	fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFileName))

	return nil
}
