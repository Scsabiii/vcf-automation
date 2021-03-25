package main

import (
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

var props Props
var privateNetworks map[string]Network

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// private networks
		for _, props := range props.privateNetworks {
			network, err := networking.NewNetwork(ctx, props.Name, &networking.NetworkArgs{
				Name: pulumi.String("private-network-" + props.Name),
			})
			if err != nil {
				return err
			}
			subnet, err := networking.NewSubnet(ctx, "subnet-"+props.Name, &networking.SubnetArgs{
				Name:      pulumi.String("subnet-" + props.Name),
				NetworkId: network.ID(),
				IpVersion: pulumi.Int(4),
				Cidr:      pulumi.String(props.Cidr),
			})
			if err != nil {
				return err
			}
			if privateNetworks == nil {
				privateNetworks = make(map[string]Network)
			}
			privateNetworks[props.Name] = Network{network, subnet}
		}

		// private router
		privateRouter, err := networking.NewRouter(ctx, "management-private-router", &networking.RouterArgs{
			Name: pulumi.String("management-private-router"),
		})
		if err != nil {
			return err
		}
		for name, network := range privateNetworks {
			_, err = networking.NewRouterInterface(ctx, "router-interface-"+name, &networking.RouterInterfaceArgs{
				RouterId: privateRouter.ID(),
				SubnetId: network.Subnet.ID(),
			})
			if err != nil {
				return err
			}
		}

		//

		// privateNetworks, err := openstack.NewNetwork(
		// 	ctx,
		// 	props.privateNetwork.Name,
		// 	props.privateNetwork.Subnets...)
		// if err != nil {
		// 	return err
		// }
		// name := "router-interface-deployment"
		// _, err = networking.NewRouterInterface(ctx, name, &networking.RouterInterfaceArgs{
		// 	RouterId: privateRouter.ID(),
		// 	SubnetId: privateNetworks.Sutbnets["deployment"].ID(),
		// })
		// if err != nil {
		// 	return err
		// }
		// compute instance
		// privateSg, err := compute.NewSecGroup(ctx, "private-sg", &compute.SecGroupArgs{
		// 	Description: pulumi.String("security group for esxi server"),
		// 	Rules: compute.SecGroupRuleArray{
		// 		&compute.SecGroupRuleArgs{
		// 			Cidr:       pulumi.String("0.0.0.0/0"),
		// 			FromPort:   pulumi.Int(1),
		// 			ToPort:     pulumi.Int(65535),
		// 			IpProtocol: pulumi.String("tcp"),
		// 		},
		// 		&compute.SecGroupRuleArgs{
		// 			Cidr:       pulumi.String("0.0.0.0/0"),
		// 			FromPort:   pulumi.Int(1),
		// 			ToPort:     pulumi.Int(65535),
		// 			IpProtocol: pulumi.String("udp"),
		// 		},
		// 		&compute.SecGroupRuleArgs{
		// 			Cidr:       pulumi.String("0.0.0.0/0"),
		// 			FromPort:   pulumi.Int(0),
		// 			ToPort:     pulumi.Int(0),
		// 			IpProtocol: pulumi.String("icmp"),
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		// port1, err := networking.NewPort(ctx, "compute-port-1", &networking.PortArgs{
		// 	AdminStateUp: pulumi.Bool(true),
		// 	NetworkId:    privateNetworks.Network.ID(),
		// 	FixedIps: networking.PortFixedIpArray{
		// 		&networking.PortFixedIpArgs{
		// 			// IpAddress: pulumi.String("10.180"),
		// 			SubnetId: privateNetworks.Sutbnets["deployment"].ID(),
		// 		},
		// 	},
		// 	SecurityGroupIds: pulumi.StringArray{
		// 		privateSg.ID(),
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		// image := "ubuntu-18.04-amd64-vmware"
		// flavor := "m1.medium"
		// _, err = compute.NewInstance(ctx, "compute-1", &compute.InstanceArgs{
		// 	FlavorName: pulumi.StringPtr(flavor),
		// 	ImageName:  pulumi.StringPtr(image),
		// 	Networks: compute.InstanceNetworkArray{
		// 		&compute.InstanceNetworkArgs{
		// 			Port: port1.ID(),
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }

		// Export the IP of the instance
		ctx.Export("vmotionSubnetName", privateNetworks["vmotion"].Subnet.Name)
		ctx.Export("edgetepSubnetName", privateNetworks["edgetep"].Subnet.Name)
		ctx.Export("vmotionSubnetID", privateNetworks["vmotion"].Subnet.ID())
		ctx.Export("edgetepSubnetID", privateNetworks["edgetep"].Subnet.ID())
		return nil
	})
}

type Props struct {
	privateNetworks []NetworkProps
	publicNetworks  []NetworkProps
}

type NetworkProps struct {
	Name string
	Cidr string
}

type Network struct {
	*networking.Network
	*networking.Subnet
}

func init() {
	privateNetworkProps := []NetworkProps{{
		Name: "vmotion",
		Cidr: "10.180.0.0/24",
	}, {
		Name: "edgetep",
		Cidr: "10.180.1.0/24",
	}, {
		Name: "hosttep",
		Cidr: "10.180.2.0/24",
	}, {
		Name: "nfs",
		Cidr: "10.180.3.0/24",
	}, {
		Name: "vsan",
		Cidr: "10.180.4.0/24",
	}, {
		Name: "vsan-witness",
		Cidr: "10.180.5.0/24",
	}, {
		Name: "deployment",
		Cidr: "10.180.6.0/24",
	}}

	// publicNetwork := openstack.NetworkProp{
	// 	Name: "management-public",
	// 	Subnets: []openstack.SubnetProp{{
	// 		Name: "management-vcf01",
	// 		Cidr: "10.180.0.0/24",
	// 	}, {
	// 		Name: "edgetep",
	// 		Cidr: "10.180.1.0/24",
	// 	}},
	// }

	props = Props{privateNetworks: privateNetworkProps}
}
