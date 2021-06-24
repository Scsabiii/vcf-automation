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
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	ProjectEsxi          ProjectType = "esxi"
	ProjectExample       ProjectType = "example-go"
	ProjectVCFManagement ProjectType = "vcf/management"
	ProjectVCFWorkload   ProjectType = "vcf/workload"
)

// ProjectType is project type
type ProjectType string

// StackProps is a empty type, a placeholder for the project specific
// properties
type StackProps interface{}

// Config is configuration of project/stack
type Config struct {
	ProjectType    ProjectType `json:"project_type" yaml:"projectType"`
	StackName      string      `json:"stack" yaml:"stack"`
	Props          Props       `json:"props" yaml:"props"`
	DependsOn      []string    `json:"depends_on,omitempty" yaml:"dependsOn"`
	baseStackProps []StackProps
}

// Props is configuration needed for pulumi projects. It holds general
// configuration for openstack services and project specific stack props
type Props struct {
	OpenstackProps OpenstackProps `json:"openstack" yaml:"openstack"`
	StackProps     StackProps     `json:"stack" yaml:"stack"`
}

// OpenstackProps
type OpenstackProps struct {
	Region string `json:"region" yaml:"region"`
	Domain string `json:"domain" yaml:"domain"`
	Tenant string `json:"tenant" yaml:"tenant"`
}

// ReadConfig reads config from config file full path (configFilePath)
func ReadConfig(configFilePath string) (*Config, error) {
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		err = fmt.Errorf("%s: %v", configFilePath, err)
		return nil, err
	}
	for _, fname := range c.DependsOn {
		dc, err := ReadConfig(path.Join(path.Dir(configFilePath), fname))
		if err != nil {
			return nil, err
		}
		c.baseStackProps = append(c.baseStackProps, dc.Props.StackProps)
	}
	return &c, nil
}

// GetProjectStackName returns stack's project type and stack name. If project
// type is composite with main and sub type, as mainProjectType/subProjectType,
// only main project type is returned.
func (c *Config) GetProjectStackName() (projectType, stackName string) {
	p := strings.Split(string(c.ProjectType), "/")
	return p[0], c.StackName
}

func (c *Config) validate() error {
	if c.ProjectType == "" {
		return fmt.Errorf("projectType not set")
	}
	if c.StackName == "" {
		return fmt.Errorf("stack not set")
	}
	return nil
}

// UnmarshalStackProps() decodes StackProps data into typed props, according to
// the actual type of props.
//
// For example:
//		p := EsxiStackProps{}
//      UnmarshalStackProps(s, p)
func UnmarshalStackProps(data StackProps, props interface{}) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, props)
}

// UnmarshalStackPropList decodes StackProps slice into typed props, according
// to the actual type of porps.
func UnmarshalStackPropList(data []StackProps, props interface{}) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, props)
}

// func writeConfig(fpath string, c *Config, overwrite bool) error {
// 	if !overwrite {
// 		m, _ := filepath.Glob(fpath)
// 		if m != nil {
// 			err := fmt.Errorf("file %q exists", fpath)
// 			return err
// 		}
// 	}

// 	// Need to unmarshal the field Props.StackProps, of type interface{}, to
// 	// the actual structure it has. Otherwise it is serialized into a raw
// 	// string
// 	switch c.ProjectType {
// 	case ProjectEsxi:
// 		p := esxi.StackProps{}
// 		err := unmarshalStackProps(c.Props.StackProps, &p)
// 		if err != nil {
// 			return err
// 		}
// 		c.Props.StackProps = p
// 	default:
// 	}
// 	b, err := yaml.Marshal(c)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(fpath, b, 0644)
// }

// MergeStackPropsToConfig merges the Props.StackProps field from s into Config c.
// NOTE: only EsxiStackProps.Nodes and EsxiStackProps.Shares are merged
// func MergeStackPropsToConfig(c *Config, s StackProps) (*Config, error) {
// 	// deep copy old config to nc
// 	nc := *c
// 	switch nc.ProjectType {
// 	case ProjectEsxi:
// 		p := esxi.StackProps{}
// 		err := unmarshalStackProps(c.Props.StackProps, &p)
// 		if err != nil {
// 			return nil, err
// 		}
// 		np := esxi.StackProps{}
// 		err = unmarshalStackProps(s, &np)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if np.Nodes != nil {
// 			p.Nodes = append(p.Nodes, np.Nodes...)
// 		}
// 		if np.Shares != nil {
// 			p.Shares = append(p.Shares, np.Shares...)
// 		}
// 		nc.Props.StackProps = p
// 	default:
// 		return nil, fmt.Errorf("merging configuration not supported")
// 	}
// 	return &nc, nil
// }
