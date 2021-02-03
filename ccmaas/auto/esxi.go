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
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"gopkg.in/yaml.v2"
)

type EsxiStack struct {
	*auto.Stack
}

func InitEsxiStack(ctx context.Context, stackName, projectDir string) (EsxiStack, error) {
	fmt.Printf("Use project %q\n", projectDir)
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return EsxiStack{}, fmt.Errorf("Failed to create or select stack: %v\n", err)
	}
	fmt.Printf("Created/Selected stack %q\n", stackName)
	return EsxiStack{&s}, nil
}

// Config stack
func (s EsxiStack) Configure(ctx context.Context, cfg Config) error {
	deployProps := cfg.Props
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", deployProps.Region)
	osProjectDomainName := deployProps.Domain
	osTenantName := deployProps.Tenant
	osUserName := deployProps.UserName
	osPassword := deployProps.Password

	// config openstack
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: deployProps.Region})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osTenantName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})

	s.SetConfig(ctx, "resourcePrefix", auto.ConfigValue{Value: deployProps.Prefix})
	s.SetConfig(ctx, "nodeSubnet", auto.ConfigValue{Value: deployProps.NodeSubnet})
	s.SetConfig(ctx, "storageSubnet", auto.ConfigValue{Value: deployProps.StorageSubnet})
	s.SetConfig(ctx, "shareNetworkUUID", auto.ConfigValue{Value: deployProps.ShareNetworkName})

	// config instance
	s.SetConfig(ctx, "numNodes", auto.ConfigValue{Value: strconv.Itoa(len(cfg.Props.Nodes))})
	for i, node := range cfg.Props.Nodes {
		s.SetConfig(ctx, fmt.Sprintf("node%02dImageName", i), auto.ConfigValue{Value: node.ImageName})
		s.SetConfig(ctx, fmt.Sprintf("node%02dFlavorName", i), auto.ConfigValue{Value: node.FlavorName})
		s.SetConfig(ctx, fmt.Sprintf("node%02dUUID", i), auto.ConfigValue{Value: node.UUID})
		s.SetConfig(ctx, fmt.Sprintf("node%02dIP", i), auto.ConfigValue{Value: node.IP})
	}
	s.SetConfig(ctx, "numShares", auto.ConfigValue{Value: strconv.Itoa(len(cfg.Props.Shares))})
	for i, share := range cfg.Props.Shares {
		s.SetConfig(ctx, fmt.Sprintf("share%02dName", i), auto.ConfigValue{Value: share.Name})
		s.SetConfig(ctx, fmt.Sprintf("share%02dSize", i), auto.ConfigValue{Value: share.Size})
	}
	return nil
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
