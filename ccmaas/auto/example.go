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

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type ExampleStack struct {
	*auto.Stack
	config Config
}

func InitExampleStack(ctx context.Context, cfg Config) ExampleStack {
	workDir := filepath.Join("projects", "example-go")

	fmt.Printf("Use project %q\n", workDir)

	s, err := auto.UpsertStackLocalSource(ctx, cfg.Name, workDir)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", cfg.Name)
	return ExampleStack{&s, cfg}
}

// Config set stack configuration
func (s ExampleStack) Config(ctx context.Context) error {
	props := s.config.Props
	osRegion := props.Region
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", osRegion)
	osProjectDomainName := props.Domain
	osProjectName := props.Project
	osUserName := props.UserName
	osPassword := props.Password
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: osRegion})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osProjectName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})
	return nil
}
