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

package stack

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/sapcc/vcf-automation/pkg/stack/esxi"
	"gopkg.in/yaml.v2"
)

const (
	ProjectEsxi          ProjectType = "esxi"
	ProjectExample       ProjectType = "example-go"
	ProjectVCFManagement ProjectType = "vcf/management"
	ProjectVCFWorkload   ProjectType = "vcf/workload"
)

// Config is configuration of project/stack
type Config struct {
	ProjectType ProjectType `json:"project_type" yaml:"projectType"`
	Stack       string      `json:"stack" yaml:"stack"`
	Props       Props       `json:"props" yaml:"props"`
	DependsOn   []string    `json:"depends_on,omitempty" yaml:"dependsOn"`
}

// ProjectType is project type
type ProjectType string

// Props is configuration needed for pulumi projects. It holds general
// configuration for openstack services and project specific stack props
type Props struct {
	OpenstackProps OpenstackProps `json:"openstack" yaml:"openstack"`
	StackProps     StackProps     `json:"stack" yaml:"stack"`
	BaseStackProps []StackProps
	Keypair        Keypair
}

// OpenstackProps
type OpenstackProps struct {
	Region string `json:"region" yaml:"region"`
	Domain string `json:"domain" yaml:"domain"`
	Tenant string `json:"tenant" yaml:"tenant"`
}

// StackProps is a empty type, a placeholder for the project specific
// properties
type StackProps interface{}

// Keypair stores ssh key pairs, which is loaded from the disk
type Keypair struct {
	publicKey  string
	privateKey string
}

func ReadConfig(configFile string) (*Config, error) {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	for _, fname := range c.DependsOn {
		dc, err := ReadConfig(path.Join(path.Dir(configFile), fname))
		if err != nil {
			return nil, err
		}
		c.Props.BaseStackProps = append(c.Props.BaseStackProps, dc.Props.StackProps)
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) validate() error {
	if c.ProjectType == "" {
		return fmt.Errorf("project not set")
	}
	if c.Stack == "" {
		return fmt.Errorf("stack not set")
	}
	return nil
}

func WriteNewConfig(fpath string, c *Config) error {
	return writeConfig(fpath, c, false)
}

func WriteConfig(fpath string, c *Config) error {
	return writeConfig(fpath, c, true)
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
	switch c.ProjectType {
	case ProjectEsxi:
		p := esxi.StackProps{}
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
	switch nc.ProjectType {
	case ProjectEsxi:
		p := esxi.StackProps{}
		err := unmarshalStackProps(c.Props.StackProps, &p)
		if err != nil {
			return nil, err
		}
		np := esxi.StackProps{}
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

// Function unmarshalStackProps() deserializes the StackProps (s) into struct
// (props), according to the actual type of (props). E.g.,
//		p := EsxiStackProps{}
//      unmarshalStackProps(s, p)
func unmarshalStackProps(s StackProps, props interface{}) error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, props)
}

// deserialize StackProps slice
func unmarshalStackPropList(s []StackProps, props interface{}) error {
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, props)
}
