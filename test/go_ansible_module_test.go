package ansible_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/apenella/go-ansible/pkg/adhoc"
	"github.com/apenella/go-ansible/pkg/options"
)

func TestAsibleModulePing(t *testing.T) {

	ansibleConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
	}

	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		Inventory:  " 127.0.0.1,",
		ModuleName: "ping",
		// Args:       "ping 127.0.0.1 -c 2",
	}

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern:           "all",
		Options:           ansibleAdhocOptions,
		ConnectionOptions: ansibleConnectionOptions,
		StdoutCallback:    "oneline",
	}

	fmt.Println("Command: ", adhoc.String())

	err := adhoc.Run(context.TODO())
	if err != nil {
		panic(err)
	}
}
