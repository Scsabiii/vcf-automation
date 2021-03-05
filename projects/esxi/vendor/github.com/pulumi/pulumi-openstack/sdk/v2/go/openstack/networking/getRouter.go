// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Use this data source to get the ID of an available OpenStack router.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
// 	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "router_1"
// 		_, err := networking.LookupRouter(ctx, &networking.LookupRouterArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupRouter(ctx *pulumi.Context, args *LookupRouterArgs, opts ...pulumi.InvokeOption) (*LookupRouterResult, error) {
	var rv LookupRouterResult
	err := ctx.Invoke("openstack:networking/getRouter:getRouter", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getRouter.
type LookupRouterArgs struct {
	// Administrative up/down status for the router (must be "true" or "false" if provided).
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// Human-readable description of the router.
	Description *string `pulumi:"description"`
	// Indicates whether or not to get a distributed router.
	Distributed *bool `pulumi:"distributed"`
	// The value that points out if the Source NAT is enabled on the router.
	EnableSnat *bool `pulumi:"enableSnat"`
	// The name of the router.
	Name *string `pulumi:"name"`
	// The region in which to obtain the V2 Neutron client.
	// A Neutron client is needed to retrieve router ids. If omitted, the
	// `region` argument of the provider is used.
	Region *string `pulumi:"region"`
	// The UUID of the router resource.
	RouterId *string `pulumi:"routerId"`
	// The status of the router (ACTIVE/DOWN).
	Status *string `pulumi:"status"`
	// The list of router tags to filter.
	Tags []string `pulumi:"tags"`
	// The owner of the router.
	TenantId *string `pulumi:"tenantId"`
}

// A collection of values returned by getRouter.
type LookupRouterResult struct {
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// The set of string tags applied on the router.
	AllTags []string `pulumi:"allTags"`
	// The availability zone that is used to make router resources highly available.
	AvailabilityZoneHints []string `pulumi:"availabilityZoneHints"`
	Description           *string  `pulumi:"description"`
	Distributed           *bool    `pulumi:"distributed"`
	// The value that points out if the Source NAT is enabled on the router.
	EnableSnat bool `pulumi:"enableSnat"`
	// The external fixed IPs of the router.
	ExternalFixedIps []GetRouterExternalFixedIp `pulumi:"externalFixedIps"`
	// The network UUID of an external gateway for the router.
	ExternalNetworkId string `pulumi:"externalNetworkId"`
	// The provider-assigned unique ID for this managed resource.
	Id       string   `pulumi:"id"`
	Name     *string  `pulumi:"name"`
	Region   *string  `pulumi:"region"`
	RouterId *string  `pulumi:"routerId"`
	Status   *string  `pulumi:"status"`
	Tags     []string `pulumi:"tags"`
	TenantId *string  `pulumi:"tenantId"`
}