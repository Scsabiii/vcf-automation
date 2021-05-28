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


from management_stack import ManagementStack
from shared_stack import VCFSharedStack
from workload_stack import WorkloadStack


# stack
stack_name = pulumi.get_stack()

# read config
config = pulumi.Config()

###################################################################################
# cloud admin provider
###################################################################################
openstack_config = pulumi.Config("openstack")
user_name = openstack_config.require("userName")
password = openstack_config.require_secret("password")
auth_url = openstack_config.require("authUrl")
region = openstack_config.require("region")
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
# key pair
###################################################################################
public_key_file = "/pulumi/avocado/etc/.ssh/id_rsa.pub"
private_key_file = "/pulumi/avocado/etc/.ssh/id_rsa"
with open(public_key_file) as f:
    pk = f.read()
    key_pair = compute.Keypair("rsa-keypair", public_key=pk)

###################################################################################
# public networks
###################################################################################
externalNetworkProps = json.loads(config.require("externalNetwork"))
managementNetworkPorps = json.loads(config.require("managementNetwork"))

if stack_name == "shared":
    ss = VCFSharedStack(key_pair, provider_cloud_admin)
    ss.proivsion()
    exit(0)

if stack_name in ("management", "dev", "ap002"):
    ms = ManagementStack(key_pair, provider_ccadmin_master)
    ms.provision()
    exit(0)

if stack_name == "workload":
    ws = WorkloadStack(key_pair, provider_ccadmin_master)
    ws.provision()
    exit(0)


###################################################################################
# load more settings
###################################################################################
helper_vm = json.loads(config.require("helperVM"))


###################################################################################
# register dns records
###################################################################################

####################################################################################
## helper vm
####################################################################################
####################################################################################
## copy file to helper
####################################################################################

# Export the IP of the instance
# pulumi.export("DapnIPNetworkID", management_network["network"].id)
# pulumi.export("DapnIPSubnetID", management_network["subnet"].id)
# pulumi.export("helperVMIP", port_1.all_fixed_ips[0])

# pulumi.export("exec_2", exec_2.results)
