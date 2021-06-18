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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "autmoation",
	Short: "vcf automation",
	Long: `automation:

deploy VCF project with pulumi`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("workdir", "", "work directory")
	viper.BindPFlag("workdir", rootCmd.PersistentFlags().Lookup("workdir"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// read in environment variables that match
	viper.SetEnvPrefix("automation")
	viper.AutomaticEnv()
}

func logErrorAndExit(e error) {
	fmt.Println("ERROR", e)
	os.Exit(1)
}

func extractProjectStack(args []string) (project, stack string) {
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
	return
}
