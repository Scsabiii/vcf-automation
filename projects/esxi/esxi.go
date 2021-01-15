package main

import (
	"fmt"

	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/compute"
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Provision struct {
	Network  *networking.Network
	Subnet   *networking.Subnet
	Port     *networking.Port
	SecGroup *compute.SecGroup
	Instance *compute.Instance
}

func newInstance(ctx *pulumi.Context, floavor, image, nodeUUID string) (*Provision, error) {
	network1, err := networking.NewNetwork(ctx, "network-esxi", &networking.NetworkArgs{
		AdminStateUp: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	subnet1, err := networking.NewSubnet(ctx, "subnet-esxi", &networking.SubnetArgs{
		Cidr:      pulumi.String("192.168.199.0/24"),
		IpVersion: pulumi.Int(4),
		NetworkId: network1.ID(),
	})
	if err != nil {
		return nil, err
	}
	secgroup1, err := compute.NewSecGroup(ctx, "sg-esxi", &compute.SecGroupArgs{
		Description: pulumi.String("cci esxi security group"),
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
	if err != nil {
		return nil, err
	}
	port1, err := networking.NewPort(ctx, "port1", &networking.PortArgs{
		AdminStateUp: pulumi.Bool(true),
		FixedIps: networking.PortFixedIpArray{
			&networking.PortFixedIpArgs{
				IpAddress: pulumi.String("192.168.199.10"),
				SubnetId:  subnet1.ID(),
			},
		},
		NetworkId: network1.ID(),
		SecurityGroupIds: pulumi.StringArray{
			secgroup1.ID(),
		},
	})
	if err != nil {
		return nil, err
	}
	instance1, err := compute.NewInstance(ctx, "instance1", &compute.InstanceArgs{
		FlavorName:            pulumi.String(floavor),
		ImageName:             pulumi.String(image),
		AvailabilityZoneHints: pulumi.String(fmt.Sprintf("::%s", nodeUUID)),
		Networks: compute.InstanceNetworkArray{
			&compute.InstanceNetworkArgs{
				Port: port1.ID(),
			},
		},
		SecurityGroups: pulumi.StringArray{
			secgroup1.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	p := Provision{
		Network:  network1,
		Subnet:   subnet1,
		Port:     port1,
		Instance: instance1,
	}
	return &p, nil
}

func newSecGroup(ctx *pulumi.Context) (*compute.SecGroup, error) {
	sg, err := compute.NewSecGroup(ctx, "sg-esxi", &compute.SecGroupArgs{
		Description: pulumi.String("cci esxi security group"),
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
	return sg, err
}
