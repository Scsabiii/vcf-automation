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
	"context"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type EsxiStack struct {
	auto.Stack
	state *EsxiState
}

type EsxiState struct {
	err              error
	refreshError     error
	NodeNetworkName  string
	NodeNetworkID    string
	StorageNetworkID string
	SecurityGroupID  string
}

type EsxiStackProps struct {
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

func InitEsxiStack(ctx context.Context, stackName, projectDir string) (*EsxiStack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to create or select stack: %v\n", err)
	}
	return &EsxiStack{Stack: s, state: &EsxiState{}}, nil
}

// Config stack
func (s EsxiStack) Configure(ctx context.Context, cfg *Config) error {
	configureOpenstack(ctx, s.Stack, cfg)

	p := EsxiStackProps{}
	err := GetStackPropsFromConfig(cfg, &p)
	if err != nil {
		return err
	}
	if p.Prefix == "" {
		p.Prefix = cfg.Stack
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

func (s EsxiStack) UpdateConfig(ctx context.Context, payload *EsxiStackProps) error {
	return nil
}

func (s EsxiStack) Refresh(ctx context.Context) error {
	_, err := s.Stack.Refresh(ctx)
	if err != nil {
		s.state.refreshError = err
		return err
	}
	// printUpdateSummary(res.Summary)
	return nil
}

func (s EsxiStack) Update(ctx context.Context) (auto.UpResult, error) {
	res, err := s.Stack.Up(ctx)
	if err != nil {
		s.state.err = err
		return auto.UpResult{}, err
	}
	return res, nil
}

func (s EsxiStack) Destroy(ctx context.Context) error {
	res, err := s.Stack.Destroy(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
	fmt.Println(res)
	return nil
}

func (s EsxiStack) GetState() interface{} {
	return s.state
}

func (s EsxiStack) GetError() error {
	return s.state.err
}

func (s EsxiStack) GenYaml(ctx context.Context, cfg *Config) ([]byte, error) {
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

func (s EsxiStack) SetState() {
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
