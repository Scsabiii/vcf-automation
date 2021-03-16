// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Manages a V2 port resource within OpenStack.
//
// ## Example Usage
// ### Simple port
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
// 		network1, err := networking.NewNetwork(ctx, "network1", &networking.NetworkArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = networking.NewPort(ctx, "port1", &networking.PortArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 			NetworkId:    network1.ID(),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
// ### Port with physical binding information
//
// ```go
// package main
//
// import (
// 	"fmt"
//
// 	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/networking"
// 	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		network1, err := networking.NewNetwork(ctx, "network1", &networking.NetworkArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = networking.NewPort(ctx, "port1", &networking.PortArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 			Binding: &networking.PortBindingArgs{
// 				HostId:   pulumi.String("b080b9cf-46e0-4ce8-ad47-0fd4accc872b"),
// 				Profile:  pulumi.String(fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v", "{\n", "  \"local_link_information\": [\n", "    {\n", "      \"switch_info\": \"info1\",\n", "      \"port_id\": \"Ethernet3/4\",\n", "      \"switch_id\": \"12:34:56:78:9A:BC\"\n", "    },\n", "    {\n", "      \"switch_info\": \"info2\",\n", "      \"port_id\": \"Ethernet3/4\",\n", "      \"switch_id\": \"12:34:56:78:9A:BD\"\n", "    }\n", "  ],\n", "  \"vlan_type\": \"allowed\"\n", "}\n", "\n")),
// 				VnicType: pulumi.String("baremetal"),
// 			},
// 			DeviceId:    pulumi.String("cdf70fcf-c161-4f24-9c70-96b3f5a54b71"),
// 			DeviceOwner: pulumi.String("baremetal:none"),
// 			NetworkId:   network1.ID(),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
// ## Notes
//
// ### Ports and Instances
//
// There are some notes to consider when connecting Instances to networks using
// Ports. Please see the `compute.Instance` documentation for further
// documentation.
//
// ## Import
//
// Ports can be imported using the `id`, e.g.
//
// ```sh
//  $ pulumi import openstack:networking/port:Port port_1 eae26a3e-1c33-4cc1-9c31-0cd729c438a1
// ```
type Port struct {
	pulumi.CustomResourceState

	// Administrative up/down status for the port
	// (must be `true` or `false` if provided). Changing this updates the
	// `adminStateUp` of an existing port.
	AdminStateUp pulumi.BoolOutput `pulumi:"adminStateUp"`
	// The collection of Fixed IP addresses on the port in the
	// order returned by the Network v2 API.
	AllFixedIps pulumi.StringArrayOutput `pulumi:"allFixedIps"`
	// The collection of Security Group IDs on the port
	// which have been explicitly and implicitly added.
	AllSecurityGroupIds pulumi.StringArrayOutput `pulumi:"allSecurityGroupIds"`
	// The collection of tags assigned on the port, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayOutput `pulumi:"allTags"`
	// An IP/MAC Address pair of additional IP
	// addresses that can be active on this port. The structure is described
	// below.
	AllowedAddressPairs PortAllowedAddressPairArrayOutput `pulumi:"allowedAddressPairs"`
	// The port binding allows to specify binding information
	// for the port. The structure is described below.
	Binding PortBindingOutput `pulumi:"binding"`
	// Human-readable description of the port. Changing
	// this updates the `description` of an existing port.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// The ID of the device attached to the port. Changing this
	// creates a new port.
	DeviceId pulumi.StringOutput `pulumi:"deviceId"`
	// The device owner of the port. Changing this creates
	// a new port.
	DeviceOwner pulumi.StringOutput `pulumi:"deviceOwner"`
	// The list of maps representing port DNS assignments.
	DnsAssignments pulumi.MapArrayOutput `pulumi:"dnsAssignments"`
	// The port DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName pulumi.StringOutput `pulumi:"dnsName"`
	// An extra DHCP option that needs to be configured
	// on the port. The structure is described below. Can be specified multiple
	// times.
	ExtraDhcpOptions PortExtraDhcpOptionArrayOutput `pulumi:"extraDhcpOptions"`
	// An array of desired IPs for
	// this port. The structure is described below.
	FixedIps PortFixedIpArrayOutput `pulumi:"fixedIps"`
	// The additional MAC address.
	MacAddress pulumi.StringOutput `pulumi:"macAddress"`
	// Name of the DHCP option.
	Name pulumi.StringOutput `pulumi:"name"`
	// The ID of the network to attach the port to. Changing
	// this creates a new port.
	NetworkId pulumi.StringOutput `pulumi:"networkId"`
	// Create a port with no fixed
	// IP address. This will also remove any fixed IPs previously set on a port. `true`
	// is the only valid value for this argument.
	NoFixedIp pulumi.BoolPtrOutput `pulumi:"noFixedIp"`
	// If set to
	// `true`, then no security groups are applied to the port. If set to `false` and
	// no `securityGroupIds` are specified, then the port will yield to the default
	// behavior of the Networking service, which is to usually apply the "default"
	// security group.
	NoSecurityGroups pulumi.BoolPtrOutput `pulumi:"noSecurityGroups"`
	// Whether to explicitly enable or disable
	// port security on the port. Port Security is usually enabled by default, so
	// omitting argument will usually result in a value of `true`. Setting this
	// explicitly to `false` will disable port security. In order to disable port
	// security, the port must not have any security groups. Valid values are `true`
	// and `false`.
	PortSecurityEnabled pulumi.BoolOutput `pulumi:"portSecurityEnabled"`
	// Reference to the associated QoS policy.
	QosPolicyId pulumi.StringOutput `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a port. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// port.
	Region pulumi.StringOutput `pulumi:"region"`
	// A list
	// of security group IDs to apply to the port. The security groups must be
	// specified by ID and not name (as opposed to how they are configured with
	// the Compute Instance).
	SecurityGroupIds pulumi.StringArrayOutput `pulumi:"securityGroupIds"`
	// A set of string tags for the port.
	Tags pulumi.StringArrayOutput `pulumi:"tags"`
	// The owner of the port. Required if admin wants
	// to create a port for another tenant. Changing this creates a new port.
	TenantId pulumi.StringOutput `pulumi:"tenantId"`
	// Map of additional options.
	ValueSpecs pulumi.MapOutput `pulumi:"valueSpecs"`
}

// NewPort registers a new resource with the given unique name, arguments, and options.
func NewPort(ctx *pulumi.Context,
	name string, args *PortArgs, opts ...pulumi.ResourceOption) (*Port, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.NetworkId == nil {
		return nil, errors.New("invalid value for required argument 'NetworkId'")
	}
	var resource Port
	err := ctx.RegisterResource("openstack:networking/port:Port", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetPort gets an existing Port resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetPort(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *PortState, opts ...pulumi.ResourceOption) (*Port, error) {
	var resource Port
	err := ctx.ReadResource("openstack:networking/port:Port", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Port resources.
type portState struct {
	// Administrative up/down status for the port
	// (must be `true` or `false` if provided). Changing this updates the
	// `adminStateUp` of an existing port.
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// The collection of Fixed IP addresses on the port in the
	// order returned by the Network v2 API.
	AllFixedIps []string `pulumi:"allFixedIps"`
	// The collection of Security Group IDs on the port
	// which have been explicitly and implicitly added.
	AllSecurityGroupIds []string `pulumi:"allSecurityGroupIds"`
	// The collection of tags assigned on the port, which have been
	// explicitly and implicitly added.
	AllTags []string `pulumi:"allTags"`
	// An IP/MAC Address pair of additional IP
	// addresses that can be active on this port. The structure is described
	// below.
	AllowedAddressPairs []PortAllowedAddressPair `pulumi:"allowedAddressPairs"`
	// The port binding allows to specify binding information
	// for the port. The structure is described below.
	Binding *PortBinding `pulumi:"binding"`
	// Human-readable description of the port. Changing
	// this updates the `description` of an existing port.
	Description *string `pulumi:"description"`
	// The ID of the device attached to the port. Changing this
	// creates a new port.
	DeviceId *string `pulumi:"deviceId"`
	// The device owner of the port. Changing this creates
	// a new port.
	DeviceOwner *string `pulumi:"deviceOwner"`
	// The list of maps representing port DNS assignments.
	DnsAssignments []map[string]interface{} `pulumi:"dnsAssignments"`
	// The port DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName *string `pulumi:"dnsName"`
	// An extra DHCP option that needs to be configured
	// on the port. The structure is described below. Can be specified multiple
	// times.
	ExtraDhcpOptions []PortExtraDhcpOption `pulumi:"extraDhcpOptions"`
	// An array of desired IPs for
	// this port. The structure is described below.
	FixedIps []PortFixedIp `pulumi:"fixedIps"`
	// The additional MAC address.
	MacAddress *string `pulumi:"macAddress"`
	// Name of the DHCP option.
	Name *string `pulumi:"name"`
	// The ID of the network to attach the port to. Changing
	// this creates a new port.
	NetworkId *string `pulumi:"networkId"`
	// Create a port with no fixed
	// IP address. This will also remove any fixed IPs previously set on a port. `true`
	// is the only valid value for this argument.
	NoFixedIp *bool `pulumi:"noFixedIp"`
	// If set to
	// `true`, then no security groups are applied to the port. If set to `false` and
	// no `securityGroupIds` are specified, then the port will yield to the default
	// behavior of the Networking service, which is to usually apply the "default"
	// security group.
	NoSecurityGroups *bool `pulumi:"noSecurityGroups"`
	// Whether to explicitly enable or disable
	// port security on the port. Port Security is usually enabled by default, so
	// omitting argument will usually result in a value of `true`. Setting this
	// explicitly to `false` will disable port security. In order to disable port
	// security, the port must not have any security groups. Valid values are `true`
	// and `false`.
	PortSecurityEnabled *bool `pulumi:"portSecurityEnabled"`
	// Reference to the associated QoS policy.
	QosPolicyId *string `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a port. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// port.
	Region *string `pulumi:"region"`
	// A list
	// of security group IDs to apply to the port. The security groups must be
	// specified by ID and not name (as opposed to how they are configured with
	// the Compute Instance).
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// A set of string tags for the port.
	Tags []string `pulumi:"tags"`
	// The owner of the port. Required if admin wants
	// to create a port for another tenant. Changing this creates a new port.
	TenantId *string `pulumi:"tenantId"`
	// Map of additional options.
	ValueSpecs map[string]interface{} `pulumi:"valueSpecs"`
}

type PortState struct {
	// Administrative up/down status for the port
	// (must be `true` or `false` if provided). Changing this updates the
	// `adminStateUp` of an existing port.
	AdminStateUp pulumi.BoolPtrInput
	// The collection of Fixed IP addresses on the port in the
	// order returned by the Network v2 API.
	AllFixedIps pulumi.StringArrayInput
	// The collection of Security Group IDs on the port
	// which have been explicitly and implicitly added.
	AllSecurityGroupIds pulumi.StringArrayInput
	// The collection of tags assigned on the port, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayInput
	// An IP/MAC Address pair of additional IP
	// addresses that can be active on this port. The structure is described
	// below.
	AllowedAddressPairs PortAllowedAddressPairArrayInput
	// The port binding allows to specify binding information
	// for the port. The structure is described below.
	Binding PortBindingPtrInput
	// Human-readable description of the port. Changing
	// this updates the `description` of an existing port.
	Description pulumi.StringPtrInput
	// The ID of the device attached to the port. Changing this
	// creates a new port.
	DeviceId pulumi.StringPtrInput
	// The device owner of the port. Changing this creates
	// a new port.
	DeviceOwner pulumi.StringPtrInput
	// The list of maps representing port DNS assignments.
	DnsAssignments pulumi.MapArrayInput
	// The port DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName pulumi.StringPtrInput
	// An extra DHCP option that needs to be configured
	// on the port. The structure is described below. Can be specified multiple
	// times.
	ExtraDhcpOptions PortExtraDhcpOptionArrayInput
	// An array of desired IPs for
	// this port. The structure is described below.
	FixedIps PortFixedIpArrayInput
	// The additional MAC address.
	MacAddress pulumi.StringPtrInput
	// Name of the DHCP option.
	Name pulumi.StringPtrInput
	// The ID of the network to attach the port to. Changing
	// this creates a new port.
	NetworkId pulumi.StringPtrInput
	// Create a port with no fixed
	// IP address. This will also remove any fixed IPs previously set on a port. `true`
	// is the only valid value for this argument.
	NoFixedIp pulumi.BoolPtrInput
	// If set to
	// `true`, then no security groups are applied to the port. If set to `false` and
	// no `securityGroupIds` are specified, then the port will yield to the default
	// behavior of the Networking service, which is to usually apply the "default"
	// security group.
	NoSecurityGroups pulumi.BoolPtrInput
	// Whether to explicitly enable or disable
	// port security on the port. Port Security is usually enabled by default, so
	// omitting argument will usually result in a value of `true`. Setting this
	// explicitly to `false` will disable port security. In order to disable port
	// security, the port must not have any security groups. Valid values are `true`
	// and `false`.
	PortSecurityEnabled pulumi.BoolPtrInput
	// Reference to the associated QoS policy.
	QosPolicyId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a port. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// port.
	Region pulumi.StringPtrInput
	// A list
	// of security group IDs to apply to the port. The security groups must be
	// specified by ID and not name (as opposed to how they are configured with
	// the Compute Instance).
	SecurityGroupIds pulumi.StringArrayInput
	// A set of string tags for the port.
	Tags pulumi.StringArrayInput
	// The owner of the port. Required if admin wants
	// to create a port for another tenant. Changing this creates a new port.
	TenantId pulumi.StringPtrInput
	// Map of additional options.
	ValueSpecs pulumi.MapInput
}

func (PortState) ElementType() reflect.Type {
	return reflect.TypeOf((*portState)(nil)).Elem()
}

type portArgs struct {
	// Administrative up/down status for the port
	// (must be `true` or `false` if provided). Changing this updates the
	// `adminStateUp` of an existing port.
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// An IP/MAC Address pair of additional IP
	// addresses that can be active on this port. The structure is described
	// below.
	AllowedAddressPairs []PortAllowedAddressPair `pulumi:"allowedAddressPairs"`
	// The port binding allows to specify binding information
	// for the port. The structure is described below.
	Binding *PortBinding `pulumi:"binding"`
	// Human-readable description of the port. Changing
	// this updates the `description` of an existing port.
	Description *string `pulumi:"description"`
	// The ID of the device attached to the port. Changing this
	// creates a new port.
	DeviceId *string `pulumi:"deviceId"`
	// The device owner of the port. Changing this creates
	// a new port.
	DeviceOwner *string `pulumi:"deviceOwner"`
	// The port DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName *string `pulumi:"dnsName"`
	// An extra DHCP option that needs to be configured
	// on the port. The structure is described below. Can be specified multiple
	// times.
	ExtraDhcpOptions []PortExtraDhcpOption `pulumi:"extraDhcpOptions"`
	// An array of desired IPs for
	// this port. The structure is described below.
	FixedIps []PortFixedIp `pulumi:"fixedIps"`
	// The additional MAC address.
	MacAddress *string `pulumi:"macAddress"`
	// Name of the DHCP option.
	Name *string `pulumi:"name"`
	// The ID of the network to attach the port to. Changing
	// this creates a new port.
	NetworkId string `pulumi:"networkId"`
	// Create a port with no fixed
	// IP address. This will also remove any fixed IPs previously set on a port. `true`
	// is the only valid value for this argument.
	NoFixedIp *bool `pulumi:"noFixedIp"`
	// If set to
	// `true`, then no security groups are applied to the port. If set to `false` and
	// no `securityGroupIds` are specified, then the port will yield to the default
	// behavior of the Networking service, which is to usually apply the "default"
	// security group.
	NoSecurityGroups *bool `pulumi:"noSecurityGroups"`
	// Whether to explicitly enable or disable
	// port security on the port. Port Security is usually enabled by default, so
	// omitting argument will usually result in a value of `true`. Setting this
	// explicitly to `false` will disable port security. In order to disable port
	// security, the port must not have any security groups. Valid values are `true`
	// and `false`.
	PortSecurityEnabled *bool `pulumi:"portSecurityEnabled"`
	// Reference to the associated QoS policy.
	QosPolicyId *string `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a port. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// port.
	Region *string `pulumi:"region"`
	// A list
	// of security group IDs to apply to the port. The security groups must be
	// specified by ID and not name (as opposed to how they are configured with
	// the Compute Instance).
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// A set of string tags for the port.
	Tags []string `pulumi:"tags"`
	// The owner of the port. Required if admin wants
	// to create a port for another tenant. Changing this creates a new port.
	TenantId *string `pulumi:"tenantId"`
	// Map of additional options.
	ValueSpecs map[string]interface{} `pulumi:"valueSpecs"`
}

// The set of arguments for constructing a Port resource.
type PortArgs struct {
	// Administrative up/down status for the port
	// (must be `true` or `false` if provided). Changing this updates the
	// `adminStateUp` of an existing port.
	AdminStateUp pulumi.BoolPtrInput
	// An IP/MAC Address pair of additional IP
	// addresses that can be active on this port. The structure is described
	// below.
	AllowedAddressPairs PortAllowedAddressPairArrayInput
	// The port binding allows to specify binding information
	// for the port. The structure is described below.
	Binding PortBindingPtrInput
	// Human-readable description of the port. Changing
	// this updates the `description` of an existing port.
	Description pulumi.StringPtrInput
	// The ID of the device attached to the port. Changing this
	// creates a new port.
	DeviceId pulumi.StringPtrInput
	// The device owner of the port. Changing this creates
	// a new port.
	DeviceOwner pulumi.StringPtrInput
	// The port DNS name. Available, when Neutron DNS extension
	// is enabled.
	DnsName pulumi.StringPtrInput
	// An extra DHCP option that needs to be configured
	// on the port. The structure is described below. Can be specified multiple
	// times.
	ExtraDhcpOptions PortExtraDhcpOptionArrayInput
	// An array of desired IPs for
	// this port. The structure is described below.
	FixedIps PortFixedIpArrayInput
	// The additional MAC address.
	MacAddress pulumi.StringPtrInput
	// Name of the DHCP option.
	Name pulumi.StringPtrInput
	// The ID of the network to attach the port to. Changing
	// this creates a new port.
	NetworkId pulumi.StringInput
	// Create a port with no fixed
	// IP address. This will also remove any fixed IPs previously set on a port. `true`
	// is the only valid value for this argument.
	NoFixedIp pulumi.BoolPtrInput
	// If set to
	// `true`, then no security groups are applied to the port. If set to `false` and
	// no `securityGroupIds` are specified, then the port will yield to the default
	// behavior of the Networking service, which is to usually apply the "default"
	// security group.
	NoSecurityGroups pulumi.BoolPtrInput
	// Whether to explicitly enable or disable
	// port security on the port. Port Security is usually enabled by default, so
	// omitting argument will usually result in a value of `true`. Setting this
	// explicitly to `false` will disable port security. In order to disable port
	// security, the port must not have any security groups. Valid values are `true`
	// and `false`.
	PortSecurityEnabled pulumi.BoolPtrInput
	// Reference to the associated QoS policy.
	QosPolicyId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a port. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// port.
	Region pulumi.StringPtrInput
	// A list
	// of security group IDs to apply to the port. The security groups must be
	// specified by ID and not name (as opposed to how they are configured with
	// the Compute Instance).
	SecurityGroupIds pulumi.StringArrayInput
	// A set of string tags for the port.
	Tags pulumi.StringArrayInput
	// The owner of the port. Required if admin wants
	// to create a port for another tenant. Changing this creates a new port.
	TenantId pulumi.StringPtrInput
	// Map of additional options.
	ValueSpecs pulumi.MapInput
}

func (PortArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*portArgs)(nil)).Elem()
}

type PortInput interface {
	pulumi.Input

	ToPortOutput() PortOutput
	ToPortOutputWithContext(ctx context.Context) PortOutput
}

func (*Port) ElementType() reflect.Type {
	return reflect.TypeOf((*Port)(nil))
}

func (i *Port) ToPortOutput() PortOutput {
	return i.ToPortOutputWithContext(context.Background())
}

func (i *Port) ToPortOutputWithContext(ctx context.Context) PortOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PortOutput)
}

func (i *Port) ToPortPtrOutput() PortPtrOutput {
	return i.ToPortPtrOutputWithContext(context.Background())
}

func (i *Port) ToPortPtrOutputWithContext(ctx context.Context) PortPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PortPtrOutput)
}

type PortPtrInput interface {
	pulumi.Input

	ToPortPtrOutput() PortPtrOutput
	ToPortPtrOutputWithContext(ctx context.Context) PortPtrOutput
}

type portPtrType PortArgs

func (*portPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Port)(nil))
}

func (i *portPtrType) ToPortPtrOutput() PortPtrOutput {
	return i.ToPortPtrOutputWithContext(context.Background())
}

func (i *portPtrType) ToPortPtrOutputWithContext(ctx context.Context) PortPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PortPtrOutput)
}

// PortArrayInput is an input type that accepts PortArray and PortArrayOutput values.
// You can construct a concrete instance of `PortArrayInput` via:
//
//          PortArray{ PortArgs{...} }
type PortArrayInput interface {
	pulumi.Input

	ToPortArrayOutput() PortArrayOutput
	ToPortArrayOutputWithContext(context.Context) PortArrayOutput
}

type PortArray []PortInput

func (PortArray) ElementType() reflect.Type {
	return reflect.TypeOf(([]*Port)(nil))
}

func (i PortArray) ToPortArrayOutput() PortArrayOutput {
	return i.ToPortArrayOutputWithContext(context.Background())
}

func (i PortArray) ToPortArrayOutputWithContext(ctx context.Context) PortArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PortArrayOutput)
}

// PortMapInput is an input type that accepts PortMap and PortMapOutput values.
// You can construct a concrete instance of `PortMapInput` via:
//
//          PortMap{ "key": PortArgs{...} }
type PortMapInput interface {
	pulumi.Input

	ToPortMapOutput() PortMapOutput
	ToPortMapOutputWithContext(context.Context) PortMapOutput
}

type PortMap map[string]PortInput

func (PortMap) ElementType() reflect.Type {
	return reflect.TypeOf((map[string]*Port)(nil))
}

func (i PortMap) ToPortMapOutput() PortMapOutput {
	return i.ToPortMapOutputWithContext(context.Background())
}

func (i PortMap) ToPortMapOutputWithContext(ctx context.Context) PortMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PortMapOutput)
}

type PortOutput struct {
	*pulumi.OutputState
}

func (PortOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Port)(nil))
}

func (o PortOutput) ToPortOutput() PortOutput {
	return o
}

func (o PortOutput) ToPortOutputWithContext(ctx context.Context) PortOutput {
	return o
}

func (o PortOutput) ToPortPtrOutput() PortPtrOutput {
	return o.ToPortPtrOutputWithContext(context.Background())
}

func (o PortOutput) ToPortPtrOutputWithContext(ctx context.Context) PortPtrOutput {
	return o.ApplyT(func(v Port) *Port {
		return &v
	}).(PortPtrOutput)
}

type PortPtrOutput struct {
	*pulumi.OutputState
}

func (PortPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Port)(nil))
}

func (o PortPtrOutput) ToPortPtrOutput() PortPtrOutput {
	return o
}

func (o PortPtrOutput) ToPortPtrOutputWithContext(ctx context.Context) PortPtrOutput {
	return o
}

type PortArrayOutput struct{ *pulumi.OutputState }

func (PortArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Port)(nil))
}

func (o PortArrayOutput) ToPortArrayOutput() PortArrayOutput {
	return o
}

func (o PortArrayOutput) ToPortArrayOutputWithContext(ctx context.Context) PortArrayOutput {
	return o
}

func (o PortArrayOutput) Index(i pulumi.IntInput) PortOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Port {
		return vs[0].([]Port)[vs[1].(int)]
	}).(PortOutput)
}

type PortMapOutput struct{ *pulumi.OutputState }

func (PortMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]Port)(nil))
}

func (o PortMapOutput) ToPortMapOutput() PortMapOutput {
	return o
}

func (o PortMapOutput) ToPortMapOutputWithContext(ctx context.Context) PortMapOutput {
	return o
}

func (o PortMapOutput) MapIndex(k pulumi.StringInput) PortOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) Port {
		return vs[0].(map[string]Port)[vs[1].(string)]
	}).(PortOutput)
}

func init() {
	pulumi.RegisterOutputType(PortOutput{})
	pulumi.RegisterOutputType(PortPtrOutput{})
	pulumi.RegisterOutputType(PortArrayOutput{})
	pulumi.RegisterOutputType(PortMapOutput{})
}