package main

import (
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/compute"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an OpenStack resource (Compute Instance)
		instance, err := compute.NewInstance(ctx, "test", &compute.InstanceArgs{
			FlavorName: pulumi.String("m1.small"),
			ImageName:  pulumi.String("ubuntu-18.04-amd64-vmware"),
			Networks: compute.InstanceNetworkArray{
				compute.InstanceNetworkArgs{
					Name: pulumi.String("d067954"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the IP of the instance
		ctx.Export("InstanceIP", instance.AccessIpV4)
		return nil
	})
}
