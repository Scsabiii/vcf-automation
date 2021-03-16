// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package networking

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pulumi/pulumi-openstack/sdk/v2/go/openstack"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type module struct {
	version semver.Version
}

func (m *module) Version() semver.Version {
	return m.version
}

func (m *module) Construct(ctx *pulumi.Context, name, typ, urn string) (r pulumi.Resource, err error) {
	switch typ {
	case "openstack:networking/addressScope:AddressScope":
		r, err = NewAddressScope(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/floatingIp:FloatingIp":
		r, err = NewFloatingIp(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/floatingIpAssociate:FloatingIpAssociate":
		r, err = NewFloatingIpAssociate(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/network:Network":
		r, err = NewNetwork(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/port:Port":
		r, err = NewPort(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/portSecGroupAssociate:PortSecGroupAssociate":
		r, err = NewPortSecGroupAssociate(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/qosBandwidthLimitRule:QosBandwidthLimitRule":
		r, err = NewQosBandwidthLimitRule(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/qosDscpMarkingRule:QosDscpMarkingRule":
		r, err = NewQosDscpMarkingRule(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/qosMinimumBandwidthRule:QosMinimumBandwidthRule":
		r, err = NewQosMinimumBandwidthRule(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/qosPolicy:QosPolicy":
		r, err = NewQosPolicy(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/quotaV2:QuotaV2":
		r, err = NewQuotaV2(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/rbacPolicyV2:RbacPolicyV2":
		r, err = NewRbacPolicyV2(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/router:Router":
		r, err = NewRouter(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/routerInterface:RouterInterface":
		r, err = NewRouterInterface(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/routerRoute:RouterRoute":
		r, err = NewRouterRoute(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/secGroup:SecGroup":
		r, err = NewSecGroup(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/secGroupRule:SecGroupRule":
		r, err = NewSecGroupRule(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/subnet:Subnet":
		r, err = NewSubnet(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/subnetPool:SubnetPool":
		r, err = NewSubnetPool(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/subnetRoute:SubnetRoute":
		r, err = NewSubnetRoute(ctx, name, nil, pulumi.URN_(urn))
	case "openstack:networking/trunk:Trunk":
		r, err = NewTrunk(ctx, name, nil, pulumi.URN_(urn))
	default:
		return nil, fmt.Errorf("unknown resource type: %s", typ)
	}

	return
}

func init() {
	version, err := openstack.PkgVersion()
	if err != nil {
		fmt.Println("failed to determine package version. defaulting to v1: %v", err)
	}
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/addressScope",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/floatingIp",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/floatingIpAssociate",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/network",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/port",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/portSecGroupAssociate",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/qosBandwidthLimitRule",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/qosDscpMarkingRule",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/qosMinimumBandwidthRule",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/qosPolicy",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/quotaV2",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/rbacPolicyV2",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/router",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/routerInterface",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/routerRoute",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/secGroup",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/secGroupRule",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/subnet",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/subnetPool",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/subnetRoute",
		&module{version},
	)
	pulumi.RegisterResourceModule(
		"openstack",
		"networking/trunk",
		&module{version},
	)
}
