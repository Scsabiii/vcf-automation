import json
from types import SimpleNamespace

import jinja2
import pulumi
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi_openstack import compute, dns, networking

from provisioners import ConnectionArgs, CopyFile, CopyFileFromString, RemoteExec


class SharedStack:
    def __init__(self, keypair_name, provider_cloud_admin):
        self.config = pulumi.Config()
        self.stack_name = pulumi.get_stack()
        self.provider_cloud_admin = provider_cloud_admin
        self.keypair = compute.get_keypair(keypair_name)

        private_key_file = (
            self.config.get("privateKeyFile") or "/pulumi/automation/etc/.ssh/id_rsa"
        )
        self.props = SimpleNamespace(
            external_network=json.loads(self.config.require("externalNetwork")),
            mgmt_network=json.loads(self.config.require("managementNetwork")),
            deploy_network=json.loads(self.config.require("deploymentNetwork")),
            helper_vm=json.loads(self.config.require("helperVM")),
            public_router_name=self.config.require("publicRouter"),
            keypair_name=keypair_name,
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
                # provider=self.provider_cloud_admin,
                delete_before_replace=True,
                protect=protect,
            ),
        )
        networking.RouterInterface(
            "router-interface-deployment",
            router_id=public_router.id,
            subnet_id=deploy_subnet.id,
            opts=ResourceOptions(
                # provider=self.provider_cloud_admin,
                delete_before_replace=True,
                protect=protect,
            ),
        )

        pulumi.export(
            "DeploymentNetwork",
            Output.all(deploy_network.name, deploy_network.id).apply(
                lambda args: f"{args[0]} ({args[1]})"
            ),
        )
        pulumi.export(
            "PublicRouter",
            Output.all(public_router.name, public_router.id).apply(
                lambda args: f"{args[0]} ({args[1]})"
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
            key_pair=self.props.keypair_name,
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

        # copy rsa key
        CopyFile(
            "copy-rsa-key",
            host_id=helper_vm.id,
            conn=conn_args,
            src=self.props.private_key_file,
            dest="/home/ccloud/esxi_rsa",
            mode="600",
            opts=ResourceOptions(depends_on=[attach_external_ip]),
        )

        # copy from path relative to the project root
        CopyFile(
            "copy-cleanup",
            host_id=helper_vm.id,
            conn=conn_args,
            src="./scripts/cleanup.sh",
            dest="/home/ccloud/cleanup.sh",
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

        pulumi.export(
            "HelperVM",
            Output.all(
                helper_vm.name, helper_vm.id, external_port.all_fixed_ips[0]
            ).apply(lambda args: f"{args[0]} ({args[1]}, {args[2]})"),
        )
