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
	"os"
	"path/filepath"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type EsxiStack struct {
	*auto.Stack
	config Config
}

func InitEsxiStack(ctx context.Context, cfg Config) EsxiStack {
	workDir := filepath.Join(".", "projects", "esxi")

	fmt.Printf("Use project %q\n", workDir)

	s, err := auto.UpsertStackLocalSource(ctx, cfg.Name, workDir)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", cfg.Name)

	return EsxiStack{&s, cfg}
}

// Config stack
func (s EsxiStack) Config(ctx context.Context) error {
	deployProps := s.config.Props
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", deployProps.Region)
	osProjectDomainName := deployProps.Domain
	osProjectName := deployProps.Project
	osUserName := deployProps.UserName
	osPassword := deployProps.Password

	// config openstack
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: deployProps.Region})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osProjectName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})

	s.SetConfig(ctx, "resourcePrefix", auto.ConfigValue{Value: deployProps.Prefix})
	s.SetConfig(ctx, "nodeSubnet", auto.ConfigValue{Value: deployProps.NodeSubnet})
	s.SetConfig(ctx, "storageSubnet", auto.ConfigValue{Value: deployProps.StorageSubnet})
	s.SetConfig(ctx, "shareNetworkUUID", auto.ConfigValue{Value: deployProps.ShareNetworkName})

	// config instance
	s.SetConfig(ctx, "numNodes", auto.ConfigValue{Value: strconv.Itoa(len(s.config.Nodes))})
	for i, node := range s.config.Nodes {
		s.SetConfig(ctx, fmt.Sprintf("node%02dImageName", i), auto.ConfigValue{Value: node.ImageName})
		s.SetConfig(ctx, fmt.Sprintf("node%02dFlavorName", i), auto.ConfigValue{Value: node.FlavorName})
		s.SetConfig(ctx, fmt.Sprintf("node%02dUUID", i), auto.ConfigValue{Value: node.UUID})
		s.SetConfig(ctx, fmt.Sprintf("node%02dIP", i), auto.ConfigValue{Value: node.IP})
	}
	s.SetConfig(ctx, "numShares", auto.ConfigValue{Value: strconv.Itoa(len(s.config.Shares))})
	for i, share := range s.config.Shares {
		s.SetConfig(ctx, fmt.Sprintf("share%02dName", i), auto.ConfigValue{Value: share.Name})
		s.SetConfig(ctx, fmt.Sprintf("share%02dSize", i), auto.ConfigValue{Value: share.Size})
	}
	return nil
}
