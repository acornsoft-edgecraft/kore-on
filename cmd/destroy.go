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

type strDestroyCmd struct {
	name    string
	dryRun  bool
	timeout int64
	target  string
	verbose bool
	step    bool
}

func destroyCmd() *cobra.Command {
	destroy := &strDestroyCmd{}
	cmd := &cobra.Command{
		Use:          "destroy [flags]",
		Short:        "destroy",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return destroy.run()
		},
	}
	f := cmd.Flags()
	f.StringVarP(&destroy.target, "target", "", "", "target. [all]")
	f.BoolVarP(&destroy.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&destroy.step, "step", "", false, "step")
	f.BoolVarP(&destroy.dryRun, "dry-run", "d", false, "dryRun")
	return cmd
}

func (c *strDestroyCmd) run() error {
	workDir, _ := os.Getwd()
	var err error = nil
	koreonToml, _ := utils.ValidateKoreonTomlConfig(workDir)
	startTime := time.Now()
	logger.Infof("Start provisioning for koreon infrastructure")

	switch c.target {
	default:
		utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, "\nDestroy koreon cluster ..."))
		if err = c.destroy(workDir, koreonToml); err != nil {
			return err
		}
		utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, fmt.Sprintf("Setup Koreon cluster Done. (%v)", (time.Duration(time.Since(startTime).Seconds())*time.Second).String())))
	}

	//infra.PrintK8sWorkResult(workDir, c.target)
	utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, "Installation Completed."))
	return nil
}

func (c *strDestroyCmd) destroy(workDir string, koreonToml model.KoreonToml) error {

	if !utils.CheckUserInput("Do you really want to destroy? Only 'yes' will be accepted to confirm: ", "yes") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}

	// # 1
	utils.CheckDocker()

	utils.CopyFilePreWork(workDir, koreonToml, conf.CMD_DESTROY)

	inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, nil)

	basicFilePath := utils.CreateBasicYaml(workDir, koreonToml, conf.CMD_DESTROY)

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
	}

	commandArgs = append(commandArgs, commandArgsVol...)
	commandArgs = append(commandArgs, commandArgsAnsible...)

	switch c.target {
	case "all":
		commandArgs = append(commandArgs, conf.ResetYaml)
	default:
		commandArgs = append(commandArgs, conf.ResetYaml)
		commandArgs = append(commandArgs, "--tags")
		commandArgs = append(commandArgs, "reset-cluster")
	}

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
