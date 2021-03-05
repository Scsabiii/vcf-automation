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

package auto

const (
	DeployEsxi    DeployType = "esxi"
	DeployExample DeployType = "example"
)

type DeployType string

type DeployProps struct {
	Region           string `yaml:"region"`
	Domain           string `yaml:"domain"`
	Tenant           string `yaml:"tenant"`
	UserName         string `yaml:"user"`
	Prefix           string `yaml:"resourcePrefix"`
	NodeSubnet       string `yaml:"nodeSubnet"`
	StorageSubnet    string `yaml:"storageSubnet"`
	ShareNetworkName string `yaml:"shareNetworkName"`
	Password         string
	Nodes            []Node  `yaml:"nodes"`
	Shares           []Share `yaml:"shares"`
}

type Node struct {
	Name   string `yaml:"name"`
	UUID   string `yaml:"uuid"`
	IP     string `yaml:"ip"`
	Image  string `yaml:"image"`
	Flavor string `yaml:"flavor"`
}

type Share struct {
	Name string `yaml:"name"`
	Size int    `yaml:"size"`
}