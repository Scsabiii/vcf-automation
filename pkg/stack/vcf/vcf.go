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

package vcf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type Stack struct {
	auto.Stack
	state StackState
}

type StackState struct {
	err error
}

type StackProps struct {
	EsxiServerImage    string            `yaml:"esxiServerImage"`
	EsxiServerFlavor   string            `yaml:"esxiServerFlavor"`
	EsxiNodes          []EsxiNode        `yaml:"esxiNodes"`
	Shares             []Share           `yaml:"shares"`
	SDDCManager        SDDCManager       `yaml:"sddcManager"`
	Nsxt               Nsxt              `yaml:"nsxt"`
	NsxtManagers       []NsxtManager     `yaml:"nsxtManagers"`
	VCenter            VCenter           `yaml:"vcenter"`
	ExternalNetwork    ExternalNetwork   `yaml:"externalNetwork"`
	ManagementNetwork  MgmtNetwork       `yaml:"managementNetwork"`
	DeploymentNetwork  DeploymentNetwork `yaml:"deploymentNetwork"`
	PrivateNetworks    []PrivateNetwork  `yaml:"privateNetworks"`
	PublicRouter       string            `yaml:"publicRouter"`
	DNSZoneName        string            `yaml:"dnsZoneName"`
	ReverseDNSZoneName string            `yaml:"reverseDnsZoneName"`
	HelperVM           HelperVM          `yaml:"helperVM"`
	KeypairFile        KeypairFile       `yaml:"keypairFile"`
	ReservedIPs        []ReservedIP      `yaml:"reservedIPs"`
}

type ExternalNetwork struct {
	Name string `json:"name,omitempty" yaml:"name"`
	ID   string `json:"id,omitempty" yaml:"id"`
}

type MgmtNetwork struct {
	NetworkName   string `yaml:"networkName" json:"name,omitempty"`
	SubnetName    string `yaml:"subnetName" json:"subnet_name,omitempty"`
	SubnetGateway string `yaml:"subnetGateway" json:"subnet_gateway,omitempty"`
	SubnetMask    string `yaml:"subnetMask" json:"subnet_mask,omitempty"`
	SubnetCidr    string `yaml:"subnetCidr" json:"subnet_cidr,omitempty"`
	VlanID        int    `yaml:"vlanID" json:"vlan_id,omitempty"`
	EsxiInterface string `yaml:"esxiInterface" json:"esxi_interface,omitempty"`
}

type DeploymentNetwork struct {
	NetworkName string `yaml:"networkName" json:"name,omitempty"`
	SubnetName  string `yaml:"subnetName" json:"subnet_name,omitempty"`
	CIDR        string `yaml:"cidr" json:"cidr,omitempty"`
	Gateway     string `yaml:"gatewayIP" json:"gateway_ip,omitempty"`
}

type HelperVM struct {
	FlavorID   string `yaml:"flavorID" json:"flavor_id,omitempty"`
	FlavorName string `yaml:"flavorName" json:"flavor_name,omitempty"`
	ImageName  string `yaml:"imageName" json:"image_name,omitempty"`
	IP         string `yaml:"ip" json:"ip,omitempty"`
}

type KeypairFile struct {
	PublicKey  string `json:"public_key,omitempty" yaml:"publicKey"`
	PrivateKey string `json:"private_key,omitempty" yaml:"privateKey"`
}

type PrivateNetwork struct {
	NetworkName   string `yaml:"networkName" json:"name,omitempty"`
	CIDR          string `yaml:"cidr" json:"cidr,omitempty"`
	VlanID        int    `yaml:"vlanID" json:"vlan_id,omitempty"`
	EsxiInterface string `yaml:"esxiInterface" json:"esxi_interface,omitempty"`
}

type EsxiNode struct {
	Name      string `yaml:"name" json:"name,omitempty"`
	ID        string `yaml:"id" json:"id,omitempty"`
	IP        string `yaml:"ip" json:"ip,omitempty"`
	ImageName string `yaml:"imageName" json:"image_name,omitempty"`
}

type ReservedIP struct {
	IP       string `yaml:"ip" json:"ip,omitempty"`
	Hostname string `yaml:"hostname" json:"hostname,omitempty"`
}

type Share struct {
	ShareName string `yaml:"shareName" json:"share_name,omitempty"`
	ShareSize int    `yaml:"shareSize" json:"share_size,omitempty"`
}

type SDDCManager struct {
	ID         string `yaml:"id" json:"id"`
	IP         string `yaml:"ip" json:"ip"`
	Hostname   string `yaml:"hostname" json:"hostname"`
	Domain     string `yaml:"domain" json:"domain"`
	EsxLicense string `yaml:"esxLicense" json:"esx_license"`
}

type Nsxt struct {
	IP       string `yaml:"ip" json:"ip,omitempty"`
	Hostname string `yaml:"hostname" json:"hostname,omitempty"`
	License  string `yaml:"license" json:"license,omitempty"`
}

type NsxtManager struct {
	IP       string `yaml:"ip" json:"ip,omitempty"`
	Hostname string `yaml:"hostname" json:"hostname,omitempty"`
}

type VCenter struct {
	IP       string `yaml:"ip" json:"ip,omitempty"`
	Hostname string `yaml:"hostname" json:"hostname,omitempty"`
}

func InitVCFStack(ctx context.Context, stackName, projectDir string) (*Stack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return nil, err
	}
	return &Stack{s, StackState{}}, nil
}

func (s *Stack) Configure(ctx context.Context, props ...StackProps) error {
	p := props[0]
	for _, q := range props[1:] {
		mergo.Merge(&p, q)
	}
	if (p.ExternalNetwork != ExternalNetwork{}) {
		if en, err := json.Marshal(p.ExternalNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "externalNetwork", auto.ConfigValue{Value: string(en)})
		}
	}
	if (p.ManagementNetwork != MgmtNetwork{}) {
		if mn, err := json.Marshal(p.ManagementNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "managementNetwork", auto.ConfigValue{Value: string(mn)})
		}
	}
	if (p.DeploymentNetwork != DeploymentNetwork{}) {
		if dn, err := json.Marshal(p.DeploymentNetwork); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "deploymentNetwork", auto.ConfigValue{Value: string(dn)})
		}
	}
	if (p.HelperVM != HelperVM{}) {
		if n, err := json.Marshal(p.HelperVM); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "helperVM", auto.ConfigValue{Value: string(n)})
		}
	}
	if p.DNSZoneName != "" {
		s.SetConfig(ctx, "dnsZoneName", auto.ConfigValue{Value: p.DNSZoneName})
	}
	if p.ReverseDNSZoneName != "" {
		s.SetConfig(ctx, "reverseDnsZoneName", auto.ConfigValue{Value: p.ReverseDNSZoneName})
	}
	if p.PublicRouter != "" {
		s.SetConfig(ctx, "publicRouter", auto.ConfigValue{Value: p.PublicRouter})
	}
	if p.KeypairFile.PublicKey != "" {
		s.SetConfig(ctx, "publicKeyFile", auto.ConfigValue{Value: p.KeypairFile.PublicKey})
	}
	if p.KeypairFile.PrivateKey != "" {
		s.SetConfig(ctx, "privateKeyFile", auto.ConfigValue{Value: p.KeypairFile.PrivateKey})
	}
	if p.PrivateNetworks != nil {
		if pn, err := json.Marshal(p.PrivateNetworks); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "privateNetworks", auto.ConfigValue{Value: string(pn)})
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
	if p.EsxiServerFlavor != "" {
		s.SetConfig(ctx, "esxiServerFlavor", auto.ConfigValue{Value: p.EsxiServerFlavor})
	}
	if p.Shares != nil {
		if n, err := json.Marshal(p.Shares); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "shares", auto.ConfigValue{Value: string(n)})
		}
	}

	// make reserved ip list
	if p.ReservedIPs == nil {
		p.ReservedIPs = []ReservedIP{}
	}
	if (p.Nsxt != Nsxt{}) {
		if n, err := json.Marshal(p.Nsxt); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "nsxt", auto.ConfigValue{Value: string(n)})
			p.ReservedIPs = append(p.ReservedIPs, ReservedIP{
				Hostname: p.Nsxt.Hostname,
				IP:       p.Nsxt.IP,
			})
		}
	}
	if p.NsxtManagers != nil {
		if n, err := json.Marshal(p.NsxtManagers); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "nsxtManagers", auto.ConfigValue{Value: string(n)})
			for _, m := range p.NsxtManagers {
				p.ReservedIPs = append(p.ReservedIPs, ReservedIP{
					Hostname: m.Hostname,
					IP:       m.IP,
				})
			}
		}
	}
	if (p.SDDCManager != SDDCManager{}) {
		if n, err := json.Marshal(p.SDDCManager); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "sddcManager", auto.ConfigValue{Value: string(n)})
			p.ReservedIPs = append(p.ReservedIPs, ReservedIP{
				Hostname: p.SDDCManager.Hostname,
				IP:       p.SDDCManager.IP,
			})
		}
	}
	if (p.VCenter != VCenter{}) {
		if n, err := json.Marshal(p.VCenter); err != nil {
			return err
		} else {
			s.SetConfig(ctx, "vcenter", auto.ConfigValue{Value: string(n)})
			p.ReservedIPs = append(p.ReservedIPs, ReservedIP{
				Hostname: p.VCenter.Hostname,
				IP:       p.VCenter.IP,
			})
		}
	}
	if n, err := json.Marshal(p.ReservedIPs); err != nil {
		return err
	} else {
		s.SetConfig(ctx, "reservedIPs", auto.ConfigValue{Value: string(n)})
	}
	return nil
}

func (s *Stack) Refresh(ctx context.Context) error {
	_, err := s.Stack.Refresh(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
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

func (s *Stack) Outputs(ctx context.Context) (auto.OutputMap, error) {
	o, err := s.Stack.Outputs(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := o["cloud-builder"]; ok {
		delete(o, "cloud-builder")
	}
	return o, nil
}

//
func (s *Stack) GetOutput(ctx context.Context, key string) (string, error) {
	o, err := s.Stack.Outputs(ctx)
	if err != nil {
		return "", err
	}
	if v, ok := o[key]; ok {
		if s, ok := v.Value.(string); ok {
			return s, nil
		} else {
			return "", fmt.Errorf("value error: %s", v.Value)
		}
	}
	return "", fmt.Errorf("key error: %s", key)
}

func (s *Stack) GetState() interface{} {
	return nil
}

func (s *Stack) GetError() error {
	return nil
}
