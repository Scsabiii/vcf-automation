"""An OpenStack Python Pulumi program"""

import pulumi
from pulumi.resource import ResourceOptions
from pulumi_openstack import compute, networking

# Create an OpenStack resource (Compute Instance)

privateNetworkProps = [
    {"name":"vmotion",      "cidr":"10.180.0.0/24"},
    {"name":"edgetep",      "cidr":"10.180.1.0/24"},
    {"name":"hosttep",      "cidr":"10.180.2.0/24"},
    {"name":"nfs",          "cidr":"10.180.3.0/24"},
    {"name":"vsan",         "cidr":"10.180.4.0/24"},
    {"name":"vsan-witness", "cidr":"10.180.5.0/24"},
    {"name":"deployment",   "cidr":"10.180.6.0/24"},
]

externalNetworkId = "430991b3-da0d-41cb-ac54-d1d532841725"
esxiImage = "vsphere-7.0U1c-amd64-baremetal"
esxiImageID = "de2dacf0-c76d-44c0-8979-472085a7512e"

esxiFlavor = "m1.medium"
esxiNodes = [
    {
        "node": "node009-bb095",
        "id": "6bdb7232-8e58-4d28-a4dc-3e3bd781c2b8"
    },
]

floatingIPPool = "FloatingIP-external-monsoon3-03"
floatingIPPoolSubnetID = "ac736737-1969-4e2c-9f6d-81b8b5278dd7"


# create private netowrks with preix "private-network-"
# and router "private-router"
privateNetworks = {}
for props in privateNetworkProps:
    network = networking.Network(props["name"],
                                 name="private-network-"+props["name"])
    subnet = networking.Subnet("subnet-"+props["name"],
                               name="subnet-"+props["name"],
                               network_id=network.id,
                               cidr=props["cidr"],
                               ip_version=4)
    privateNetworks[props["name"]] = {"network": network, "subnet": subnet}

privateRouter = networking.Router("management-private-router",
                                  external_network_id=externalNetworkId)

for name, network in privateNetworks.items():
    networking.RouterInterface("router-interface-"+name,
                               router_id=privateRouter.id,
                               subnet_id=network["subnet"].id)

# esxi installation
# nodeName = "esxi-"+esxiNodes[0]["node"]
# nodeHint = "::"+esxiNodes[0]["id"]
# instance = compute.Instance(nodeName, 
#                             name=nodeName,
#                             availability_zone_hints=nodeHint,
#                             flavor_name=esxiFlavor,
#                             image_id=esxiImageID,
#                             networks=[
#                                 compute.InstanceNetworkArgs(
#                                     name=privateNetworks["deployment"]["network"].name,
#                                     fixed_ip_v4="10.180.6.6"),
#                             ])

# create a compute instance and a helper vm
# execute powershell script on the vm 
compute_instance = compute.Instance("compute-instance",
                                    flavor_name="m1.small",
                                    image_id="0db3dc89-671a-4745-a35f-fe6b99ea3d8a", #ubuntu 20.04
                                    networks=[
                                        compute.InstanceNetworkArgs(
                                            name=privateNetworks["deployment"]["network"].name,
                                            fixed_ip_v4="10.180.6.10",
                                        )
                                    ])
helper_instance = compute.Instance("helper-instance",
                                    flavor_name="m1.small",
                                    image_id="0db3dc89-671a-4745-a35f-fe6b99ea3d8a", #ubuntu 20.04
                                    networks=[
                                        compute.InstanceNetworkArgs(
                                            name=privateNetworks["deployment"]["network"].name,
                                        )
                                    ],
                                    opts=pulumi.ResourceOptions(
                                         depends_on=[compute_instance]
                                    ))
instance_fip = networking.FloatingIp('helper-fip',
                                     pool=floatingIPPool,
                                     subnet_id=floatingIPPoolSubnetID)
compute.FloatingIpAssociate('helper-fip-associate',
                           fixed_ip=helper_instance.access_ip_v4,
                           floating_ip=instance_fip.address,
                           instance_id=helper_instance.id)





# Export the IP of the instance
pulumi.export("vmotionSubnetName", privateNetworks["vmotion"]["subnet"].name)
pulumi.export("vmotionSubnetID", privateNetworks["vmotion"]["subnet"].id)

