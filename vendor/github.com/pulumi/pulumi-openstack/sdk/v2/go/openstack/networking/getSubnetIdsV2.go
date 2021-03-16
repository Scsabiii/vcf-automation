// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Use this data source to get a list of Openstack Subnet IDs matching the
// specified criteria.
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
// 		opt0 := "public"
// 		_, err := networking.GetSubnetIdsV2(ctx, &networking.GetSubnetIdsV2Args{
// 			NameRegex: &opt0,
// 			Tags: []string{
// 				"public",
// 			},
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func GetSubnetIdsV2(ctx *pulumi.Context, args *GetSubnetIdsV2Args, opts ...pulumi.InvokeOption) (*GetSubnetIdsV2Result, error) {
	var rv GetSubnetIdsV2Result
	err := ctx.Invoke("openstack:networking/getSubnetIdsV2:getSubnetIdsV2", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getSubnetIdsV2.
type GetSubnetIdsV2Args struct {
	// The CIDR of the subnet.
	Cidr *string `pulumi:"cidr"`
	// Human-readable description of the subnet.
	Description *string `pulumi:"description"`
	// If the subnet has DHCP enabled.
	DhcpEnabled *bool `pulumi:"dhcpEnabled"`
	// The IP of the subnet's gateway.
	GatewayIp *string `pulumi:"gatewayIp"`
	// The IP version of the subnet (either 4 or 6).
	IpVersion *int `pulumi:"ipVersion"`
	// The IPv6 address mode. Valid values are
	// `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.
	Ipv6AddressMode *string `pulumi:"ipv6AddressMode"`
	// The IPv6 Router Advertisement mode. Valid values
	// are `dhcpv6-stateful`, `dhcpv6-stateless`, or `slaac`.
	Ipv6RaMode *string `pulumi:"ipv6RaMode"`
	// The name of the subnet.
	Name      *string `pulumi:"name"`
	NameRegex *string `pulumi:"nameRegex"`
	// The ID of the network the subnet belongs to.
	NetworkId *string `pulumi:"networkId"`
	// The region in which to obtain the V2 Neutron client.
	// A Neutron client is needed to retrieve subnet ids. If omitted, the
	// `region` argument of the provider is used.
	Region *string `pulumi:"region"`
	// Order the results in either `asc` or `desc`.
	// Defaults to none.
	SortDirection *string `pulumi:"sortDirection"`
	// Sort subnets based on a certain key. Defaults to none.
	SortKey *string `pulumi:"sortKey"`
	// The ID of the subnetpool associated with the subnet.
	SubnetpoolId *string `pulumi:"subnetpoolId"`
	// The list of subnet tags to filter.
	Tags []string `pulumi:"tags"`
	// The owner of the subnet.
	TenantId *string `pulumi:"tenantId"`
}

// A collection of values returned by getSubnetIdsV2.
type GetSubnetIdsV2Result struct {
	Cidr        *string `pulumi:"cidr"`
	Description *string `pulumi:"description"`
	DhcpEnabled *bool   `pulumi:"dhcpEnabled"`
	GatewayIp   *string `pulumi:"gatewayIp"`
	// The provider-assigned unique ID for this managed resource.
	Id              string   `pulumi:"id"`
	Ids             []string `pulumi:"ids"`
	IpVersion       *int     `pulumi:"ipVersion"`
	Ipv6AddressMode *string  `pulumi:"ipv6AddressMode"`
	Ipv6RaMode      string   `pulumi:"ipv6RaMode"`
	Name            *string  `pulumi:"name"`
	NameRegex       *string  `pulumi:"nameRegex"`
	NetworkId       *string  `pulumi:"networkId"`
	Region          string   `pulumi:"region"`
	SortDirection   *string  `pulumi:"sortDirection"`
	SortKey         *string  `pulumi:"sortKey"`
	SubnetpoolId    *string  `pulumi:"subnetpoolId"`
	Tags            []string `pulumi:"tags"`
	TenantId        *string  `pulumi:"tenantId"`
}
