import json

import pulumi
from pulumi import resource
from pulumi.config import ConfigMissingError
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi.stack_reference import StackReference
from pulumi_openstack import Provider, compute, dns, networking
from pulumi_openstack.compute import keypair

from provisioners import ConnectionArgs, RemoteExec

from .vcf_stack import VCFStack, resources_cache


class ManagementStack(VCFStack):
    def __init__(self, provider_cloud_admin, provider_ccadmin_master) -> None:
        super(ManagementStack, self).__init__(
            provider_cloud_admin, provider_ccadmin_master
        )

    def provision(self):
        super(ManagementStack, self).provision()
        self._provision_private_router()
        self._provision_private_networks()
        self._provision_reserved_names()
        self._provision_esxi_dns_recrods()
        self._provision_shares()

        self._provision_esxi_servers()
        for s in self.resources.esxi_servers:
            self._configure_esxi_server(s)
