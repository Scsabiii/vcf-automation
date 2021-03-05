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

package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/compute"
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/sharedfilesystem"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

type Network struct {
	Network *networking.Network
	Subnet  *networking.Subnet
}

type EsxiInstance struct {
	Instance *compute.Instance
	SecGroup *compute.SecGroup
}

type NodeProps struct {
	ID     int    `json:"id"`
	Image  string `json:"image"`
	Flavor string `json:"flavor"`
	UUID   string `json:"uuid"`
	IP     string `json:"ip"`
}

type ShareProps struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

func newEsxiStack(ctx *pulumi.Context) error {
	// Read stack configuration
	conf := config.New(ctx, "")
	prefix := conf.Get("resourcePrefix")
	nodeSubent := conf.Get("nodeSubnet")
	storageSubnet := conf.Get("storageSubnet")

	nodePropsStr := conf.Get("nodes")
	nodeProps := make([]NodeProps, 0)
	err := json.Unmarshal([]byte(nodePropsStr), &nodeProps)
	if err != nil {
		return err
	}

	sharePropsStr := conf.Get("shares")
	shareProps := make([]ShareProps, 0)
	err = json.Unmarshal([]byte(sharePropsStr), &shareProps)
	if err != nil {
		return err
	}

	// Create Instance
	esxiNetwork, err := newEsxiNetwork(ctx, prefix, nodeSubent)
	if err != nil {
		return err
	}
	sg, err := newEsxiSecGroup(ctx, prefix)
	if err != nil {
		return err
	}
	nodes := make([]*compute.Instance, 0)
	for _, np := range nodeProps {
		port, err := newEsxiPort(ctx, prefix, &np, esxiNetwork, sg)
		if err != nil {
			return err
		}
		instance, err := newComputeInstance(ctx, prefix, &np, port)
		if err != nil {
			return err
		}
		nodes = append(nodes, instance)
	}

	// Create NFS Shares
	storageNetwork, err := newStroageNetwork(ctx, prefix, storageSubnet)
	if err != nil {
		return err
	}
	shareNetwork, err := newShareNetwork(ctx, prefix, storageNetwork)
	if err != nil {
		return err
	}
	shares := make([]*sharedfilesystem.Share, 0)
	for _, s := range shareProps {
		nfs, err := newNFSShare(ctx, prefix, &s, shareNetwork)
		if err != nil {
			return err
		}
		shares = append(shares, nfs)
	}

	// Export
	ctx.Export("EsxiNetworkName", esxiNetwork.Network.Name)
	ctx.Export("EsxiNetworkID", esxiNetwork.Network.ID())
	ctx.Export("EsxiSubnetName", esxiNetwork.Subnet.Name)
	ctx.Export("EsxiSubnetID", esxiNetwork.Subnet.ID())
	ctx.Export("EsxiSubnetCIDR", esxiNetwork.Subnet.Cidr)
	ctx.Export("StorageNetworkName", storageNetwork.Network.Name)
	ctx.Export("StorageNetworkID", storageNetwork.Network.ID())
	ctx.Export("StorageSubnetName", storageNetwork.Subnet.Name)
	ctx.Export("StorageSubnetID", storageNetwork.Subnet.ID())
	ctx.Export("StorageSubnetCIDR", storageNetwork.Subnet.Cidr)

	for i, n := range nodes {
		ctx.Export(fmt.Sprintf("EsxiInstance%02dName", i), n.Name)
		ctx.Export(fmt.Sprintf("EsxiInstance%02dID", i), n.ID())
		ctx.Export(fmt.Sprintf("EsxiInstance%02dIP", i), n.AccessIpV4)
	}
	for i, s := range shares {
		ctx.Export(fmt.Sprintf("Share%02dName", i), s.Name)
		ctx.Export(fmt.Sprintf("Share%02dID", i), s.ID())
	}
	return nil
}

func newStroageNetwork(ctx *pulumi.Context, prefix, cidr string) (sn *Network, err error) {
	name := fmt.Sprintf("%s-storage-network", prefix)
	network, err := networking.NewNetwork(ctx, name, &networking.NetworkArgs{
		AdminStateUp: pulumi.Bool(true),
	})
	if err != nil {
		return
	}
	name = fmt.Sprintf("%s-storage-subnet", prefix)
	subnet, err := networking.NewSubnet(ctx, name, &networking.SubnetArgs{
		Cidr:      pulumi.String(cidr),
		IpVersion: pulumi.Int(4),
		NetworkId: network.ID(),
	})
	if err != nil {
		return
	}
	return &Network{network, subnet}, nil
}

func newEsxiNetwork(ctx *pulumi.Context, prefix, subnetCidr string) (n *Network, err error) {
	name := fmt.Sprintf("%s-network", prefix)
	network, err := networking.NewNetwork(ctx, name, &networking.NetworkArgs{
		AdminStateUp: pulumi.Bool(true),
	})
	if err != nil {
		return
	}
	name = fmt.Sprintf("%s-subnet", prefix)
	subnet, err := networking.NewSubnet(ctx, name, &networking.SubnetArgs{
		Cidr:      pulumi.String(subnetCidr),
		IpVersion: pulumi.Int(4),
		NetworkId: network.ID(),
	})
	if err != nil {
		return
	}
	return &Network{network, subnet}, nil
}

func newEsxiSecGroup(ctx *pulumi.Context, prefix string) (*compute.SecGroup, error) {
	name := fmt.Sprintf("%s-sg", prefix)
	return compute.NewSecGroup(ctx, name, &compute.SecGroupArgs{
		Description: pulumi.String("security group for esxi server"),
		Rules: compute.SecGroupRuleArray{
			&compute.SecGroupRuleArgs{
				Cidr:       pulumi.String("0.0.0.0/0"),
				FromPort:   pulumi.Int(1),
				ToPort:     pulumi.Int(65535),
				IpProtocol: pulumi.String("tcp"),
			},
			&compute.SecGroupRuleArgs{
				Cidr:       pulumi.String("0.0.0.0/0"),
				FromPort:   pulumi.Int(1),
				ToPort:     pulumi.Int(65535),
				IpProtocol: pulumi.String("udp"),
			},
			&compute.SecGroupRuleArgs{
				Cidr:       pulumi.String("0.0.0.0/0"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				IpProtocol: pulumi.String("icmp"),
			},
		},
	})
}

func newEsxiPort(ctx *pulumi.Context, prefix string, node *NodeProps, n *Network, sg *compute.SecGroup) (*networking.Port, error) {
	name := fmt.Sprintf("%s-port-%d", prefix, node.ID)
	return networking.NewPort(ctx, name, &networking.PortArgs{
		AdminStateUp: pulumi.Bool(true),
		NetworkId:    n.Network.ID(),
		FixedIps: networking.PortFixedIpArray{
			&networking.PortFixedIpArgs{
				IpAddress: pulumi.String(node.IP),
				SubnetId:  n.Subnet.ID(),
			},
		},
		SecurityGroupIds: pulumi.StringArray{
			sg.ID(),
		},
	})
}

func newComputeInstance(ctx *pulumi.Context, prefix string, node *NodeProps, port *networking.Port) (*compute.Instance, error) {
	name := fmt.Sprintf("%s-instance-%d", prefix, node.ID)
	return compute.NewInstance(ctx, name, &compute.InstanceArgs{
		FlavorName: pulumi.String(node.Flavor),
		ImageName:  pulumi.String(node.Image),
		Networks: compute.InstanceNetworkArray{
			&compute.InstanceNetworkArgs{
				Port: port.ID(),
			},
		},
	})
}

func newEsxiInstance(ctx *pulumi.Context, prefix string, node *NodeProps, port *networking.Port) (*compute.Instance, error) {
	name := fmt.Sprintf("%s-esxi-instance", prefix)
	return compute.NewInstance(ctx, name, &compute.InstanceArgs{
		FlavorName:            pulumi.String(node.Flavor),
		ImageName:             pulumi.String(node.Image),
		AvailabilityZoneHints: pulumi.String(fmt.Sprintf("::%s", node.UUID)),
		Networks: compute.InstanceNetworkArray{
			&compute.InstanceNetworkArgs{
				Port: port.ID(),
			},
		},
	})
}
func newShareNetwork(ctx *pulumi.Context, prefix string, n *Network) (*sharedfilesystem.ShareNetwork, error) {
	name := fmt.Sprintf("%s-share-network", prefix)
	return sharedfilesystem.NewShareNetwork(ctx, name, &sharedfilesystem.ShareNetworkArgs{
		NeutronNetId:    n.Network.ID(),
		NeutronSubnetId: n.Subnet.ID(),
	})
}

func newNFSShare(ctx *pulumi.Context, prefix string, p *ShareProps, sn *sharedfilesystem.ShareNetwork) (*sharedfilesystem.Share, error) {
	shareName := fmt.Sprintf("%s-%s", prefix, p.Name)
	return sharedfilesystem.NewShare(ctx, shareName, &sharedfilesystem.ShareArgs{
		Description:    pulumi.String("test share description"),
		ShareNetworkId: sn.ID(),
		ShareProto:     pulumi.String("NFS"),
		Size:           pulumi.Int(p.Size),
	})
}
