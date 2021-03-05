// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Manages a V2 Neutron QoS policy resource within OpenStack.
//
// ## Example Usage
// ### Create a QoS Policy
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
// 		_, err := networking.NewQosPolicy(ctx, "qosPolicy1", &networking.QosPolicyArgs{
// 			Description: pulumi.String("bw_limit"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type QosPolicy struct {
	pulumi.CustomResourceState

	// The collection of tags assigned on the QoS policy, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayOutput `pulumi:"allTags"`
	// The time at which QoS policy was created.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// The human-readable description for the QoS policy.
	// Changing this updates the description of the existing QoS policy.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Indicates whether the QoS policy is default
	// QoS policy or not. Changing this updates the default status of the existing
	// QoS policy.
	IsDefault pulumi.BoolPtrOutput `pulumi:"isDefault"`
	// The name of the QoS policy. Changing this updates the name of
	// the existing QoS policy.
	Name pulumi.StringOutput `pulumi:"name"`
	// The owner of the QoS policy. Required if admin wants to
	// create a QoS policy for another project. Changing this creates a new QoS policy.
	ProjectId pulumi.StringOutput `pulumi:"projectId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron Qos policy. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// QoS policy.
	Region pulumi.StringOutput `pulumi:"region"`
	// The revision number of the QoS policy.
	RevisionNumber pulumi.IntOutput `pulumi:"revisionNumber"`
	// Indicates whether this QoS policy is shared across
	// all projects. Changing this updates the shared status of the existing
	// QoS policy.
	Shared pulumi.BoolPtrOutput `pulumi:"shared"`
	// A set of string tags for the QoS policy.
	Tags pulumi.StringArrayOutput `pulumi:"tags"`
	// The time at which QoS policy was created.
	UpdatedAt pulumi.StringOutput `pulumi:"updatedAt"`
	// Map of additional options.
	ValueSpecs pulumi.MapOutput `pulumi:"valueSpecs"`
}

// NewQosPolicy registers a new resource with the given unique name, arguments, and options.
func NewQosPolicy(ctx *pulumi.Context,
	name string, args *QosPolicyArgs, opts ...pulumi.ResourceOption) (*QosPolicy, error) {
	if args == nil {
		args = &QosPolicyArgs{}
	}
	var resource QosPolicy
	err := ctx.RegisterResource("openstack:networking/qosPolicy:QosPolicy", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetQosPolicy gets an existing QosPolicy resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetQosPolicy(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *QosPolicyState, opts ...pulumi.ResourceOption) (*QosPolicy, error) {
	var resource QosPolicy
	err := ctx.ReadResource("openstack:networking/qosPolicy:QosPolicy", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering QosPolicy resources.
type qosPolicyState struct {
	// The collection of tags assigned on the QoS policy, which have been
	// explicitly and implicitly added.
	AllTags []string `pulumi:"allTags"`
	// The time at which QoS policy was created.
	CreatedAt *string `pulumi:"createdAt"`
	// The human-readable description for the QoS policy.
	// Changing this updates the description of the existing QoS policy.
	Description *string `pulumi:"description"`
	// Indicates whether the QoS policy is default
	// QoS policy or not. Changing this updates the default status of the existing
	// QoS policy.
	IsDefault *bool `pulumi:"isDefault"`
	// The name of the QoS policy. Changing this updates the name of
	// the existing QoS policy.
	Name *string `pulumi:"name"`
	// The owner of the QoS policy. Required if admin wants to
	// create a QoS policy for another project. Changing this creates a new QoS policy.
	ProjectId *string `pulumi:"projectId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron Qos policy. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// QoS policy.
	Region *string `pulumi:"region"`
	// The revision number of the QoS policy.
	RevisionNumber *int `pulumi:"revisionNumber"`
	// Indicates whether this QoS policy is shared across
	// all projects. Changing this updates the shared status of the existing
	// QoS policy.
	Shared *bool `pulumi:"shared"`
	// A set of string tags for the QoS policy.
	Tags []string `pulumi:"tags"`
	// The time at which QoS policy was created.
	UpdatedAt *string `pulumi:"updatedAt"`
	// Map of additional options.
	ValueSpecs map[string]interface{} `pulumi:"valueSpecs"`
}

type QosPolicyState struct {
	// The collection of tags assigned on the QoS policy, which have been
	// explicitly and implicitly added.
	AllTags pulumi.StringArrayInput
	// The time at which QoS policy was created.
	CreatedAt pulumi.StringPtrInput
	// The human-readable description for the QoS policy.
	// Changing this updates the description of the existing QoS policy.
	Description pulumi.StringPtrInput
	// Indicates whether the QoS policy is default
	// QoS policy or not. Changing this updates the default status of the existing
	// QoS policy.
	IsDefault pulumi.BoolPtrInput
	// The name of the QoS policy. Changing this updates the name of
	// the existing QoS policy.
	Name pulumi.StringPtrInput
	// The owner of the QoS policy. Required if admin wants to
	// create a QoS policy for another project. Changing this creates a new QoS policy.
	ProjectId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron Qos policy. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// QoS policy.
	Region pulumi.StringPtrInput
	// The revision number of the QoS policy.
	RevisionNumber pulumi.IntPtrInput
	// Indicates whether this QoS policy is shared across
	// all projects. Changing this updates the shared status of the existing
	// QoS policy.
	Shared pulumi.BoolPtrInput
	// A set of string tags for the QoS policy.
	Tags pulumi.StringArrayInput
	// The time at which QoS policy was created.
	UpdatedAt pulumi.StringPtrInput
	// Map of additional options.
	ValueSpecs pulumi.MapInput
}

func (QosPolicyState) ElementType() reflect.Type {
	return reflect.TypeOf((*qosPolicyState)(nil)).Elem()
}

type qosPolicyArgs struct {
	// The human-readable description for the QoS policy.
	// Changing this updates the description of the existing QoS policy.
	Description *string `pulumi:"description"`
	// Indicates whether the QoS policy is default
	// QoS policy or not. Changing this updates the default status of the existing
	// QoS policy.
	IsDefault *bool `pulumi:"isDefault"`
	// The name of the QoS policy. Changing this updates the name of
	// the existing QoS policy.
	Name *string `pulumi:"name"`
	// The owner of the QoS policy. Required if admin wants to
	// create a QoS policy for another project. Changing this creates a new QoS policy.
	ProjectId *string `pulumi:"projectId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron Qos policy. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// QoS policy.
	Region *string `pulumi:"region"`
	// Indicates whether this QoS policy is shared across
	// all projects. Changing this updates the shared status of the existing
	// QoS policy.
	Shared *bool `pulumi:"shared"`
	// A set of string tags for the QoS policy.
	Tags []string `pulumi:"tags"`
	// Map of additional options.
	ValueSpecs map[string]interface{} `pulumi:"valueSpecs"`
}

// The set of arguments for constructing a QosPolicy resource.
type QosPolicyArgs struct {
	// The human-readable description for the QoS policy.
	// Changing this updates the description of the existing QoS policy.
	Description pulumi.StringPtrInput
	// Indicates whether the QoS policy is default
	// QoS policy or not. Changing this updates the default status of the existing
	// QoS policy.
	IsDefault pulumi.BoolPtrInput
	// The name of the QoS policy. Changing this updates the name of
	// the existing QoS policy.
	Name pulumi.StringPtrInput
	// The owner of the QoS policy. Required if admin wants to
	// create a QoS policy for another project. Changing this creates a new QoS policy.
	ProjectId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron Qos policy. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new
	// QoS policy.
	Region pulumi.StringPtrInput
	// Indicates whether this QoS policy is shared across
	// all projects. Changing this updates the shared status of the existing
	// QoS policy.
	Shared pulumi.BoolPtrInput
	// A set of string tags for the QoS policy.
	Tags pulumi.StringArrayInput
	// Map of additional options.
	ValueSpecs pulumi.MapInput
}

func (QosPolicyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*qosPolicyArgs)(nil)).Elem()
}