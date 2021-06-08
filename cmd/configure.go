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

var configureCmd = &cobra.Command{
	Use:   "configure [projectName/stackName]",
	Short: "configure project/stack",
	Long:  `Read automation's yaml configuration file and generate pulumi project/stack's config file`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		workdir := viper.GetString("work_dir")
		projectRoot := viper.GetString("project_root")
		configdir := viper.GetString("config_dir")
		if projectRoot == "" {
			projectRoot = path.Join(workdir, "projects")
		}
		if configdir == "" {
			configdir = path.Join(workdir, "etc")
		}
		projectName, stackName := extractProjectStack(args)
		fname := fmt.Sprintf("%s-%s.yaml", projectName, stackName)
		fpath := path.Join(configdir, fname)
		c, err := stack.NewControllerFromConfigFile(projectRoot, fpath)
		if err != nil {
			logErrorAndExit(err)
		}
		err = c.InitStack(ctx)
		if err != nil {
			logErrorAndExit(err)
		}
		err = c.ConfigureStack(ctx)
		if err != nil {
			logErrorAndExit(err)
		}
		fmt.Printf("successfully configured the stack %s in project %s\n", stackName, projectName)
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
