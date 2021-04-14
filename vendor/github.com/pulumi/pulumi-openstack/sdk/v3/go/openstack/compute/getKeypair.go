// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package compute

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Use this data source to get the ID and public key of an OpenStack keypair.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/compute"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		_, err := compute.LookupKeypair(ctx, &compute.LookupKeypairArgs{
// 			Name: "sand",
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupKeypair(ctx *pulumi.Context, args *LookupKeypairArgs, opts ...pulumi.InvokeOption) (*LookupKeypairResult, error) {
	var rv LookupKeypairResult
	err := ctx.Invoke("openstack:compute/getKeypair:getKeypair", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getKeypair.
type LookupKeypairArgs struct {
	// The unique name of the keypair.
	Name string `pulumi:"name"`
	// The region in which to obtain the V2 Compute client.
	// If omitted, the `region` argument of the provider is used.
	Region *string `pulumi:"region"`
}

// A collection of values returned by getKeypair.
type LookupKeypairResult struct {
	// The fingerprint of the OpenSSH key.
	Fingerprint string `pulumi:"fingerprint"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// See Argument Reference above.
	Name string `pulumi:"name"`
	// The OpenSSH-formatted public key of the keypair.
	PublicKey string `pulumi:"publicKey"`
	// See Argument Reference above.
	Region string `pulumi:"region"`
}