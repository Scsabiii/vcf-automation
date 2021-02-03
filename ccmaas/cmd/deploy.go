/******************************************************************************
*
*  Copyright 2021 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"ccmaas/auto"

	"github.com/spf13/cobra"
)

var (
	destroy bool
	outputs bool
	yaml    bool
	stack   auto.Stack
	ctl     auto.Controller
)

var deployCmd = &cobra.Command{
	Use:   "deploy [projectName/stackName]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var project, stack string
		projectStackNames := strings.Split(args[0], "/")
		if len(projectStackNames) == 1 {
			project = "esxi"
			stack = projectStackNames[0]
		} else if len(projectStackNames) == 2 {
			project = projectStackNames[0]
			stack = projectStackNames[1]
		} else {
			err := fmt.Errorf("arg must be of format [projectName/][stackName]")
			fmt.Println(err)
			os.Exit(1)
		}
		deploy(project, stack)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&destroy, "destory", "d", false, "Destory stack")
	deployCmd.Flags().BoolVarP(&outputs, "outputs", "o", false, "Outputs of stack")
	deployCmd.Flags().BoolVarP(&yaml, "yaml", "y", false, "Yaml output")
}

func deploy(projectName, stackName string) {
	ctx := context.Background()
	ctl, err := auto.NewController(workDir, projectName, stackName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	stack, err = ctl.InitStack(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if outputs {
		fmt.Println()
		fmt.Println("Outputs")
		fmt.Println("-------")
		auto.PrintStackOutputs(ctx, stack)
	} else if yaml {
		yamlOutput, err := stack.GenYaml(ctx, ctl.Config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("Yaml Outputs")
		fmt.Println("-------")
		fmt.Println(string(yamlOutput))
	} else {
		if destroy {
			ctl.DestoryStack(ctx, stack)
		} else {
			ctl.UpdateStack(ctx, stack)
		}
	}
}
