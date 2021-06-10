import json

import pulumi
from pulumi.invoke import InvokeOptions
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi_openstack import Provider, compute, dns, networking

from provisioners import ConnectionArgs, RemoteExec
from .vcf_stack import VCFStack, resources_cache


class WorkloadStack(VCFStack):
    def __init__(self, provider_cloud_admin, provider_ccadmin_master) -> None:
        super(WorkloadStack, self).__init__(
            provider_cloud_admin, provider_ccadmin_master
        )

    def provision(self):
        self._provision_private_router()
        self._provision_private_networks()
        self._provision_esxi_dns_recrods()

        self._provision_esxi_servers()
        for s in self.resources.esxi_servers:
            self._configure_esxi_server(s)

    @resources_cache("private_router")
    def _provision_private_router(self):
        return networking.Router(
            "private-router-" + self.stack_name,
            name="private-router-" + self.stack_name,
            opts=ResourceOptions(delete_before_replace=True),
        )

    @resources_cache("private_networks")
    def _provision_private_networks(self):
        private_networks = {}
        for props in self.props.private_networks:
            network = networking.Network("private-network-" + props["name"])
            subnet = networking.Subnet(
                "subnet-" + props["name"],
                network_id=network.id,
                cidr=props["cidr"],
                ip_version=4,
                opts=ResourceOptions(delete_before_replace=True),
            )
            networking.RouterInterface(
                "router-interface-" + props["name"],
                router_id=self.resources.private_router.id,
                subnet_id=subnet.id,
                opts=ResourceOptions(delete_before_replace=True),
            )
            private_networks[props["name"]] = {
                "network": network,
                "subnet": subnet,
                "vlan_id": props["vlan_id"],
            }
        return private_networks

    def _provision_dns_record(self, dns_name, ipaddr):
        dns_zone = dns.get_dns_zone(name=self.props.dns_zone_name)
        reverse_dns_zone = dns.get_dns_zone(
            name=self.props.reverse_dns_zone_name,
            opts=InvokeOptions(provider=self.provider_ccadmin_master),
        )
        dns_name = dns_name + "." + self.props.dns_zone_name
        r = dns.RecordSet(
            dns_name,
            name=dns_name,
            records=[ipaddr],
            type="A",
            ttl=1800,
            zone_id=dns_zone.id,
            opts=ResourceOptions(delete_before_replace=True),
        )
        rr = dns.RecordSet(
            "reverse-" + dns_name,
            name=ipaddr.split(".")[-1] + "." + self.props.reverse_dns_zone_name,
            records=[dns_name],
            type="PTR",
            ttl=1800,
            zone_id=reverse_dns_zone.id,
            opts=ResourceOptions(
                provider=self.provider_ccadmin_master,
                delete_before_replace=True,
                depends_on=[r],
            ),
        )

    def _provision_esxi_dns_recrods(self):
        for n in self.props.esxi_nodes:
            node_name, node_ip = n["name"], n["ip"]
            self._provision_dns_record("esxi-" + node_name, node_ip)
