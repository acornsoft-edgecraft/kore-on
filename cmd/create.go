package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kore-on/pkg/conf"
	"kore-on/pkg/model"
	"kore-on/pkg/utils"
	"log"
	"os"
	"syscall"
	"time"
)

type strCreateCmd struct {
	dryRun  bool
	verbose bool
	step    bool
}

func createCmd() *cobra.Command {
	create := &strCreateCmd{}

	cmd := &cobra.Command{
		Use:          "create [flags]",
		Short:        "Install kubernetes cluster, registry",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&create.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&create.step, "step", "", false, "step")
	f.BoolVarP(&create.dryRun, "dry-run", "d", false, "dryRun")

	return cmd
}

func (c *strCreateCmd) run() error {

	//if !utils.CheckUserInput("Do you really want to create? Only 'yes' will be accepted to confirm: ", "yes") {
	//	fmt.Println("nothing to changed. exit")
	//	os.Exit(1)
	//}

	workDir, _ := os.Getwd()
	var err error = nil
	koreonToml, _ := utils.ValidateKoreonTomlConfig(workDir)
	startTime := time.Now()
	logger.Infof("Start provisioning for cloud infrastructure")

	utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, "\nSetup Koreon cluster ..."))
	if err = c.create(workDir, koreonToml); err != nil {
		return err
	}
	utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, fmt.Sprintf("Setup Koreon cluster Done. (%v)", (time.Duration(time.Since(startTime).Seconds())*time.Second).String())))

	return nil
}

func (c *strCreateCmd) create(workDir string, koreonToml model.KoreonToml) error {
	// # 1
	utils.CheckDocker()

	utils.CopyFilePreWork(workDir, koreonToml, conf.CMD_CREATE)

	inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, nil)

	basicFilePath := utils.CreateBasicYaml(workDir, koreonToml, conf.CMD_CREATE)

	commandArgs := []string{
		"docker",
		"run",
		"--name",
		conf.KoreonImageName,
		"--rm",
		"--privileged",
		"-it",
	}

	commandArgsVol := []string{
		"-v",
		fmt.Sprintf("%s:%s", workDir, conf.WorkDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/"+conf.KoreonDestDir, conf.Inventory+"/files"),
		"-v",
		fmt.Sprintf("%s:%s", inventoryFilePath, conf.InventoryIni),
		"-v",
		fmt.Sprintf("%s:%s", basicFilePath, conf.BasicYaml),
	}

	commandArgsAnsible := []string{
		conf.KoreonImage,
		"ansible-playbook",
		"-i",
		conf.InventoryIni,
		"-u",
		koreonToml.NodePool.Security.SSHUserID, //수정
		"--private-key",
		conf.KoreonDestDir + "/" + conf.IdRsa,
		conf.CreateYaml,
	}

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsAnsible...)

	if c.verbose {
		commandArgs = append(commandArgs, "-v")
	}

	if c.step {
		commandArgs = append(commandArgs, "--step")
	}

	if c.dryRun {
		commandArgs = append(commandArgs, "-C")
		commandArgs = append(commandArgs, "-D")
	}

	if koreonToml.Koreon.DebugMode {
		fmt.Printf("%s \n", commandArgs)
	}

	err := syscall.Exec(conf.DockerBin, commandArgs, os.Environ())
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return nil
}
