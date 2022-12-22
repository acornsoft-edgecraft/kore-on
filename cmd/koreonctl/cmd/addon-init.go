package cmd

import (
	"fmt"
	"io/ioutil"
	"kore-on/pkg/config"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"
	"path/filepath"
	"time"

	"kore-on/cmd/koreonctl/conf"

	"github.com/spf13/cobra"
)

type strAddonInitCmd struct {
	verbose bool
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

	f := cmd.Flags()
	f.BoolVar(&addonInit.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strAddonInitCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strAddonInitCmd) init(workDir string) error {
	currTime := time.Now()

	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"
	addOnConfigFile := conf.AddOnConfigFile

	if !utils.CheckUserInput("Do you really want to init? \nIs this ok [y/N]: ", "y") {
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
