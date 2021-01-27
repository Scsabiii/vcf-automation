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

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name   string      `yaml:"name"`
	Type   DeployType  `yaml:"type"`
	Props  DeployProps `yaml:"props"`
	Nodes  []Node      `yaml:"nodes"`
	Shares []Share     `yaml:"shares"`
}

func ReadConfig(fpath string, c *Config) error {
	yamlBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlBytes, c)
}

func (c Config) Write(fpath string) error {
	if c.Name == "" {
		return fmt.Errorf("config name required")
	}
	if c.Props.Domain == "" {
		return fmt.Errorf("domain name required")
	}
	if c.Props.Project == "" {
		return fmt.Errorf("project name required")
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fpath, b, 0644)
}

func (c Config) AddNode(n Node) error {
	for _, nn := range c.Nodes {
		if nn.Name == n.Name {
			return fmt.Errorf("node %q exists in config %q", n.Name, c.Name)
		}
	}
	c.Nodes = append(c.Nodes, n)
	return nil
}
