package utils

import (
	"bytes"
	"fmt"
	"kore-on/pkg/model"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func CheckKoreonToml(value *model.KoreOnToml) error {

	return nil
}

func CheckCommand(cmd *cobra.Command) error {
	cmdCheck := cmd.Commands()

	if os.Args[1] == cmd.Name() {
		if len(os.Args[1:]) < 2 {
			args := append([]string{"koreonctl"}, os.Args[1:]...)
			buf := new(bytes.Buffer)
			cmd.SetErr(buf)
			fmt.Println("Error: unknown command", args)
			fmt.Println(fmt.Sprintf("Run 'koreonctl %s --help' for usage.", cmd.Name()))
			os.Exit(1)
		}
		subcmd := os.Args[2]
		for _, cv := range cmdCheck {
			if cv.Name() != subcmd && string(subcmd[0]) != "-" {
				strContains := ""
				errMessage := ""
				for _, v := range cmdCheck {
					if strings.Contains(v.Name(), subcmd) {
						strContains = v.Name()
						break
					}
				}
				args := append([]string{"koreonctl"}, os.Args[1:]...)
				buf := new(bytes.Buffer)
				cmd.SetErr(buf)
				fmt.Println("Error: unknown command", args)
				if strContains != "" {
					errMessage = fmt.Sprintf("Did you mean this?\n\t%s\n\nRun 'koreonctl %s --help' for usage.", strContains, cmd.Name())
				} else {
					errMessage = fmt.Sprintf("Run 'koreonctl %s --help' for usage.", cmd.Name())
				}
				fmt.Println(errMessage)
				os.Exit(1)
			}
		}
	}

	return nil
}
