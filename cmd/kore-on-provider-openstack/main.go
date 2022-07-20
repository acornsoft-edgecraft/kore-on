package main

import (
	"os"

	"kore-on/cmd/kore-on-provider-openstack/app"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	ctx := signals.SetupSignalHandler()
	if err := app.NewProviderOpenstackCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
