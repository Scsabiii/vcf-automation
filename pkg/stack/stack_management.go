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

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type ManagementStack struct {
	auto.Stack
	state ManagementStackState
}

type ManagementStackState struct {
	err error
}

type ManagementStackProps struct {
	ExternalNetwork    MgmtDomainExternalNetwork   `yaml:"externalNetwork"`
	ManagementNetwork  MgmtDomainMgmtNetwork       `yaml:"managementNetwork"`
	DeploymentNetwork  MgmtDomainDeploymentNetwork `yaml:"deploymentNetwork"`
	HelperVM           HelperVM                    `yaml:"helperVM"`
	DNSZoneName        string                      `yaml:"dnsZoneName"`
	PublicRouter       string                      `yaml:"publicRouter"`
	ReverseDNSZoneName string                      `yaml:"reverseDnsZoneName"`
	EsxiServerImage    string                      `yaml:"esxiServerImange"`
	EsxiServerFlavorID string                      `yaml:"esxiServerFlavorID"`
	EsxiNodes          []MgmtDomainEsxiNode        `yaml:"esxiNodes"`
	PrivateNetworks    []MgmtDomainPrivateNetwork  `yaml:"privateNetworks"`
	ReservedIPs        []RerservedIP               `yaml:"reservedIPs"`
}

type MgmtDomainExternalNetwork struct {
	Name string `json:"name,omitempty" yaml:"name"`
	ID   string `json:"id,omitempty" yaml:"id"`
}

type MgmtDomainMgmtNetwork struct {
	NetworkName   string `yaml:"networkName" json:"name,omitempty"`
	SubnetName    string `yaml:"subnetName" json:"subnet_name,omitempty"`
	SubnetGateway string `yaml:"subnetGateway" json:"subnet_gateway,omitempty"`
	SubnetMask    string `yaml:"subnetMask" json:"subnet_mask,omitempty"`
	VlanID        int    `yaml:"vlanID" json:"vlan_id,omitempty"`
	EsxiInterface string `yaml:"esxiInterface" json:"esxi_interface,omitempty"`
}

type MgmtDomainDeploymentNetwork struct {
	NetworkName string `yaml:"networkName" json:"name,omitempty"`
	SubnetName  string `yaml:"subnetName" json:"subnet_name,omitempty"`
	CIDR        string `yaml:"cidr" json:"cidr,omitempty"`
	Gateway     string `yaml:"gatewayIP" json:"gateway_ip,omitempty"`
}

type MgmtDomainPrivateNetwork struct {
	NetworkName   string `yaml:"networkName" json:"name,omitempty"`
	CIDR          string `yaml:"cidr" json:"cidr,omitempty"`
	VlanID        int    `yaml:"vlanID" json:"vlan_id,omitempty"`
	EsxiInterface string `yaml:"esxiInterface" json:"esxi_interface,omitempty"`
}

type MgmtDomainEsxiNode struct {
	Name      string `yaml:"name" json:"name,omitempty"`
	ID        string `yaml:"id" json:"id,omitempty"`
	IP        string `yaml:"ip" json:"ip,omitempty"`
	ImageName string `yaml:"imageName" json:"image_name,omitempty"`
}

type HelperVM struct {
	FlavorID   string `yaml:"flavorID" json:"flavor_id,omitempty"`
	FlavorName string `yaml:"flavorName" json:"flavor_name,omitempty"`
	ImageName  string `yaml:"imageName" json:"image_name,omitempty"`
	IP         string `yaml:"ip" json:"ip,omitempty"`
}

type RerservedIP struct {
	IP   string `yaml:"ip" json:"ip,omitempty"`
	Name string `yaml:"name" json:"name,omitempty"`
}

func InitManagementStack(ctx context.Context, stackName, projectDir string) (*ManagementStack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return nil, err
	}
	return &ManagementStack{s, ManagementStackState{}}, nil
}

func (s *ManagementStack) Configure(ctx context.Context, cfg *Config) error {
	err := configureKeypair(ctx, s.Stack, cfg)
	if err != nil {
		return err
	}
	err = configureOpenstack(ctx, s.Stack, cfg)
	if err != nil {
		return err
	}
	p := ManagementStackProps{}
	err = GetStackPropsFromConfig(cfg, &p)
	if err != nil {
		return err
	}
	for _, n := range p.EsxiNodes {
		err = validateIronicNodes(n.Name, n.ID)
		if err != nil {
			return err
		}
	}
	if (p.ExternalNetwork != MgmtDomainExternalNetwork{}) {
		if en, err := json.Marshal(p.ExternalNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "externalNetwork", auto.ConfigValue{Value: string(en)})
		}
	}
	if (p.ManagementNetwork != MgmtDomainMgmtNetwork{}) {
		if mn, err := json.Marshal(p.ManagementNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "managementNetwork", auto.ConfigValue{Value: string(mn)})
		}
	}
	if (p.DeploymentNetwork != MgmtDomainDeploymentNetwork{}) {
		if dn, err := json.Marshal(p.DeploymentNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "deploymentNetwork", auto.ConfigValue{Value: string(dn)})
		}
	}
	if p.PrivateNetworks != nil {
		if pn, err := json.Marshal(p.PrivateNetworks); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "privateNetworks", auto.ConfigValue{Value: string(pn)})
		}
	}
	if (p.HelperVM != HelperVM{}) {
		if n, err := json.Marshal(p.HelperVM); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "helperVM", auto.ConfigValue{Value: string(n)})
		}
	}
	if p.ReservedIPs != nil {
		if n, err := json.Marshal(p.ReservedIPs); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "reservedIPs", auto.ConfigValue{Value: string(n)})
		}
	}
	if p.EsxiNodes != nil {
		if n, err := json.Marshal(p.EsxiNodes); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "esxiNodes", auto.ConfigValue{Value: string(n)})
		}
	}
	if p.EsxiServerImage != "" {
		s.SetConfig(ctx, "esxiServerImage", auto.ConfigValue{Value: p.EsxiServerImage})
	}
	if p.EsxiServerFlavorID != "" {
		s.SetConfig(ctx, "esxiServerFlavorID", auto.ConfigValue{Value: p.EsxiServerFlavorID})
	}
	if p.DNSZoneName != "" {
		s.SetConfig(ctx, "dnsZoneName", auto.ConfigValue{Value: p.DNSZoneName})
	}
	if p.ReverseDNSZoneName != "" {
		s.SetConfig(ctx, "reverseDnsZoneName", auto.ConfigValue{Value: p.ReverseDNSZoneName})
	}
	if p.PublicRouter != "" {
		s.SetConfig(ctx, "publicRouter", configValue(p.PublicRouter))
	}
	return nil
}

func (s *ManagementStack) GenYaml(ctx context.Context, cfg *Config) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (s *ManagementStack) Refresh(ctx context.Context) error {
	_, err := s.Stack.Refresh(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
	return nil
}

func (s *ManagementStack) Update(ctx context.Context) (auto.UpResult, error) {
	res, err := s.Stack.Up(ctx)
	if err != nil {
		s.state.err = err
		return auto.UpResult{}, err
	}
	return res, nil
}

func (s *ManagementStack) GetState() interface{} {
	return nil
}

func (s *ManagementStack) GetError() error {
	return nil
}
