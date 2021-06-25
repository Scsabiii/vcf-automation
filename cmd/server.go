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
	"github.com/sapcc/vcf-automation/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "start server",
	Long:  `start server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(port)
	},
}

func init() {
	viper.SetEnvPrefix("automation")
	viper.SetDefault("port", 8080)
	viper.AutomaticEnv()

	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", viper.GetInt("port"), "server port")
}
