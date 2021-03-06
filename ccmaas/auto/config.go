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

// Config is configuration of project
type Config struct {
	Stack   string     `yaml:"name"`
	Project DeployType `yaml:"type"`
	Props   Props      `yaml:"props"`
}

type Props struct {
	Openstack OpenstackProps `yaml:"openstack"`
	Stack     interface{}    `yaml:"stack"`
}

type OpenstackProps struct {
	Region   string `yaml:"region"`
	Domain   string `yaml:"domain"`
	Tenant   string `yaml:"tenant"`
	UserName string `yaml:"user"`
	Password string
}

func ReadConfig(fpath string) (*Config, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Read(fpath string) error {
	yamlBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlBytes, c); err != nil {
		return err
	}
	return nil
}

func (c *Config) Write(fpath string) error {
	if c.Project == "" {
		return fmt.Errorf("%q: %v", "Project", ErrStringEmpty)
	}
	if c.Stack == "" {
		return fmt.Errorf("%q: %v", "Stack", ErrStringEmpty)
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fpath, b, 0644)
}

// func (c *Config) AddNode(n Node) error {
// 	for _, nn := range c.Props.Nodes {
// 		if nn.Name == n.Name {
// 			return fmt.Errorf("%q: %v", n.Name, ErrNodeExists)
// 		}
// 	}
// 	c.Props.Nodes = append(c.Props.Nodes, n)
// 	return nil
// }

// func (c *Config) validate() error {
// 	if c.Props.Domain == "" {
// 		return fmt.Errorf("%q: %v", "Props.Domain", ErrStringEmpty)
// 	}
// 	if c.Props.Tenant == "" {
// 		return fmt.Errorf("%q: %v", "Props.Tenant", ErrStringEmpty)
// 	}
// 	return nil
// }
