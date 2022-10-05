package common

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type strInitCmd struct {
}

func InitCmd() *cobra.Command {
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

	fmt.Println("Getting Strart time: ", currTime)
	fmt.Println("currDir: ", currDir)
	return nil
}
