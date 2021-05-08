"""An OpenStack Python Pulumi program"""

import datetime
from types import MappingProxyType
from jinja2 import Template
import json
import pulumi
from pulumi.config import ConfigMissingError
from pulumi.invoke import InvokeOptions
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi_openstack import Provider, provider
from pulumi_openstack import compute, networking, dns

from provisioners import (
    ConnectionArgs,
    CopyFile,
    CopyFileFromString,
    RemoteExec,
)

# read config
config = pulumi.Config()
# public_key = config.require("publicKey")
# private_key = config.require_secret("privateKey")
public_key_file = "/pulumi/avocado/etc/.ssh/id_rsa.pub"
private_key_file = "/pulumi/avocado/etc/.ssh/id_rsa"

privateNetworkProps = json.loads(config.require("privateNetworks"))
deploymentNetworkProps = json.loads(config.require("deploymentNetwork"))
managementNetworkPorps = json.loads(config.require("managementNetwork"))
management_network_subnet_gateway = managementNetworkPorps["subnet_gateway"]
management_network_subnet_mask = managementNetworkPorps["subnet_mask"]

externalNetworkProps = json.loads(config.require("externalNetwork"))
external_network = externalNetworkProps["name"]
external_network_id = externalNetworkProps["id"]

esxi_image = config.require("esxiServerImange")
esxi_flavor_id = config.require("esxiServerFlavorID")
esxi_nodes = json.loads(config.require("esxiNodes"))

helper_vm = json.loads(config.require("helperVM"))
dns_zone_name = config.require("dnsZoneName")
reverse_dns_zone_name = config.require("reverseDnsZoneName")

try:
    reserved_ips = json.loads(config.require("reservedIPs"))
except ConfigMissingError:
    reserved_ips = []

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

###################################################################################
# provider: ccadmin/master
###################################################################################
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
# ssh key
###################################################################################
key_pair = compute.Keypair("rsa-keypair", public_key=public_key)

###################################################################################
# public networks
###################################################################################
publicRouter = networking.Router(
    "mgmtdomain-public-router",
    external_network_id=external_network_id,
    opts=ResourceOptions(delete_before_replace=True),
)

network_1 = networking.Network(deploymentNetworkProps["name"] + "-network")
subnet_1 = networking.Subnet(
    "subnet-" + deploymentNetworkProps["name"],
    network_id=network_1.id,
    cidr=deploymentNetworkProps["cidr"],
    ip_version=4,
    opts=ResourceOptions(delete_before_replace=True),
)
networking.RouterInterface(
    "router-interface-deployment",
    router_id=publicRouter.id,
    subnet_id=subnet_1.id,
    opts=ResourceOptions(provider=provider_cloud_admin, delete_before_replace=True),
)

network_2 = networking.get_network(name=managementNetworkPorps["name"])
subnet_2 = networking.get_subnet(name=managementNetworkPorps["subnet_name"])
networking.RouterInterface(
    "router-interface-management",
    router_id=publicRouter.id,
    subnet_id=subnet_2.id,
    opts=ResourceOptions(provider=provider_cloud_admin, delete_before_replace=True),
)
# NOTE: router interface needs cloud admin provider, since DapnIP network is
# not owned by current project

deployment_network = {"network": network_1, "subnet": subnet_1}
management_network = {
    "network": network_2,
    "subnet": subnet_2,
    "vlan_id": managementNetworkPorps["vlan_id"],
}

###################################################################################
# private networks
###################################################################################
privateRouter = networking.Router(
    "mgmtdomain-private-router",
    opts=ResourceOptions(delete_before_replace=True),
)

private_networks = {}
for props in privateNetworkProps:
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
        router_id=privateRouter.id,
        subnet_id=subnet.id,
        opts=ResourceOptions(delete_before_replace=True),
    )
    private_networks[props["name"]] = {
        "network": network,
        "subnet": subnet,
        "vlan_id": props["vlan_id"],
    }


###################################################################################
# register dns records
###################################################################################
dns_zone = dns.get_dns_zone(name=dns_zone_name)
reverse_dns_zone = dns.get_dns_zone(
    name=reverse_dns_zone_name, opts=InvokeOptions(provider=provider_ccadmin_master)
)
for r in reserved_ips:
    ipaddr, name = r["ip"], r["name"]
    networking.Port(
        "reserved-port-" + ipaddr,
        network_id=management_network["network"].id,
        fixed_ips=[
            networking.PortFixedIpArgs(
                subnet_id=management_network["subnet"].id,
                ip_address=ipaddr,
            )
        ],
        opts=ResourceOptions(delete_before_replace=True),
    )
    dns_name = name + "." + dns_zone_name
    dns.RecordSet(
        dns_name,
        name=dns_name,
        records=[ipaddr],
        type="A",
        ttl=1800,
        zone_id=dns_zone.id,
        opts=ResourceOptions(delete_before_replace=True),
    )
    dns.RecordSet(
        "reverse-" + dns_name,
        name=ipaddr.split(".")[-1] + "." + reverse_dns_zone_name,
        records=[dns_name],
        type="PTR",
        ttl=1800,
        zone_id=reverse_dns_zone.id,
        opts=ResourceOptions(
            provider=provider_ccadmin_master,
            delete_before_replace=True,
        ),
    )
for n in esxi_nodes:
    node_name, node_ip = n["name"], n["ip"]
    dns_name = node_name + "." + dns_zone_name
    dns.RecordSet(
        dns_name,
        name=dns_name,
        records=[node_ip],
        type="A",
        ttl=1800,
        zone_id=dns_zone.id,
        opts=ResourceOptions(delete_before_replace=True),
    )
    dns.RecordSet(
        "reverse-" + dns_name,
        name=node_ip.split(".")[-1] + "." + reverse_dns_zone_name,
        records=[dns_name],
        type="PTR",
        ttl=1800,
        zone_id=reverse_dns_zone.id,
        opts=ResourceOptions(
            provider=provider_ccadmin_master,
            delete_before_replace=True,
        ),
    )


###################################################################################
# helper vm
###################################################################################
helper_image = helper_vm["image_name"]
helper_flavor = helper_vm["flavor_id"]
helper_ip = helper_vm["ip"]

sg = compute.SecGroup(
    "helper-vm-sg",
    description="allow ssh",
    rules=[
        compute.SecGroupRuleArgs(
            cidr="0.0.0.0/0", from_port=22, to_port=22, ip_protocol="tcp"
        )
    ],
)
port_1 = networking.Port(
    "helper-vm-external-port",
    network_id=management_network["network"].id,
    fixed_ips=[
        networking.PortFixedIpArgs(
            subnet_id=management_network["subnet"].id, ip_address=helper_ip
        )
    ],
    security_group_ids=[sg.id],
)
init_script = r"""#!/bin/bash
echo 'net.ipv4.conf.default.rp_filter = 2' >> /etc/sysctl.conf
echo 'net.ipv4.conf.all.rp_filter = 2' >> /etc/sysctl.conf
/usr/sbin/sysctl -p /etc/sysctl.conf
"""
helper_vm = compute.Instance(
    "helper-vm",
    flavor_id=helper_flavor,
    image_name=helper_image,
    networks=[
        compute.InstanceNetworkArgs(name=deployment_network["network"].name),
    ],
    key_pair=key_pair.name,
    user_data=init_script,
    opts=ResourceOptions(depends_on=[key_pair], delete_before_replace=True),
)
helper_vm_attach_external_ip = compute.InterfaceAttach(
    "helper_vm_attatch",
    instance_id=helper_vm.id,
    port_id=port_1.id,
    opts=ResourceOptions(delete_before_replace=True),
)
helper_vm_external_ip = port_1.all_fixed_ips[0]

###################################################################################
# copy file to helper
###################################################################################
conn_args = ConnectionArgs(
    host=helper_vm_external_ip,
    username="ccloud",
    private_key=private_key,
)
exec_install_pwsh = RemoteExec(
    "install-powershell",
    host_id=helper_vm.id,
    conn=conn_args,
    commands=[
        "[ ! -f packages-microsoft-prod.deb ] && wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb",
        "sudo dpkg -i packages-microsoft-prod.deb",
        "sudo apt-get update",
        "sudo apt-get install -y powershell",
        "pwsh -Command Set-PSRepository -Name 'PSGallery' -InstallationPolicy Trusted",
        "pwsh -Command Install-Module VMware.PowerCLI",
    ],
    opts=ResourceOptions(depends_on=[helper_vm, helper_vm_attach_external_ip]),
)

with open("./config.sh") as f:
    template = Template(f.read())
    config_script = template.render(
        private_networks=privateNetworkProps, management_network=managementNetworkPorps
    )
copy_1 = CopyFileFromString(
    "copy_file_1",
    host_id=helper_vm.id,
    conn=conn_args,
    from_str=config_script,
    dest="/home/ccloud/config.sh",
    opts=ResourceOptions(depends_on=[helper_vm, helper_vm_attach_external_ip]),
)
###################################################################################
# esxi installation
###################################################################################
for n in esxi_nodes:
    node_name, node_id, node_ip = n["name"], n["id"], n["ip"]

    parent_port = networking.Port(
        node_name + "-deployment",
        network_id=deployment_network["network"].id,
        opts=ResourceOptions(depends_on=[deployment_network["subnet"]]),
    )

    if n.get("image_name") is not None:
        esxi_instance = compute.Instance(
            "esxi-" + node_name,
            name="esxi-" + node_name,
            availability_zone_hints=f"::{node_id}",
            flavor_id=esxi_flavor_id,
            image_name=n.get("image_name"),
            networks=[compute.InstanceNetworkArgs(port=parent_port.id)],
            key_pair=key_pair.name,
            opts=ResourceOptions(delete_before_replace=True),
        )
    else:
        esxi_instance = compute.Instance(
            "esxi-" + node_name,
            name="esxi-" + node_name,
            availability_zone_hints=f"::{node_id}",
            flavor_id=esxi_flavor_id,
            image_name=esxi_image,
            networks=[compute.InstanceNetworkArgs(port=parent_port.id)],
            key_pair=key_pair.name,
            opts=ResourceOptions(delete_before_replace=True),
        )

    subport_vmotion = networking.Port(
        node_name + "-vmotion",
        admin_state_up=True,
        network_id=private_networks["vmotion"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["vmotion"]["subnet"]]),
    )
    subport_edgetep = networking.Port(
        node_name + "-edgetep",
        network_id=private_networks["edgetep"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["edgetep"]["subnet"]]),
    )
    subport_hosttep = networking.Port(
        node_name + "-hosttep",
        network_id=private_networks["hosttep"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["hosttep"]["subnet"]]),
    )
    subport_nfs = networking.Port(
        node_name + "-nfs",
        network_id=private_networks["nfs"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["nfs"]["subnet"]]),
    )
    subport_vsan = networking.Port(
        node_name + "-vsan",
        network_id=private_networks["vsan"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["vsan"]["subnet"]]),
    )
    subport_vsanwitness = networking.Port(
        node_name + "-vsanwitness",
        network_id=private_networks["vsanwitness"]["network"].id,
        opts=ResourceOptions(depends_on=[private_networks["vsanwitness"]["subnet"]]),
    )
    subport_management = networking.Port(
        node_name + "-management-vcf01",
        network_id=management_network["network"].id,
        fixed_ips=[
            networking.PortFixedIpArgs(
                subnet_id=management_network["subnet"].id, ip_address=node_ip
            )
        ],
    )

    trunk = networking.trunk.Trunk(
        node_name + "-trunk",
        name=node_name + "-trunk",
        admin_state_up=True,
        port_id=parent_port.id,
        sub_ports=[
            networking.TrunkSubPortArgs(
                port_id=subport_vmotion.id,
                segmentation_id=private_networks["vmotion"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_edgetep.id,
                segmentation_id=private_networks["edgetep"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_hosttep.id,
                segmentation_id=private_networks["hosttep"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_nfs.id,
                segmentation_id=private_networks["nfs"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_vsan.id,
                segmentation_id=private_networks["vsan"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_vsanwitness.id,
                segmentation_id=private_networks["vsanwitness"]["vlan_id"],
                segmentation_type="vlan",
            ),
            networking.TrunkSubPortArgs(
                port_id=subport_management.id,
                segmentation_id=management_network["vlan_id"],
                segmentation_type="vlan",
            ),
        ],
        opts=ResourceOptions(depends_on=[esxi_instance]),
    )

    command_str = Output.all(
        parent_port.all_fixed_ips, subport_management.all_fixed_ips
    ).apply(
        lambda ports: "pwsh /home/ccloud/config.sh -LocalIP {} -IP {} -Gateway {} -Netmask {} 2>/tmp/error.txt; cat /tmp/error.txt >&2".format(
            ports[0][0],
            ports[1][0],
            management_network_subnet_gateway,
            management_network_subnet_mask,
        )
    )
    RemoteExec(
        "configure-esxi-host-" + node_name,
        host_id=esxi_instance.id,
        conn=conn_args,
        commands=[command_str],
        opts=ResourceOptions(
            depends_on=[copy_1, esxi_instance, trunk, exec_install_pwsh]
        ),
    )

    pulumi.export("EsxiHostIP-" + node_name, subport_management.all_fixed_ips[0])


# exec_2 = RemoteExec(
#     "configure-esxi-host",
#     host_id=helper_vm.id,
#     conn=conn_args,
#     commands=[
#         # redirect powershell script error to a file then to error stream
#         "pwsh /home/ccloud/config.sh -LocalIP {} -IP {} -Gateway {} -Netmask {} 2>/tmp/error.txt; cat /tmp/error.txt >&2".format(
#             "10.180.6.5",
#             "10.237.209.20",
#             "10.180.6.5",
#             "10.180.6.5",
#         )
#     ],
#     opts=ResourceOptions(depends_on=[copy_1]),
# )
# pwsh config.sh -LocalIP 10.180.6.5 -IP 10.237.209.20 -Gateway 10.237.209.1 -Netmask 255.255.255.128


# Export the IP of the instance
# pulumi.export("vmotionSubnetName",       privateNetworks["vmotion"]["subnet"].name)
# pulumi.export("vmotionSubnetID",         privateNetworks["vmotion"]["subnet"].id)
pulumi.export("DapnIPNetworkID", management_network["network"].id)
pulumi.export("DapnIPSubnetID", management_network["subnet"].id)
pulumi.export("helperVMIP", port_1.all_fixed_ips[0])
pulumi.export("deploymentNetwork", deploymentNetworkProps["name"])

# pulumi.export("exec_2", exec_2.results)
