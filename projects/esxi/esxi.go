package main

import (
	"fmt"
	"strconv"

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
	Image  string
	Flavor string
	UUID   string
}

type ShareProps struct {
	Name string
	Size int
}

func newEsxiStack(ctx *pulumi.Context) error {
	// Read stack configuration
	conf := config.New(ctx, "")
	prefix := conf.Get("resourcePrefix")
	np := &NodeProps{
		Image:  conf.Require("imageName"),
		Flavor: conf.Require("flavorName"),
		UUID:   conf.Require("nodeUUID"),
	}

	// Create Instance
	esxiNetwork, err := newEsxiNetwork(ctx, prefix, "192.168.199.0/24")
	if err != nil {
		return err
	}
	sg, err := newEsxiSecGroup(ctx, prefix)
	if err != nil {
		return err
	}
	port, err := newEsxiPort(ctx, prefix, "192.168.199.10", esxiNetwork, sg)
	if err != nil {
		return err
	}
	// instance, err := newMetalInstance(ctx, np, port)
	instance, err := newComputeInstance(ctx, prefix, np, port)
	if err != nil {
		return err
	}

	// Create NFS Shares
	storageNetwork, err := newStroageNetwork(ctx, prefix, "192.168.200.0/24")
	if err != nil {
		return err
	}
	shareNetwork, err := newShareNetwork(ctx, prefix, storageNetwork)
	if err != nil {
		return err
	}
	shareProps, err := readShareProps(conf)
	if err != nil {
		return err
	}
	shares := make([]*sharedfilesystem.Share, 0)
	for _, s := range shareProps {
		nfs, err := newNFSShare(ctx, prefix, s, shareNetwork)
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

	ctx.Export("EsxiInstanceName", instance.Name)
	ctx.Export("EsxiInstanceID", instance.ID())
	ctx.Export("EsxiInstanceIP", instance.AccessIpV4)

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

func newMetalInstance(ctx *pulumi.Context, node *NodeProps,
	n *Network, sg *compute.SecGroup, ipAddr string) (*compute.Instance, error) {
	port, err := networking.NewPort(ctx, "esxi-port", &networking.PortArgs{
		AdminStateUp: pulumi.Bool(true),
		NetworkId:    n.Network.ID(),
		FixedIps: networking.PortFixedIpArray{
			&networking.PortFixedIpArgs{
				IpAddress: pulumi.String(ipAddr),
				SubnetId:  n.Subnet.ID(),
			},
		},
		SecurityGroupIds: pulumi.StringArray{
			sg.ID(),
		},
	})
	if err != nil {
		return nil, err
	}
	return compute.NewInstance(ctx, "esxi-metal-instance", &compute.InstanceArgs{
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

func newEsxiPort(ctx *pulumi.Context, prefix, ipAddr string, n *Network, sg *compute.SecGroup) (*networking.Port, error) {
	name := fmt.Sprintf("%s-port", prefix)
	return networking.NewPort(ctx, name, &networking.PortArgs{
		AdminStateUp: pulumi.Bool(true),
		NetworkId:    n.Network.ID(),
		FixedIps: networking.PortFixedIpArray{
			&networking.PortFixedIpArgs{
				IpAddress: pulumi.String(ipAddr),
				SubnetId:  n.Subnet.ID(),
			},
		},
		SecurityGroupIds: pulumi.StringArray{
			sg.ID(),
		},
	})
}

func newComputeInstance(ctx *pulumi.Context, prefix string, node *NodeProps, port *networking.Port) (*compute.Instance, error) {
	name := fmt.Sprintf("%s-instance", prefix)
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

func readShareProps(conf *config.Config) ([]*ShareProps, error) {
	ns, err := strconv.Atoi(conf.Require("numShares"))
	if err != nil {
		return nil, err
	}
	sp := make([]*ShareProps, ns)
	for i := 0; i < ns; i++ {
		size, err := strconv.Atoi(conf.Require(fmt.Sprintf("share%02dSize", i)))
		if err != nil {
			return nil, err
		}
		sp[i] = &ShareProps{
			Name: conf.Require(fmt.Sprintf("share%02dName", i)),
			Size: size,
		}
	}
	return sp, nil
}

func parseNodeProps(conf *config.Config) ([]*NodeProps, error) {
	nn, err := strconv.Atoi(conf.Require("numNodes"))
	if err != nil {
		return nil, err
	}
	np := make([]*NodeProps, nn)
	return np, nil
}
