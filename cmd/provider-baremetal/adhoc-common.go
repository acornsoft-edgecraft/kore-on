package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/execute/measure"
	"github.com/apenella/go-ansible/pkg/options"
)

// Check Helm Repo Login - ansible ad-hoc used
func checkHelmRepoLogin(id string, pw string, commandArgs string) error {
	var err error

	buff := new(bytes.Buffer)

	ansibleConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
	}

	executorTimeMeasurement := measure.NewExecutorTimeMeasurement(
		execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
	)

	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		Inventory:  " 127.0.0.1,",
		ModuleName: "command",
		Args:       commandArgs,
	}

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern:           "all",
		Exec:              executorTimeMeasurement,
		Options:           ansibleAdhocOptions,
		ConnectionOptions: ansibleConnectionOptions,
		StdoutCallback:    "json",
	}

	err = adhoc.Run(context.TODO())
	if err != nil {
		firstIndex := strings.Index(buff.String(), "stderr")
		lastIndex := strings.LastIndex(buff.String(), "stderr_lines")
		tempStr := buff.String()[firstIndex+9 : lastIndex-1]
		tmpLastIndex := strings.LastIndex(tempStr, ",")
		return fmt.Errorf(tempStr[0:tmpLastIndex])
	}

	return nil
}
