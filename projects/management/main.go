package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/sapcc/avocado-automation/pkg/openstack"
)

var props Props

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		privateNetwork, err := openstack.NewNetwork(ctx, "management", props.managementNetworkProp.Subnets...)
		if err != nil {
			return err
		}

		// // Create an OpenStack resource (Compute Instance)
		// instance, err := compute.NewInstance(ctx, "test", &compute.InstanceArgs{
		// 	FlavorName: pulumi.String("s1-2"),
		// 	ImageName:  pulumi.String("Ubuntu 16.04"),
		// })

		// Export the IP of the instance
		ctx.Export("vmotionSubnetName", privateNetwork.Sutbnets[0].Name)
		ctx.Export("vmotionSubnetID", privateNetwork.Sutbnets[0].ID())
		ctx.Export("edgetepSubnetName", privateNetwork.Sutbnets[1].Name)
		ctx.Export("edgetepSubnetID", privateNetwork.Sutbnets[1].ID())
		return nil
	})
}

type Props struct {
	managementNetworkProp openstack.NetworkProp
}

func init() {
	managementNetworkProp := openstack.NetworkProp{
		Name: "management",
		Subnets: []openstack.SubnetProp{{
			Name: "vmotion",
			Cidr: "10.180.0.0/24",
		}, {
			Name: "edgetep",
			Cidr: "10.180.0.1/24",
		}},
	}
	props = Props{managementNetworkProp}
}
