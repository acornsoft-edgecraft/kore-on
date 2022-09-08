package main

import (
	"os"

	"kore-on/cmd/kore-on-provider-baremetal/app"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := app.NewProviderBaremetalCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
