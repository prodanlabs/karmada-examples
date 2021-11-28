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
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	KarmadaVersion = "unknown"
	Version        = "unknown"
	GitCommitID    = "unknown"
)

//var Verbose bool

// newCmdVersion output version information.
func newCmdVersion() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the current versions information.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s-karmada:%s\nGitCommitID: %s\nBuildDate: %s\nGoVersion: %s\nPlatform: %s/%s\n",
				Version, KarmadaVersion, GitCommitID, time.Now().UTC().Format("2006-01-02T15:04:05Z"), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}
	//cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	return cmd
}
