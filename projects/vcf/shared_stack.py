import json
from os import name
import jinja2

import pulumi
from pulumi.config import ConfigMissingError
from pulumi.resource import ResourceOptions
from pulumi_openstack import compute, dns, networking
from types import SimpleNamespace

from provisioners import (
    ConnectionArgs,
    CopyFile,
    CopyFileFromString,
    RemoteExec,
)


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


class VCFSharedStack:
    def __init__(self, key_pair, provider_cloud_admin):
        self.config = pulumi.Config()
        self.stack_name = pulumi.get_stack()
        self.key_pair = key_pair
        self.provider_cloud_admin = provider_cloud_admin

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
            helper_vm=json.loads(self.config.require("helperVM")),
            public_router_name=self.config.require("publicRouter"),
            public_key_file=public_key_file,
            private_key_file=private_key_file,
        )
        mgmt_network = networking.get_network(name=self.props.mgmt_network["name"])
        mgmt_subnet = networking.get_subnet(name=self.props.mgmt_network["subnet_name"])
        self.resources = SimpleNamespace(
            mgmt_network=mgmt_network,
            mgmt_subnet=mgmt_subnet,
        )

    def proivsion(self):
        self._provision_network(protect=True)
        self._provision_helper_vm()

    def _provision_network(self, protect=False):
        deploy_network = networking.Network(
            self.props.deploy_network["name"],
            name=self.props.deploy_network["name"],
            opts=ResourceOptions(delete_before_replace=True, protect=protect),
        )
        deploy_subnet = networking.Subnet(
            self.props.deploy_network["subnet_name"],
            name=self.props.deploy_network["subnet_name"],
            network_id=deploy_network.id,
            cidr=self.props.deploy_network["cidr"],
            ip_version=4,
            opts=ResourceOptions(delete_before_replace=True, protect=protect),
        )
        public_router = networking.Router(
            self.props.public_router_name,
            name=self.props.public_router_name,
            external_network_id=self.props.external_network["id"],
            opts=ResourceOptions(delete_before_replace=True, protect=protect),
        )
        networking.RouterInterface(
            "router-interface-management",
            router_id=public_router.id,
            subnet_id=self.resources.mgmt_subnet.id,
            opts=ResourceOptions(
                provider=self.provider_cloud_admin,
                delete_before_replace=True,
                protect=protect,
            ),
        )
        networking.RouterInterface(
            "router-interface-deployement",
            router_id=public_router.id,
            subnet_id=deploy_subnet.id,
            opts=ResourceOptions(
                provider=self.provider_cloud_admin,
                delete_before_replace=True,
                protect=protect,
            ),
        )

    def _provision_helper_vm(self):
        init_script = r"""#!/bin/bash
echo 'net.ipv4.conf.default.rp_filter = 2' >> /etc/sysctl.conf
echo 'net.ipv4.conf.all.rp_filter = 2' >> /etc/sysctl.conf
/usr/sbin/sysctl -p /etc/sysctl.conf
"""
        sg = compute.SecGroup(
            "helper-vm-sg",
            description="allow ssh",
            rules=[
                compute.SecGroupRuleArgs(
                    cidr="0.0.0.0/0", from_port=22, to_port=22, ip_protocol="tcp"
                )
            ],
        )
        external_port = networking.Port(
            "helper-vm-external-port",
            network_id=self.resources.mgmt_network.id,
            fixed_ips=[
                networking.PortFixedIpArgs(
                    subnet_id=self.resources.mgmt_subnet.id,
                    ip_address=self.props.helper_vm["ip"],
                )
            ],
            security_group_ids=[sg.id],
        )
        helper_vm = compute.Instance(
            "helper-vm",
            name="helper-vm",
            flavor_id=self.props.helper_vm["flavor_id"],
            image_name=self.props.helper_vm["image_name"],
            networks=[
                compute.InstanceNetworkArgs(name=self.props.deploy_network["name"]),
            ],
            key_pair=self.key_pair.name,
            user_data=init_script,
            opts=ResourceOptions(
                delete_before_replace=True,
                ignore_changes=["image_name"],
            ),
        )
        attach_external_ip = compute.InterfaceAttach(
            "helper-vm-attatch",
            instance_id=helper_vm.id,
            port_id=external_port.id,
            opts=ResourceOptions(delete_before_replace=True, depends_on=[helper_vm]),
        )

        # configure helper vm
        conn_args = ConnectionArgs(
            host=self.props.helper_vm["ip"],
            username="ccloud",
            private_key_file=self.props.private_key_file,
        )
        exec_install_pwsh = RemoteExec(
            "install-powershell",
            host_id=helper_vm.id,
            conn=conn_args,
            commands=[
                "[ ! -f packages-microsoft-prod.deb ] && wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb || true",
                "sudo dpkg -i packages-microsoft-prod.deb",
                "sudo apt-get update",
                "echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections",
                "sudo apt-get install -y -q powershell",
                "pwsh -Command Set-PSRepository -Name 'PSGallery' -InstallationPolicy Trusted",
                "pwsh -Command Install-Module VMware.PowerCLI",
                "pwsh -Command Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Confirm:0",
                "pwsh -Command Set-PowerCLIConfiguration -Scope User -ParticipateInCEIP 0 -Confirm:0",
            ],
            opts=ResourceOptions(depends_on=[attach_external_ip]),
        )

        # send files to helper vm
        CopyFile(
            "copy-remove-vmk0",
            host_id=helper_vm.id,
            conn=conn_args,
            src="./scripts/cleanup.sh",
            dest="/home/ccloud/cleanup.sh",
            opts=ResourceOptions(depends_on=[attach_external_ip]),
        )
        CopyFile(
            "copy-rsa-key",
            host_id=helper_vm.id,
            conn=conn_args,
            src=self.props.private_key_file,
            dest="/home/ccloud/esxi_rsa",
            mode="600",
            opts=ResourceOptions(depends_on=[attach_external_ip]),
        )
        with open("./scripts/config.sh") as f:
            template = jinja2.Template(f.read())
            config_script = template.render(
                management_network=self.props.mgmt_network,
            )
            CopyFileFromString(
                "copy-config-sh",
                host_id=helper_vm.id,
                conn=conn_args,
                from_str=config_script,
                dest="/home/ccloud/config.sh",
                opts=ResourceOptions(depends_on=[attach_external_ip]),
            )
