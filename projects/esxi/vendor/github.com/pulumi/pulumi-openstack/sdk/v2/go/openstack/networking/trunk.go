// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Manages a networking V2 trunk resource within OpenStack.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack/compute"
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
// 		_, err = networking.NewSubnet(ctx, "subnet1", &networking.SubnetArgs{
// 			Cidr:       pulumi.String("192.168.1.0/24"),
// 			EnableDhcp: pulumi.Bool(true),
// 			IpVersion:  pulumi.Int(4),
// 			NetworkId:  network1.ID(),
// 			NoGateway:  pulumi.Bool(true),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		parentPort1, err := networking.NewPort(ctx, "parentPort1", &networking.PortArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 			NetworkId:    network1.ID(),
// 		}, pulumi.DependsOn([]pulumi.Resource{
// 			"openstack_networking_subnet_v2.subnet_1",
// 		}))
// 		if err != nil {
// 			return err
// 		}
// 		subport1, err := networking.NewPort(ctx, "subport1", &networking.PortArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 			NetworkId:    network1.ID(),
// 		}, pulumi.DependsOn([]pulumi.Resource{
// 			"openstack_networking_subnet_v2.subnet_1",
// 		}))
// 		if err != nil {
// 			return err
// 		}
// 		trunk1, err := networking.NewTrunk(ctx, "trunk1", &networking.TrunkArgs{
// 			AdminStateUp: pulumi.Bool(true),
// 			PortId:       parentPort1.ID(),
// 			SubPorts: networking.TrunkSubPortArray{
// 				&networking.TrunkSubPortArgs{
// 					PortId:           subport1.ID(),
// 					SegmentationId:   pulumi.Int(1),
// 					SegmentationType: pulumi.String("vlan"),
// 				},
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = compute.NewInstance(ctx, "instance1", &compute.InstanceArgs{
// 			Networks: compute.InstanceNetworkArray{
// 				&compute.InstanceNetworkArgs{
// 					Port: trunk1.PortId,
// 				},
// 			},
// 			SecurityGroups: pulumi.StringArray{
// 				pulumi.String("default"),
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type Trunk struct {
	pulumi.CustomResourceState

	// Administrative up/down status for the trunk
	// (must be "true" or "false" if provided). Changing this updates the
	// `adminStateUp` of an existing trunk.
	AdminStateUp pulumi.BoolPtrOutput `pulumi:"adminStateUp"`
	// The collection of tags assigned on the trunk, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayOutput `pulumi:"allTags"`
	// Human-readable description of the trunk. Changing this
	// updates the name of the existing trunk.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// A unique name for the trunk. Changing this
	// updates the `name` of an existing trunk.
	Name pulumi.StringOutput `pulumi:"name"`
	// The ID of the port to be made a subport of the trunk.
	PortId pulumi.StringOutput `pulumi:"portId"`
	// The region in which to obtain the V2 networking client.
	// A networking client is needed to create a trunk. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// trunk.
	Region pulumi.StringOutput `pulumi:"region"`
	// The set of ports that will be made subports of the trunk.
	// The structure of each subport is described below.
	SubPorts TrunkSubPortArrayOutput `pulumi:"subPorts"`
	// A set of string tags for the port.
	Tags pulumi.StringArrayOutput `pulumi:"tags"`
	// The owner of the Trunk. Required if admin wants
	// to create a trunk on behalf of another tenant. Changing this creates a new trunk.
	TenantId pulumi.StringOutput `pulumi:"tenantId"`
}

// NewTrunk registers a new resource with the given unique name, arguments, and options.
func NewTrunk(ctx *pulumi.Context,
	name string, args *TrunkArgs, opts ...pulumi.ResourceOption) (*Trunk, error) {
	if args == nil || args.PortId == nil {
		return nil, errors.New("missing required argument 'PortId'")
	}
	if args == nil {
		args = &TrunkArgs{}
	}
	var resource Trunk
	err := ctx.RegisterResource("openstack:networking/trunk:Trunk", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetTrunk gets an existing Trunk resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetTrunk(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *TrunkState, opts ...pulumi.ResourceOption) (*Trunk, error) {
	var resource Trunk
	err := ctx.ReadResource("openstack:networking/trunk:Trunk", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Trunk resources.
type trunkState struct {
	// Administrative up/down status for the trunk
	// (must be "true" or "false" if provided). Changing this updates the
	// `adminStateUp` of an existing trunk.
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// The collection of tags assigned on the trunk, which have been
	// explicitly and implicitly added.
	AllTags []string `pulumi:"allTags"`
	// Human-readable description of the trunk. Changing this
	// updates the name of the existing trunk.
	Description *string `pulumi:"description"`
	// A unique name for the trunk. Changing this
	// updates the `name` of an existing trunk.
	Name *string `pulumi:"name"`
	// The ID of the port to be made a subport of the trunk.
	PortId *string `pulumi:"portId"`
	// The region in which to obtain the V2 networking client.
	// A networking client is needed to create a trunk. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// trunk.
	Region *string `pulumi:"region"`
	// The set of ports that will be made subports of the trunk.
	// The structure of each subport is described below.
	SubPorts []TrunkSubPort `pulumi:"subPorts"`
	// A set of string tags for the port.
	Tags []string `pulumi:"tags"`
	// The owner of the Trunk. Required if admin wants
	// to create a trunk on behalf of another tenant. Changing this creates a new trunk.
	TenantId *string `pulumi:"tenantId"`
}

type TrunkState struct {
	// Administrative up/down status for the trunk
	// (must be "true" or "false" if provided). Changing this updates the
	// `adminStateUp` of an existing trunk.
	AdminStateUp pulumi.BoolPtrInput
	// The collection of tags assigned on the trunk, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayInput
	// Human-readable description of the trunk. Changing this
	// updates the name of the existing trunk.
	Description pulumi.StringPtrInput
	// A unique name for the trunk. Changing this
	// updates the `name` of an existing trunk.
	Name pulumi.StringPtrInput
	// The ID of the port to be made a subport of the trunk.
	PortId pulumi.StringPtrInput
	// The region in which to obtain the V2 networking client.
	// A networking client is needed to create a trunk. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// trunk.
	Region pulumi.StringPtrInput
	// The set of ports that will be made subports of the trunk.
	// The structure of each subport is described below.
	SubPorts TrunkSubPortArrayInput
	// A set of string tags for the port.
	Tags pulumi.StringArrayInput
	// The owner of the Trunk. Required if admin wants
	// to create a trunk on behalf of another tenant. Changing this creates a new trunk.
	TenantId pulumi.StringPtrInput
}

func (TrunkState) ElementType() reflect.Type {
	return reflect.TypeOf((*trunkState)(nil)).Elem()
}

type trunkArgs struct {
	// Administrative up/down status for the trunk
	// (must be "true" or "false" if provided). Changing this updates the
	// `adminStateUp` of an existing trunk.
	AdminStateUp *bool `pulumi:"adminStateUp"`
	// Human-readable description of the trunk. Changing this
	// updates the name of the existing trunk.
	Description *string `pulumi:"description"`
	// A unique name for the trunk. Changing this
	// updates the `name` of an existing trunk.
	Name *string `pulumi:"name"`
	// The ID of the port to be made a subport of the trunk.
	PortId string `pulumi:"portId"`
	// The region in which to obtain the V2 networking client.
	// A networking client is needed to create a trunk. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// trunk.
	Region *string `pulumi:"region"`
	// The set of ports that will be made subports of the trunk.
	// The structure of each subport is described below.
	SubPorts []TrunkSubPort `pulumi:"subPorts"`
	// A set of string tags for the port.
	Tags []string `pulumi:"tags"`
	// The owner of the Trunk. Required if admin wants
	// to create a trunk on behalf of another tenant. Changing this creates a new trunk.
	TenantId *string `pulumi:"tenantId"`
}

// The set of arguments for constructing a Trunk resource.
type TrunkArgs struct {
	// Administrative up/down status for the trunk
	// (must be "true" or "false" if provided). Changing this updates the
	// `adminStateUp` of an existing trunk.
	AdminStateUp pulumi.BoolPtrInput
	// Human-readable description of the trunk. Changing this
	// updates the name of the existing trunk.
	Description pulumi.StringPtrInput
	// A unique name for the trunk. Changing this
	// updates the `name` of an existing trunk.
	Name pulumi.StringPtrInput
	// The ID of the port to be made a subport of the trunk.
	PortId pulumi.StringInput
	// The region in which to obtain the V2 networking client.
	// A networking client is needed to create a trunk. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// trunk.
	Region pulumi.StringPtrInput
	// The set of ports that will be made subports of the trunk.
	// The structure of each subport is described below.
	SubPorts TrunkSubPortArrayInput
	// A set of string tags for the port.
	Tags pulumi.StringArrayInput
	// The owner of the Trunk. Required if admin wants
	// to create a trunk on behalf of another tenant. Changing this creates a new trunk.
	TenantId pulumi.StringPtrInput
}

func (TrunkArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*trunkArgs)(nil)).Elem()
}
