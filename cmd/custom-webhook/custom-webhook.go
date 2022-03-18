package main

import (
	"os"

	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"

	"github.com/prodanlabs/karmada-examples/cmd/custom-webhook/app"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	ctx := apiserver.SetupSignalContext()
	if err := app.NewWebhookCommand(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
