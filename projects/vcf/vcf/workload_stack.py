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
        self._provision_keypair()
        self._provision_deployment_network(True)
        self._provision_deployment_subnet(True)
        self._provision_router(True)
        self._provision_helper_vm()
        self._configure_helper_vm()

        self._provision_private_router()
        self._provision_private_networks()
        self._provision_reserved_names()
        self._provision_esxi_dns_recrods()
        self._provision_shares()

        self._provision_esxi_servers()
        for s in self.resources.esxi_servers:
            self._configure_esxi_server(s)
