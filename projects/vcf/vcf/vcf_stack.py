import json
from types import SimpleNamespace

import pulumi
from pulumi.config import ConfigMissingError
from pulumi.resource import ResourceOptions
from pulumi_openstack import compute, dns, networking


def resources_cache(name):
    def inner(fn):
        def wrapper(self, *args, **kwargs):
            res = fn(self, *args, **kwargs)
            setattr(self.resources, name, res)
            return res

        return wrapper

    return inner


class VCFStack:
    def __init__(self) -> None:
        self.config = pulumi.Config()
        self.stack_name = pulumi.get_stack()

        try:
            private_networks = json.loads(self.config.require("privateNetworks"))
        except ConfigMissingError:
            private_networks = []
        try:
            esxi_nodes = json.loads(self.config.require("esxiNodes"))
        except ConfigMissingError:
            esxi_nodes = []
        try:
            reserved_ips = json.loads(self.config.require("reservedIPs"))
        except ConfigMissingError:
            reserved_ips = []
        public_key_file = (
            self.config.get("publicKeyFile") or "/pulumi/automation/etc/.ssh/id_rsa.pub"
        )
        private_key_file = (
            self.config.get("privateKeyFile") or "/pulumi/automation/etc/.ssh/id_rsa"
        )

        self.props = SimpleNamespace(
            external_network=json.loads(self.config.require("externalNetwork")),
            mgmt_network=json.loads(self.config.require("managementNetwork")),
            deploy_network=json.loads(self.config.require("deploymentNetwork")),
            dns_zone_name=self.config.require("dnsZoneName"),
            reverse_dns_zone_name=self.config.require("reverseDnsZoneName"),
            helper_vm=json.loads(self.config.require("helperVM")),
            private_networks=private_networks,
            esxi_nodes=esxi_nodes,
            reserved_ips=reserved_ips,
            esxi_image=self.config.get("esxiServerImage"),
            esxi_flavor_id=self.config.get("esxiServerFlavorID"),
            public_key_file=public_key_file,
            private_key_file=private_key_file,
        )

        deploy_network = networking.get_network(name=self.props.deploy_network["name"])
        deploy_subnet = networking.get_subnet(
            name=self.props.deploy_network["subnet_name"],
            network_id=deploy_network.id,
            cidr=self.props.deploy_network["cidr"],
            ip_version=4,
        )
        mgmt_network = networking.get_network(name=self.props.mgmt_network["name"])
        mgmt_subnet = networking.get_subnet(name=self.props.mgmt_network["subnet_name"])

        self.resources = SimpleNamespace(
            deploy_network=deploy_network,
            deploy_subnet=deploy_subnet,
            mgmt_network=mgmt_network,
            mgmt_subnet=mgmt_subnet,
        )
