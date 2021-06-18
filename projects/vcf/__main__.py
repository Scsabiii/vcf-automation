"""An OpenStack Python Pulumi program"""

import pulumi
from pulumi_openstack import Provider

from vcf import ManagementStack, WorkloadStack


# stack
config = pulumi.Config()
stack_name = pulumi.get_stack()
stack_type = config.require("stackType")

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
