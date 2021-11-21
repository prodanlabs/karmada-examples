/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewKaAdmCommand returns cobra.Command to run kaAdm command
func NewKaAdmCommand() *cobra.Command {

	cmds := &cobra.Command{
		Use:   "kaadm",
		Short: "kaadm: easy and fast installation of karmada.",
		Long: dedent.Dedent(`
		┌─────────────────────────────────────────────────────────────────────────┐
		│ KAAdm                                                                   │
		│ Easy and fast installation of karmada.                                  │
		│                                                                         │
		│ Example usage:                                                          │
		│   Install karmada on kubernetes.                                        │
		│   (run kaadm on the kubernetes master node)                             │
		│   # kaadm install --master=xxx.xxx.xxx.xxx                              │
		│                                                                         │
		│   Install karmada on Linux.                                             │
		│   (Not yet implemented)                                                 │
		└─────────────────────────────────────────────────────────────────────────┘

`),

		SilenceErrors: true,
		SilenceUsage:  true,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },

	}

	cmds.ResetFlags()
	cmds.AddCommand(newCmdVersion())
	cmds.AddCommand(newCmdInstall())
	return cmds
}
