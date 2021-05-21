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

package esxi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type Stack struct {
	*auto.Stack
	state StackState
}

type StackState struct {
	err              error
	refreshError     error
	NodeNetworkName  string
	NodeNetworkID    string
	StorageNetworkID string
	SecurityGroupID  string
}

type StackProps struct {
	Prefix           string  `yaml:"resourcePrefix"`
	NodeSubnet       string  `yaml:"nodeSubnet"`
	StorageSubnet    string  `yaml:"storageSubnet"`
	ShareNetworkName string  `yaml:"shareNetworkName"`
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

func InitEsxiStack(ctx context.Context, stackName, projectDir string) (*Stack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create or select stack: %v", err)
	}
	return &Stack{&s, StackState{}}, nil
}

// Config Esxi Project specific properties
func (s *Stack) Configure(ctx context.Context, props ...StackProps) error {
	p := props[0]
	for _, q := range props[1:] {
		mergo.Merge(&p, q)
	}
	if p.NodeSubnet == "" {
		return fmt.Errorf("Config.Props.Stack.NodeSubnet not set")
	}
	if p.StorageSubnet == "" {
		return fmt.Errorf("Config.Props.Stack.StorageSubnet not set")
	}

	// set project settings
	s.SetConfig(ctx, "resourcePrefix", auto.ConfigValue{Value: p.Prefix})
	s.SetConfig(ctx, "nodeSubnet", auto.ConfigValue{Value: p.NodeSubnet})
	s.SetConfig(ctx, "storageSubnet", auto.ConfigValue{Value: p.StorageSubnet})
	s.SetConfig(ctx, "shareNetworkUUID", auto.ConfigValue{Value: p.ShareNetworkName})
	nodes, err := json.Marshal(p.Nodes)
	if err != nil {
		return err
	}
	s.SetConfig(ctx, "nodes", auto.ConfigValue{Value: string(nodes)})
	shares, err := json.Marshal(p.Shares)
	if err != nil {
		return err
	}
	s.SetConfig(ctx, "shares", auto.ConfigValue{Value: string(shares)})

	return nil
}

func (s *Stack) UpdateConfig(ctx context.Context, payload *StackProps) error {
	return nil
}

func (s *Stack) Refresh(ctx context.Context) error {
	_, err := s.Stack.Refresh(ctx)
	if err != nil {
		s.state.refreshError = err
		return err
	}
	// printUpdateSummary(res.Summary)
	return nil
}

func (s *Stack) Update(ctx context.Context) (auto.UpResult, error) {
	res, err := s.Stack.Up(ctx)
	if err != nil {
		s.state.err = err
		return auto.UpResult{}, err
	}
	return res, nil
}

func (s *Stack) Destroy(ctx context.Context) error {
	res, err := s.Stack.Destroy(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
	fmt.Println(res)
	return nil
}

func (s *Stack) GetState() interface{} {
	return s.state
}

func (s *Stack) GetError() error {
	return s.state.err
}

func (s *Stack) GenYaml(ctx context.Context, p *StackProps) ([]byte, error) {
	// outputs, err := s.Outputs(ctx)
	// if err != nil {
	// 	fmt.Printf("PrintYaml: %v\n", err)
	// 	return nil, err
	// }
	// nodes := make([]NodeOutput, len(p.Nodes))
	// for i := 0; i < len(cfg.Props.Stack.Nodes); i++ {
	// 	id, err := lookupOutput(outputs, fmt.Sprintf("EsxiInstance%02dID", i))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil, err
	// 	}
	// 	nodes[i].ID = id
	// 	ip, err := lookupOutput(outputs, fmt.Sprintf("EsxiInstance%02dIP", i))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil, err
	// 	}
	// 	nodes[i].IP = ip
	// }
	// res, err := yaml.Marshal(YamlOutput{nodes})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	// return res, nil
	return nil, nil
}

func lookupOutput(outputs auto.OutputMap, key string) (string, error) {
	for k, v := range outputs {
		if k == key {
			// TODO validate value type
			return v.Value.(string), nil
		}
	}
	err := fmt.Errorf("Key %q not found", key)
	return "", err
}

func (s *Stack) SetState() {
	// for k, v := range res.Outputs {
	// 	switch k {
	// 	case "EsxiNetworkName":
	// 		if vv, ok := v.(string); ok {
	// 			s.state.NodeNetworkName = vv
	// 		}
	// 	case "EsxiNetworkID":
	// 		if vv, ok := v.(string); ok {
	// 			s.state.NodeNetworkID = vv
	// 		}
	// 	default:
	// 	}
	// }
}
