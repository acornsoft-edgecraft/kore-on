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

type strPrepareAirgapCmd struct {
	dryRun  bool
	verbose bool
	step    bool
}

func prepareAirgapCmd() *cobra.Command {
	airgap := &strPrepareAirgapCmd{}
	cmd := &cobra.Command{
		Use:          "prepare-airgap [flags]",
		Short:        "prepare-airgap",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgap.run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&airgap.verbose, "verbose", "v", false, "verbose")
	f.BoolVarP(&airgap.step, "step", "", false, "step")
	f.BoolVarP(&airgap.dryRun, "dry-run", "d", false, "dryRun")
	return cmd
}

func (c *strPrepareAirgapCmd) run() error {
	workDir, _ := os.Getwd()
	var err error = nil
	koreonToml, _ := utils.ValidateKoreonTomlConfig(workDir)
	startTime := time.Now()

	utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, "\nPrepare airgap ..."))
	if err = c.prepareAirgap(workDir, koreonToml); err != nil {
		return err
	}
	utils.PrintInfo(fmt.Sprintf(conf.SUCCESS_FORMAT, fmt.Sprintf("Prepare airgap Done. (%v)", (time.Duration(time.Since(startTime).Seconds())*time.Second).String())))
	return nil
}

func (c *strPrepareAirgapCmd) prepareAirgap(workDir string, koreonToml model.KoreonToml) error {
	// # 1
	utils.CheckDocker()

	utils.CopyFilePreWork(workDir, koreonToml, conf.CMD_PREPARE_AIREGAP)

	inventoryFilePath := utils.CreateInventoryFile(workDir, koreonToml, nil)

	basicFilePath := utils.CreateBasicYaml(workDir, koreonToml, conf.CMD_PREPARE_AIREGAP)

	commandArgs := []string{
		"docker",
		"run",
		"--name",
		conf.KoreonImageName,
		"--rm",
		"--privileged",
		"-it",
		"-v",
		fmt.Sprintf("%s:%s", workDir, conf.WorkDir),
		"-v",
		fmt.Sprintf("%s:%s", workDir+"/"+conf.KoreonDestDir, conf.Inventory+"/files"),
		"-v",
		fmt.Sprintf("%s:%s", inventoryFilePath, conf.InventoryIni),
		"-v",
		fmt.Sprintf("%s:%s", basicFilePath, conf.BasicYaml),
		conf.KoreonImage,
		"ansible-playbook",
		"-i",
		conf.InventoryIni,
		"-u",
		koreonToml.NodePool.Security.SSHUserID, //수정
		"--private-key",
		conf.KoreonDestDir + "/" + conf.IdRsa,
		conf.PrepareAirgapYaml,
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
