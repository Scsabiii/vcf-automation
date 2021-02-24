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
	"ccmaas/auto"
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [projectName/stackName]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var ctx = context.Background()
		var workDir = viper.GetString("workdir")
		project, stack := extractProjectStack(args)
		c, err := auto.NewController(workDir, project, stack)
		if err != nil {
			logErrorAndExit(err)
		}
		err = c.InitStack(ctx)
		if err != nil {
			logErrorAndExit(err)
		}
		err = c.Destory(ctx)
		if err != nil {
			logErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
