package karmadactl

import (
	"github.com/prodanlabs/karmada-examples/pkg/karmadactl/options"
	"github.com/spf13/cobra"
)

const ctlCommandName = "custom-karmadactl"

// NewCustomKarmadaCtlCommand new `custom-karmadactl` command
func NewCustomKarmadaCtlCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  ctlCommandName,
		Long: "custom karmadactl",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	opts := options.NewGlobalOptions()
	opts.AddFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(NewLogsPull(ctlCommandName, opts))
	return rootCmd
}
