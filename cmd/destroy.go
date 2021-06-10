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
	"path"

	"github.com/sapcc/vcf-automation/pkg/stack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [projectName/stackName]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		workdir := viper.GetString("workdir")
		prjdir := path.Join(workdir, "projects")
		etcdir := path.Join(workdir, "etc")
		projectName, stackName := extractProjectStack(args)
		fname := fmt.Sprintf("%s-%s.yaml", projectName, stackName)
		fpath := path.Join(etcdir, fname)
		c, err := stack.NewControllerFromConfigFile(prjdir, fpath)
		if err != nil {
			logErrorAndExit(err)
		}
		err = c.InitStack(ctx)
		if err != nil {
			logErrorAndExit(err)
		}
		// err = c.DestoryStack(ctx)
		// if err != nil {
		// 	logErrorAndExit(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}