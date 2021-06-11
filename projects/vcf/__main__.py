"""An OpenStack Python Pulumi program"""

import datetime
from types import MappingProxyType
from jinja2 import Template
import json
import pulumi
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi_openstack import Provider, provider
from pulumi_openstack import compute, networking, dns
from pulumi_openstack.networking.get_network import get_network
from pulumi_openstack.networking.get_subnet import get_subnet


from vcf import ManagementStack, SharedStack, WorkloadStack


# stack
stack_name = pulumi.get_stack()
stack_type = stack_name.split("-")[0]

###################################################################################
# ccadmin/cloud_admin and ccadmin/master provider
###################################################################################
openstack_config = pulumi.Config("openstack")
auth_url = openstack_config.require("authUrl")
region = openstack_config.require("region")
user_name = openstack_config.require("userName")
password = openstack_config.require_secret("password")
provider_cloud_admin = Provider(
    "cloud_admin",
    user_name=user_name,
    password=password,
    auth_url=auth_url,
    insecure=True,
    project_domain_name="ccadmin",
    user_domain_name="ccadmin",
    tenant_name="cloud_admin",
)
provider_ccadmin_master = Provider(
    "ccadmin_master",
    user_name=user_name,
    password=password,
    auth_url=auth_url,
    insecure=True,
    project_domain_name="ccadmin",
    user_domain_name="ccadmin",
    tenant_name="master",
)


###################################################################################
# provision
###################################################################################
if stack_type == "management":
    ms = ManagementStack(provider_cloud_admin, provider_ccadmin_master)
    ms.provision()
    exit(0)

if stack_type == "workload":
    ws = WorkloadStack(provider_cloud_admin, provider_ccadmin_master)
    ws.provision()
    exit(0)
