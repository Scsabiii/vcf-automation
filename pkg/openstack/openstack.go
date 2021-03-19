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

package openstack

import (
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Network struct {
	Network  *networking.Network
	Sutbnets map[string]*networking.Subnet
}

type NetworkProp struct {
	Name    string
	Subnets []SubnetProp
}

type SubnetProp struct {
	Name string
	Cidr string
}

// NewNetwork creates network and subnets. Network physical name is same as the
// given name, subent names are auto generated
func NewNetwork(ctx *pulumi.Context, networkName string, subnetProps ...SubnetProp) (*Network, error) {
	network, err := networking.NewNetwork(ctx, networkName, &networking.NetworkArgs{
		AdminStateUp: pulumi.Bool(true),
		Name:         pulumi.String(networkName),
	})
	if err != nil {
		return nil, err
	}
	subnets := make(map[string]*networking.Subnet, 0)
	for _, prop := range subnetProps {
		s, err := networking.NewSubnet(ctx, prop.Name, &networking.SubnetArgs{
			Name:      pulumi.String(prop.Name),
			Cidr:      pulumi.String(prop.Cidr),
			IpVersion: pulumi.Int(4),
			NetworkId: network.ID(),
		})
		if err != nil {
			return nil, err
		}
		subnets[prop.Name] = s
	}
	return &Network{network, subnets}, nil
}
