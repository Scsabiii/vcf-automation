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
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func GetConfig(stack string) (cfg Config, err error) {
	cfgPath := path.Join("etc", fmt.Sprintf("%s.yaml", stack))
	yamlBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		return
	}
	return
}

func AddNode(stack string, n Node) error {
	c, err := GetConfig(stack)
	if err != nil {
		return err
	}
	err = c.AddNode(n)
	if err != nil {
		return err
	}
	err = c.Save(true)
	if err != nil {
		return err
	}
	return nil
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

func (c Config) Save(overwrite bool) error {
	if c.Name == "" {
		return fmt.Errorf("config name required")
	}
	if c.Props.Domain == "" {
		return fmt.Errorf("domain name required")
	}
	if c.Props.Project == "" {
		return fmt.Errorf("project name required")
	}
	fpath := path.Join("etc", fmt.Sprintf("%s.yaml", c.Name))
	return writeConfigToFile(c, fpath, overwrite)
}

func writeConfigToFile(cfg Config, fpath string, overwrite bool) error {
	if !overwrite {
		if fileExists(fpath) {
			return fmt.Errorf("file %q exists", fpath)
		}
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fpath, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
