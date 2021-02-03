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

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type ExampleStack struct {
	*auto.Stack
}

func InitExampleStack(ctx context.Context, stackName, projectDir string) (ExampleStack, error) {
	fmt.Printf("Use project %q\n", projectDir)
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		e := fmt.Errorf("failed to create/select stack: %v", err)
		return ExampleStack{}, e
	}
	fmt.Printf("Created/Selected stack %q\n", stackName)
	return ExampleStack{&s}, nil
}

// Config set stack configuration
func (s ExampleStack) Configure(ctx context.Context, cfg Config) error {
	props := cfg.Props
	osRegion := props.Region
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", osRegion)
	osProjectDomainName := props.Domain
	osTenantName := props.Tenant
	osUserName := props.UserName
	osPassword := props.Password
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: osRegion})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osTenantName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})
	return nil
}

func (s ExampleStack) GenYaml(ctx context.Context, cfg Config) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
