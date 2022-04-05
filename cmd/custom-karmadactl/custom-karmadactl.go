package main

import (
	"os"

	"github.com/prodanlabs/karmada-examples/pkg/karmadactl"
)

func main() {
	if err := karmadactl.NewCustomKarmadaCtlCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
