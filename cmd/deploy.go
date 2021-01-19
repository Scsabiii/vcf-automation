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
	"path"

	"ccmaas/auto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	destroy bool
	outputs bool
	yaml    bool
	stack   auto.Stack
)

var deployCmd = &cobra.Command{
	Use:   "deploy [cfgFileName]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploy(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&destroy, "destory", "d", false, "Destory stack")
	deployCmd.Flags().BoolVarP(&outputs, "outputs", "o", false, "Outputs of stack")
	deployCmd.Flags().BoolVarP(&yaml, "yaml", "y", false, "Yaml output")
}

func deploy(cfgName string) {
	ctx := context.Background()
	cfg := readConfig(cfgName)

	// Initial pulumi stack
	switch cfg.Type {
	case auto.DeployExample:
		cfg.Stack = "dev"
		stack = auto.InitExampleStack(ctx, cfg)
	case auto.DeployEsxi:
		cfg.Stack = cfgName
		stack = auto.InitEsxiStack(ctx, cfg)
	}

	if outputs {
		fmt.Println()
		fmt.Println("Outputs")
		fmt.Println("-------")
		auto.PrintOutputs(ctx, stack)
	} else if yaml {
		yamlOutput := auto.GenYaml(ctx, cfg, stack)
		fmt.Println()
		fmt.Println("Yaml Outputs")
		fmt.Println("-------")
		fmt.Println(string(yamlOutput))
	} else {
		auto.RunStack(ctx, stack, destroy)
	}
}

func readConfig(cfgName string) auto.Config {
	cfgPath := path.Join("etc", fmt.Sprintf("%s.yaml", cfgName))
	cfg, err := auto.ReadConfig(cfgPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if cfg.Type == "" {
		cfg.Type = auto.DeployEsxi
	}
	cfg.Props.Password = viper.GetString("os_password")
	return cfg
}
