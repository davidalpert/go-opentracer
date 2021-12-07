// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/davidalpert/gopentracer/internal/utils"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Execute builds the default root command and invokes it with os.Args
func Execute() {
	rootCmd := NewRootCmd()

	rootCmd.SetArgs(os.Args[1:]) // without program

	err := rootCmd.Execute()

	utils.ExitIfErr(err)
}

// NewRootCmd creates the root command with default arguments
func NewRootCmd() *cobra.Command {
	return NewDefaultRootCommandWithArgs(os.Args, os.Stdin, os.Stdout, os.Stderr)
}

// NewDefaultRootCommandWithArgs creates the root command with explicit arguments (exposing a seam decoupled from the environment)
func NewDefaultRootCommandWithArgs(args []string, in io.Reader, out, errout io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "gopentracer",
		//Aliases:       []string{},
		Short:         "gopentracer executes a shell command in an open trace",
		SilenceUsage:  true,
		SilenceErrors: true,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//RunE: func(cmd *cobra.Command, args []string) error {
		//},
	}

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//bindPersistentFlags(rootCmd)

	// Cobra also supports local Flags(), which will only run
	// when this action is called directly.
	//bindLocalFlags(rootCmd)

	rootCmd.AddCommand(NewCmdRun())
	rootCmd.AddCommand(NewCmdVersion())

	return rootCmd
}

func init() {
	// cobra.OnInitialize does not invoke these handlers directly but appends
	// them to an array of initializers which are invoked when the command's
	// Run field is executed; this is why they can be configured in advance.
	//
	// Any startup code which needs config needs to delay until after
	// cobra (with its flag bindings) has been initialized
	//cobra.OnInitialize(appInitialized)
}
