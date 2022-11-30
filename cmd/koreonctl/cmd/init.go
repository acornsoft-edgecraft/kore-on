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

type strInitCmd struct {
	verbose bool
}

func initCmd() *cobra.Command {
	init := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get configuration file",
		Long:         "This command downloads a sample file that can set machine information and installation information.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&init.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strInitCmd) run() error {

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strInitCmd) init(workDir string) error {
	// Doker check
	utils.CheckDocker()

	currTime := time.Now()

	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"
	koreOnConfigFile := conf.KoreOnConfigFile

	if !utils.CheckUserInput("Do you really want to init? \nIs this ok [y/N]: ", "y") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}

	koreOnConfigFilePath, err := filepath.Abs(koreOnConfigFile)
	if err != nil {
		ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))
	} else {
		fmt.Println("Previous " + koreOnConfigFile + " file exist and it will be backup")
		os.Rename(koreOnConfigFilePath, koreOnConfigFilePath+"_"+currTime.Format("20060102150405"))
		ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))

	}
	return nil
}
