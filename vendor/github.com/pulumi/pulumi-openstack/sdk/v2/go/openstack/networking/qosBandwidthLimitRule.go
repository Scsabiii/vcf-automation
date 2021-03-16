// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Manages a V2 Neutron QoS bandwidth limit rule resource within OpenStack.
//
// ## Example Usage
// ### Create a QoS Policy with some bandwidth limit rule
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
// 		qosPolicy1, err := networking.NewQosPolicy(ctx, "qosPolicy1", &networking.QosPolicyArgs{
// 			Description: pulumi.String("bw_limit"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = networking.NewQosBandwidthLimitRule(ctx, "bwLimitRule1", &networking.QosBandwidthLimitRuleArgs{
// 			Direction:    pulumi.String("egress"),
// 			MaxBurstKbps: pulumi.Int(300),
// 			MaxKbps:      pulumi.Int(3000),
// 			QosPolicyId:  qosPolicy1.ID(),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
//
// ## Import
//
// QoS bandwidth limit rules can be imported using the `qos_policy_id/bandwidth_limit_rule` format, e.g.
//
// ```sh
//  $ pulumi import openstack:networking/qosBandwidthLimitRule:QosBandwidthLimitRule bw_limit_rule_1 d6ae28ce-fcb5-4180-aa62-d260a27e09ae/46dfb556-b92f-48ce-94c5-9a9e2140de94
// ```
type QosBandwidthLimitRule struct {
	pulumi.CustomResourceState

	// The direction of traffic. Defaults to "egress". Changing this updates the direction of the
	// existing QoS bandwidth limit rule.
	Direction pulumi.StringPtrOutput `pulumi:"direction"`
	// The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
	// maximum burst size in kilobits of the existing QoS bandwidth limit rule.
	MaxBurstKbps pulumi.IntPtrOutput `pulumi:"maxBurstKbps"`
	// The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
	// maximum kilobits per second of the existing QoS bandwidth limit rule.
	MaxKbps pulumi.IntOutput `pulumi:"maxKbps"`
	// The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
	QosPolicyId pulumi.StringOutput `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
	Region pulumi.StringOutput `pulumi:"region"`
}

// NewQosBandwidthLimitRule registers a new resource with the given unique name, arguments, and options.
func NewQosBandwidthLimitRule(ctx *pulumi.Context,
	name string, args *QosBandwidthLimitRuleArgs, opts ...pulumi.ResourceOption) (*QosBandwidthLimitRule, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.MaxKbps == nil {
		return nil, errors.New("invalid value for required argument 'MaxKbps'")
	}
	if args.QosPolicyId == nil {
		return nil, errors.New("invalid value for required argument 'QosPolicyId'")
	}
	var resource QosBandwidthLimitRule
	err := ctx.RegisterResource("openstack:networking/qosBandwidthLimitRule:QosBandwidthLimitRule", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetQosBandwidthLimitRule gets an existing QosBandwidthLimitRule resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetQosBandwidthLimitRule(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *QosBandwidthLimitRuleState, opts ...pulumi.ResourceOption) (*QosBandwidthLimitRule, error) {
	var resource QosBandwidthLimitRule
	err := ctx.ReadResource("openstack:networking/qosBandwidthLimitRule:QosBandwidthLimitRule", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering QosBandwidthLimitRule resources.
type qosBandwidthLimitRuleState struct {
	// The direction of traffic. Defaults to "egress". Changing this updates the direction of the
	// existing QoS bandwidth limit rule.
	Direction *string `pulumi:"direction"`
	// The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
	// maximum burst size in kilobits of the existing QoS bandwidth limit rule.
	MaxBurstKbps *int `pulumi:"maxBurstKbps"`
	// The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
	// maximum kilobits per second of the existing QoS bandwidth limit rule.
	MaxKbps *int `pulumi:"maxKbps"`
	// The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
	QosPolicyId *string `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
	Region *string `pulumi:"region"`
}

type QosBandwidthLimitRuleState struct {
	// The direction of traffic. Defaults to "egress". Changing this updates the direction of the
	// existing QoS bandwidth limit rule.
	Direction pulumi.StringPtrInput
	// The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
	// maximum burst size in kilobits of the existing QoS bandwidth limit rule.
	MaxBurstKbps pulumi.IntPtrInput
	// The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
	// maximum kilobits per second of the existing QoS bandwidth limit rule.
	MaxKbps pulumi.IntPtrInput
	// The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
	QosPolicyId pulumi.StringPtrInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
	Region pulumi.StringPtrInput
}

func (QosBandwidthLimitRuleState) ElementType() reflect.Type {
	return reflect.TypeOf((*qosBandwidthLimitRuleState)(nil)).Elem()
}

type qosBandwidthLimitRuleArgs struct {
	// The direction of traffic. Defaults to "egress". Changing this updates the direction of the
	// existing QoS bandwidth limit rule.
	Direction *string `pulumi:"direction"`
	// The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
	// maximum burst size in kilobits of the existing QoS bandwidth limit rule.
	MaxBurstKbps *int `pulumi:"maxBurstKbps"`
	// The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
	// maximum kilobits per second of the existing QoS bandwidth limit rule.
	MaxKbps int `pulumi:"maxKbps"`
	// The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
	QosPolicyId string `pulumi:"qosPolicyId"`
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
	Region *string `pulumi:"region"`
}

// The set of arguments for constructing a QosBandwidthLimitRule resource.
type QosBandwidthLimitRuleArgs struct {
	// The direction of traffic. Defaults to "egress". Changing this updates the direction of the
	// existing QoS bandwidth limit rule.
	Direction pulumi.StringPtrInput
	// The maximum burst size in kilobits of a QoS bandwidth limit rule. Changing this updates the
	// maximum burst size in kilobits of the existing QoS bandwidth limit rule.
	MaxBurstKbps pulumi.IntPtrInput
	// The maximum kilobits per second of a QoS bandwidth limit rule. Changing this updates the
	// maximum kilobits per second of the existing QoS bandwidth limit rule.
	MaxKbps pulumi.IntInput
	// The QoS policy reference. Changing this creates a new QoS bandwidth limit rule.
	QosPolicyId pulumi.StringInput
	// The region in which to obtain the V2 Networking client.
	// A Networking client is needed to create a Neutron QoS bandwidth limit rule. If omitted, the
	// `region` argument of the provider is used. Changing this creates a new QoS bandwidth limit rule.
	Region pulumi.StringPtrInput
}

func (QosBandwidthLimitRuleArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*qosBandwidthLimitRuleArgs)(nil)).Elem()
}

type QosBandwidthLimitRuleInput interface {
	pulumi.Input

	ToQosBandwidthLimitRuleOutput() QosBandwidthLimitRuleOutput
	ToQosBandwidthLimitRuleOutputWithContext(ctx context.Context) QosBandwidthLimitRuleOutput
}

func (*QosBandwidthLimitRule) ElementType() reflect.Type {
	return reflect.TypeOf((*QosBandwidthLimitRule)(nil))
}

func (i *QosBandwidthLimitRule) ToQosBandwidthLimitRuleOutput() QosBandwidthLimitRuleOutput {
	return i.ToQosBandwidthLimitRuleOutputWithContext(context.Background())
}

func (i *QosBandwidthLimitRule) ToQosBandwidthLimitRuleOutputWithContext(ctx context.Context) QosBandwidthLimitRuleOutput {
	return pulumi.ToOutputWithContext(ctx, i).(QosBandwidthLimitRuleOutput)
}

func (i *QosBandwidthLimitRule) ToQosBandwidthLimitRulePtrOutput() QosBandwidthLimitRulePtrOutput {
	return i.ToQosBandwidthLimitRulePtrOutputWithContext(context.Background())
}

func (i *QosBandwidthLimitRule) ToQosBandwidthLimitRulePtrOutputWithContext(ctx context.Context) QosBandwidthLimitRulePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(QosBandwidthLimitRulePtrOutput)
}

type QosBandwidthLimitRulePtrInput interface {
	pulumi.Input

	ToQosBandwidthLimitRulePtrOutput() QosBandwidthLimitRulePtrOutput
	ToQosBandwidthLimitRulePtrOutputWithContext(ctx context.Context) QosBandwidthLimitRulePtrOutput
}

type qosBandwidthLimitRulePtrType QosBandwidthLimitRuleArgs

func (*qosBandwidthLimitRulePtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**QosBandwidthLimitRule)(nil))
}

func (i *qosBandwidthLimitRulePtrType) ToQosBandwidthLimitRulePtrOutput() QosBandwidthLimitRulePtrOutput {
	return i.ToQosBandwidthLimitRulePtrOutputWithContext(context.Background())
}

func (i *qosBandwidthLimitRulePtrType) ToQosBandwidthLimitRulePtrOutputWithContext(ctx context.Context) QosBandwidthLimitRulePtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(QosBandwidthLimitRulePtrOutput)
}

// QosBandwidthLimitRuleArrayInput is an input type that accepts QosBandwidthLimitRuleArray and QosBandwidthLimitRuleArrayOutput values.
// You can construct a concrete instance of `QosBandwidthLimitRuleArrayInput` via:
//
//          QosBandwidthLimitRuleArray{ QosBandwidthLimitRuleArgs{...} }
type QosBandwidthLimitRuleArrayInput interface {
	pulumi.Input

	ToQosBandwidthLimitRuleArrayOutput() QosBandwidthLimitRuleArrayOutput
	ToQosBandwidthLimitRuleArrayOutputWithContext(context.Context) QosBandwidthLimitRuleArrayOutput
}

type QosBandwidthLimitRuleArray []QosBandwidthLimitRuleInput

func (QosBandwidthLimitRuleArray) ElementType() reflect.Type {
	return reflect.TypeOf(([]*QosBandwidthLimitRule)(nil))
}

func (i QosBandwidthLimitRuleArray) ToQosBandwidthLimitRuleArrayOutput() QosBandwidthLimitRuleArrayOutput {
	return i.ToQosBandwidthLimitRuleArrayOutputWithContext(context.Background())
}

func (i QosBandwidthLimitRuleArray) ToQosBandwidthLimitRuleArrayOutputWithContext(ctx context.Context) QosBandwidthLimitRuleArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(QosBandwidthLimitRuleArrayOutput)
}

// QosBandwidthLimitRuleMapInput is an input type that accepts QosBandwidthLimitRuleMap and QosBandwidthLimitRuleMapOutput values.
// You can construct a concrete instance of `QosBandwidthLimitRuleMapInput` via:
//
//          QosBandwidthLimitRuleMap{ "key": QosBandwidthLimitRuleArgs{...} }
type QosBandwidthLimitRuleMapInput interface {
	pulumi.Input

	ToQosBandwidthLimitRuleMapOutput() QosBandwidthLimitRuleMapOutput
	ToQosBandwidthLimitRuleMapOutputWithContext(context.Context) QosBandwidthLimitRuleMapOutput
}

type QosBandwidthLimitRuleMap map[string]QosBandwidthLimitRuleInput

func (QosBandwidthLimitRuleMap) ElementType() reflect.Type {
	return reflect.TypeOf((map[string]*QosBandwidthLimitRule)(nil))
}

func (i QosBandwidthLimitRuleMap) ToQosBandwidthLimitRuleMapOutput() QosBandwidthLimitRuleMapOutput {
	return i.ToQosBandwidthLimitRuleMapOutputWithContext(context.Background())
}

func (i QosBandwidthLimitRuleMap) ToQosBandwidthLimitRuleMapOutputWithContext(ctx context.Context) QosBandwidthLimitRuleMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(QosBandwidthLimitRuleMapOutput)
}

type QosBandwidthLimitRuleOutput struct {
	*pulumi.OutputState
}

func (QosBandwidthLimitRuleOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*QosBandwidthLimitRule)(nil))
}

func (o QosBandwidthLimitRuleOutput) ToQosBandwidthLimitRuleOutput() QosBandwidthLimitRuleOutput {
	return o
}

func (o QosBandwidthLimitRuleOutput) ToQosBandwidthLimitRuleOutputWithContext(ctx context.Context) QosBandwidthLimitRuleOutput {
	return o
}

func (o QosBandwidthLimitRuleOutput) ToQosBandwidthLimitRulePtrOutput() QosBandwidthLimitRulePtrOutput {
	return o.ToQosBandwidthLimitRulePtrOutputWithContext(context.Background())
}

func (o QosBandwidthLimitRuleOutput) ToQosBandwidthLimitRulePtrOutputWithContext(ctx context.Context) QosBandwidthLimitRulePtrOutput {
	return o.ApplyT(func(v QosBandwidthLimitRule) *QosBandwidthLimitRule {
		return &v
	}).(QosBandwidthLimitRulePtrOutput)
}

type QosBandwidthLimitRulePtrOutput struct {
	*pulumi.OutputState
}

func (QosBandwidthLimitRulePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**QosBandwidthLimitRule)(nil))
}

func (o QosBandwidthLimitRulePtrOutput) ToQosBandwidthLimitRulePtrOutput() QosBandwidthLimitRulePtrOutput {
	return o
}

func (o QosBandwidthLimitRulePtrOutput) ToQosBandwidthLimitRulePtrOutputWithContext(ctx context.Context) QosBandwidthLimitRulePtrOutput {
	return o
}

type QosBandwidthLimitRuleArrayOutput struct{ *pulumi.OutputState }

func (QosBandwidthLimitRuleArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]QosBandwidthLimitRule)(nil))
}

func (o QosBandwidthLimitRuleArrayOutput) ToQosBandwidthLimitRuleArrayOutput() QosBandwidthLimitRuleArrayOutput {
	return o
}

func (o QosBandwidthLimitRuleArrayOutput) ToQosBandwidthLimitRuleArrayOutputWithContext(ctx context.Context) QosBandwidthLimitRuleArrayOutput {
	return o
}

func (o QosBandwidthLimitRuleArrayOutput) Index(i pulumi.IntInput) QosBandwidthLimitRuleOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) QosBandwidthLimitRule {
		return vs[0].([]QosBandwidthLimitRule)[vs[1].(int)]
	}).(QosBandwidthLimitRuleOutput)
}

type QosBandwidthLimitRuleMapOutput struct{ *pulumi.OutputState }

func (QosBandwidthLimitRuleMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]QosBandwidthLimitRule)(nil))
}

func (o QosBandwidthLimitRuleMapOutput) ToQosBandwidthLimitRuleMapOutput() QosBandwidthLimitRuleMapOutput {
	return o
}

func (o QosBandwidthLimitRuleMapOutput) ToQosBandwidthLimitRuleMapOutputWithContext(ctx context.Context) QosBandwidthLimitRuleMapOutput {
	return o
}

func (o QosBandwidthLimitRuleMapOutput) MapIndex(k pulumi.StringInput) QosBandwidthLimitRuleOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) QosBandwidthLimitRule {
		return vs[0].(map[string]QosBandwidthLimitRule)[vs[1].(string)]
	}).(QosBandwidthLimitRuleOutput)
}

func init() {
	pulumi.RegisterOutputType(QosBandwidthLimitRuleOutput{})
	pulumi.RegisterOutputType(QosBandwidthLimitRulePtrOutput{})
	pulumi.RegisterOutputType(QosBandwidthLimitRuleArrayOutput{})
	pulumi.RegisterOutputType(QosBandwidthLimitRuleMapOutput{})
}