package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		image := conf.Require("imageName")
		floavor := conf.Require("flavorName")
		nodeUUID := conf.Require("nodeUUID")

		// Create Instance
		p, err := newInstance(ctx, floavor, image, nodeUUID)
		if err != nil {
			return err
		}

		// // Export the IP of the instance
		// ctx.Export("instanceIP", instance.AccessIpV4)
		ctx.Export("NetworkName", p.Network.Name)
		ctx.Export("SubnetName", p.Subnet.Name)
		ctx.Export("PortName", p.Port.Name)
		ctx.Export("InstanceIP", p.Instance.AccessIpV4)
		return nil
	})
}
