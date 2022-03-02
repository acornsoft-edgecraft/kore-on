package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"kore-on/pkg/conf"
	"kore-on/pkg/utils"
	"os"
	"time"
)

type strInitCmd struct {
}

func initCmd() *cobra.Command {
	create := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "get config file",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}
	return cmd
}

func (c *strInitCmd) run() error {
	currDir, err := os.Getwd()
	currTime := time.Now()

	if err != nil {
		return err
	}

	fmt.Println("Getting koreon.toml file ...")
	if utils.FileExists(currDir + "/" + conf.KoreonConfigFile) {
		fmt.Println("Previous " + conf.KoreonConfigFile + " file exist and it will be backup")
		os.Rename(conf.KoreonConfigFile, conf.KoreonConfigFile+"_"+currTime.Format("20060102150405"))
	}
	ioutil.WriteFile(currDir+"/"+conf.KoreonConfigFile, []byte(conf.Template), 0600)
	fmt.Printf(conf.SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", conf.KoreonConfigFile))

	return nil
}
