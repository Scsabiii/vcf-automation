// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Use this data source to get the ID of an available OpenStack floating IP.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-openstack/sdk/v3/go/openstack/networking"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "192.168.0.4"
// 		_, err := networking.LookupFloatingIp(ctx, &networking.LookupFloatingIpArgs{
// 			Address: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupFloatingIp(ctx *pulumi.Context, args *LookupFloatingIpArgs, opts ...pulumi.InvokeOption) (*LookupFloatingIpResult, error) {
	var rv LookupFloatingIpResult
	err := ctx.Invoke("openstack:networking/getFloatingIp:getFloatingIp", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getFloatingIp.
type LookupFloatingIpArgs struct {
	// The IP address of the floating IP.
	Address *string `pulumi:"address"`
	// Human-readable description of the floating IP.
	Description *string `pulumi:"description"`
	// The specific IP address of the internal port which should be associated with the floating IP.
	FixedIp *string `pulumi:"fixedIp"`
	// The name of the pool from which the floating IP belongs to.
	Pool *string `pulumi:"pool"`
	// The ID of the port the floating IP is attached.
	PortId *string `pulumi:"portId"`
	// The region in which to obtain the V2 Neutron client.
	// A Neutron client is needed to retrieve floating IP ids. If omitted, the
	// `region` argument of the provider is used.
	Region *string `pulumi:"region"`
	// status of the floating IP (ACTIVE/DOWN).
	Status *string `pulumi:"status"`
	// The list of floating IP tags to filter.
	Tags []string `pulumi:"tags"`
	// The owner of the floating IP.
	TenantId *string `pulumi:"tenantId"`
}

// A collection of values returned by getFloatingIp.
type LookupFloatingIpResult struct {
	Address *string `pulumi:"address"`
	// A set of string tags applied on the floating IP.
	AllTags     []string `pulumi:"allTags"`
	Description *string  `pulumi:"description"`
	// The floating IP DNS domain. Available, when Neutron DNS
	// extension is enabled.
	DnsDomain string `pulumi:"dnsDomain"`
	// The floating IP DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName string  `pulumi:"dnsName"`
	FixedIp *string `pulumi:"fixedIp"`
	// The provider-assigned unique ID for this managed resource.
	Id       string   `pulumi:"id"`
	Pool     *string  `pulumi:"pool"`
	PortId   *string  `pulumi:"portId"`
	Region   *string  `pulumi:"region"`
	Status   *string  `pulumi:"status"`
	Tags     []string `pulumi:"tags"`
	TenantId *string  `pulumi:"tenantId"`
}