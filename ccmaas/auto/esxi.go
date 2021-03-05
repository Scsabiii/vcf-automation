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
	"context"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type EsxiStack struct {
	*auto.Stack
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

func InitEsxiStack(ctx context.Context, stackName, projectDir string) (EsxiStack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return EsxiStack{}, fmt.Errorf("Failed to create or select stack: %v\n", err)
	}
	return EsxiStack{Stack: &s, state: &EsxiState{}}, nil
}

// Config stack
func (s EsxiStack) Configure(ctx context.Context, cfg Config) error {
	if cfg.Props.Region == "" {
		return fmt.Errorf("Config.Props.Region not set")
	}
	if cfg.Props.Domain == "" {
		return fmt.Errorf("Config.Props.Domain not set")
	}
	if cfg.Props.Tenant == "" {
		return fmt.Errorf("Config.Props.Tenant not set")
	}
	if cfg.Props.UserName == "" {
		return fmt.Errorf("Config.Props.UserName not set")
	}
	if cfg.Props.NodeSubnet == "" {
		return fmt.Errorf("Config.Props.NodeSubnet not set")
	}
	if cfg.Props.StorageSubnet == "" {
		return fmt.Errorf("Config.Props.StorageSubnet not set")
	}

	osPassword := viper.GetString("os_password")
	if osPassword == "" {
		return fmt.Errorf("env variable CCMAAS_OS_PASSWORD not configured")
	}

	osRegion := cfg.Props.Region
	osProjectDomainName := cfg.Props.Domain
	osTenantName := cfg.Props.Tenant
	osUserName := cfg.Props.UserName
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", osRegion)

	// config openstack
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: osRegion})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osTenantName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})

	s.SetConfig(ctx, "resourcePrefix", auto.ConfigValue{Value: cfg.Props.Prefix})
	s.SetConfig(ctx, "nodeSubnet", auto.ConfigValue{Value: cfg.Props.NodeSubnet})
	s.SetConfig(ctx, "storageSubnet", auto.ConfigValue{Value: cfg.Props.StorageSubnet})
	s.SetConfig(ctx, "shareNetworkUUID", auto.ConfigValue{Value: cfg.Props.ShareNetworkName})

	nodes, err := json.Marshal(cfg.Props.Nodes)
	if err != nil {
		return err
	}
	s.SetConfig(ctx, "nodes", auto.ConfigValue{Value: string(nodes)})

	shares, err := json.Marshal(cfg.Props.Shares)
	if err != nil {
		return err
	}
	s.SetConfig(ctx, "shares", auto.ConfigValue{Value: string(shares)})

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

func (s EsxiStack) Update(ctx context.Context) error {
	res, err := s.Stack.Up(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
	// printUpdateSummary(res.Summary)
	printStackOutputs(res.Outputs)
	return nil
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

func (s EsxiStack) State() interface{} {
	return s.state
}

func (s EsxiStack) Error() error {
	return s.state.err
}

func (s EsxiStack) GenYaml(ctx context.Context, cfg Config) ([]byte, error) {
	outputs, err := s.Outputs(ctx)
	if err != nil {
		fmt.Printf("PrintYaml: %v\n", err)
		return nil, err
	}
	nodes := make([]NodeOutput, len(cfg.Props.Nodes))
	for i := 0; i < len(cfg.Props.Nodes); i++ {
		id, err := lookupOutput(outputs, fmt.Sprintf("EsxiInstance%02dID", i))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		nodes[i].ID = id
		ip, err := lookupOutput(outputs, fmt.Sprintf("EsxiInstance%02dIP", i))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		nodes[i].IP = ip
	}
	res, err := yaml.Marshal(YamlOutput{nodes})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
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