// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Use this data source to get the details of a running server
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/compute"
// 	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		_, err := compute.GetInstanceV2(ctx, &compute.GetInstanceV2Args{
// 			Id: "2ba26dc6-a12d-4889-8f25-794ea5bf4453",
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func GetInstanceV2(ctx *pulumi.Context, args *GetInstanceV2Args, opts ...pulumi.InvokeOption) (*GetInstanceV2Result, error) {
	var rv GetInstanceV2Result
	err := ctx.Invoke("openstack:compute/getInstanceV2:getInstanceV2", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getInstanceV2.
type GetInstanceV2Args struct {
	// The UUID of the instance
	Id string `pulumi:"id"`
	// An array of maps, detailed below.
	Networks []GetInstanceV2Network `pulumi:"networks"`
	Region   *string                `pulumi:"region"`
	// The user data added when the server was created.
	UserData *string `pulumi:"userData"`
}

// A collection of values returned by getInstanceV2.
type GetInstanceV2Result struct {
	// The first IPv4 address assigned to this server.
	AccessIpV4 string `pulumi:"accessIpV4"`
	// The first IPv6 address assigned to this server.
	AccessIpV6 string `pulumi:"accessIpV6"`
	// The availability zone of this server.
	AvailabilityZone string `pulumi:"availabilityZone"`
	// The flavor ID used to create the server.
	FlavorId string `pulumi:"flavorId"`
	// The flavor name used to create the server.
	FlavorName string `pulumi:"flavorName"`
	Id         string `pulumi:"id"`
	// The image ID used to create the server.
	ImageId string `pulumi:"imageId"`
	// The name of the key pair assigned to this server.
	KeyPair string `pulumi:"keyPair"`
	// A set of key/value pairs made available to the server.
	Metadata map[string]interface{} `pulumi:"metadata"`
	// The name of the network
	Name string `pulumi:"name"`
	// An array of maps, detailed below.
	Networks []GetInstanceV2Network `pulumi:"networks"`
	Region   string                 `pulumi:"region"`
	// An array of security group names associated with this server.
	SecurityGroups []string `pulumi:"securityGroups"`
	// A set of string tags assigned to this server.
	Tags []string `pulumi:"tags"`
	// The user data added when the server was created.
	UserData string `pulumi:"userData"`
}
