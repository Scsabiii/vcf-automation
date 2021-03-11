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
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	DeployEsxi    ProjectType = "esxi"
	DeployExample ProjectType = "example"
)

// Config is configuration of project/stack
type Config struct {
	Project ProjectType `json:"project" yaml:"project"`
	Stack   string      `json:"stack" yaml:"stack"`
	Props   Props       `json:"props" yaml:"props"`
}

// ProjectType is project type
type ProjectType string

// Props is configuration needed for pulumi projects. It holds general
// configuration for openstack services and project specific stack props
type Props struct {
	OpenstackProps OpenstackProps `json:"openstack" yaml:"openstack"`
	StackProps     StackProps     `json:"stack" yaml:"stack"`
}

// StackProps is a empty type, and should be used as project specific
// StackProps when possible
type StackProps interface{}

// OpenstackProps
type OpenstackProps struct {
	Region string `json:"region" yaml:"region"`
	Domain string `json:"domain" yaml:"domain"`
	Tenant string `json:"tenant" yaml:"tenant"`
}

// FileName generates the configuration file name with yaml extension
func (c *Config) FileName() string {
	return fmt.Sprintf("%s-%s.yaml", c.Project, c.Stack)
}

func readConfig(fpath string) (*Config, error) {
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

func writeConfig(fpath string, c *Config, overwrite bool) error {
	if !overwrite {
		m, _ := filepath.Glob(fpath)
		if m != nil {
			err := fmt.Errorf("file %q exists", fpath)
			return err
		}
	}

	// Need to unmarshal the field Props.StackProps, of type interface{}, to
	// the actual structure it has. Otherwise it is serialized into a raw
	// string
	switch c.Project {
	case DeployEsxi:
		p := EsxiStackProps{}
		err := unmarshalStackProps(c.Props.StackProps, &p)
		if err != nil {
			return err
		}
		c.Props.StackProps = p
	default:
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fpath, b, 0644)
}

// MergeStackPropsToConfig merges the Props.StackProps field from s into Config c.
// NOTE: only EsxiStackProps.Nodes and EsxiStackProps.Shares are merged
func MergeStackPropsToConfig(c *Config, s StackProps) (*Config, error) {
	// deep copy old config to nc
	nc := *c
	switch nc.Project {
	case DeployEsxi:
		p := EsxiStackProps{}
		err := unmarshalStackProps(c.Props.StackProps, &p)
		if err != nil {
			return nil, err
		}
		np := EsxiStackProps{}
		err = unmarshalStackProps(s, &np)
		if err != nil {
			return nil, err
		}
		if np.Nodes != nil {
			p.Nodes = append(p.Nodes, np.Nodes...)
		}
		if np.Shares != nil {
			p.Shares = append(p.Shares, np.Shares...)
		}
		nc.Props.StackProps = p
	default:
		return nil, fmt.Errorf("merging configuration not supported")
	}
	return &nc, nil
}

// unmarshalStackProps deserializes the StackProps s into props, whose actual
// type is assigned before calling this function.
// E.g.,
//		p := EsxiStackProps{}
//      unmarshalStackProps(s, p)
func unmarshalStackProps(s StackProps, props interface{}) error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, props)
}

func getStackPropsFromConfig(cfg *Config, props interface{}) error {
	return unmarshalStackProps(cfg.Props.StackProps, props)
}
