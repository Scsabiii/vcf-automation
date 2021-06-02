import json

import pulumi
from pulumi.invoke import InvokeOptions
from pulumi.output import Output
from pulumi.resource import ResourceOptions
from pulumi_openstack import Provider, compute, dns, networking

from provisioners import ConnectionArgs, RemoteExec
from .vcf_stack import VCFStack, resources_cache


class WorkloadStack(VCFStack):
    def __init__(self, key_pair, provider_ccadmin_master) -> None:
        super(WorkloadStack, self).__init__()
        self.key_pair = key_pair
        self.provider_ccadmin_master = provider_ccadmin_master

    def provision(self):
        self._provision_private_router()
        self._provision_private_networks()
        self._provision_esxi_dns_recrods()
        self._provision_esxi_servers()

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

    def _provision_esxi_servers(self):
        for n in self.props.esxi_nodes:
            node_name, node_id, node_ip = n["name"], n["id"], n["ip"]
            parent_port = networking.Port(
                node_name + "-deployment",
                network_id=self.resources.deploy_network.id,
            )
            instance = compute.Instance(
                "esxi-" + node_name,
                name="esxi-" + node_name,
                availability_zone_hints=f"::{node_id}",
                flavor_id=self.props.esxi_flavor_id,
                image_name=self.props.esxi_image,
                networks=[compute.InstanceNetworkArgs(port=parent_port.id)],
                key_pair=self.key_pair.name,
                opts=ResourceOptions(
                    delete_before_replace=True, ignore_changes=["image_name"]
                ),
            )
            subport_vmotion = networking.Port(
                node_name + "-vmotion",
                admin_state_up=True,
                network_id=self.resources.private_networks["vmotion"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[self.resources.private_networks["vmotion"]["subnet"]]
                ),
            )
            subport_edgetep = networking.Port(
                node_name + "-edgetep",
                network_id=self.resources.private_networks["edgetep"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[self.resources.private_networks["edgetep"]["subnet"]]
                ),
            )
            subport_hosttep = networking.Port(
                node_name + "-hosttep",
                network_id=self.resources.private_networks["hosttep"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[self.resources.private_networks["hosttep"]["subnet"]]
                ),
            )
            subport_nfs = networking.Port(
                node_name + "-nfs",
                network_id=self.resources.private_networks["nfs"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[self.resources.private_networks["nfs"]["subnet"]]
                ),
            )
            subport_vsan = networking.Port(
                node_name + "-vsan",
                network_id=self.resources.private_networks["vsan"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[self.resources.private_networks["vsan"]["subnet"]]
                ),
            )
            subport_vsanwitness = networking.Port(
                node_name + "-vsanwitness",
                network_id=self.resources.private_networks["vsanwitness"]["network"].id,
                opts=ResourceOptions(
                    depends_on=[
                        self.resources.private_networks["vsanwitness"]["subnet"]
                    ]
                ),
            )
            subport_management = networking.Port(
                node_name + "-management-vcf01",
                network_id=self.resources.mgmt_network.id,
                fixed_ips=[
                    networking.PortFixedIpArgs(
                        subnet_id=self.resources.mgmt_subnet.id, ip_address=node_ip
                    )
                ],
            )
            pn = self.resources.private_networks
            trunk = networking.trunk.Trunk(
                node_name + "-trunk",
                name=node_name + "-trunk",
                admin_state_up=True,
                port_id=parent_port.id,
                sub_ports=[
                    networking.TrunkSubPortArgs(
                        port_id=subport_vmotion.id,
                        segmentation_id=pn["vmotion"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_edgetep.id,
                        segmentation_id=pn["edgetep"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_hosttep.id,
                        segmentation_id=pn["hosttep"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_nfs.id,
                        segmentation_id=pn["nfs"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_vsan.id,
                        segmentation_id=pn["vsan"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_vsanwitness.id,
                        segmentation_id=pn["vsanwitness"]["vlan_id"],
                        segmentation_type="vlan",
                    ),
                    networking.TrunkSubPortArgs(
                        port_id=subport_management.id,
                        segmentation_id=self.props.mgmt_network["vlan_id"],
                        segmentation_type="vlan",
                    ),
                ],
                opts=ResourceOptions(depends_on=[instance]),
            )

            self._configure_esxi_node(instance, node_name, node_ip)

    def _configure_esxi_node(self, instance, node_name, node_ip):
        # set password
        command_set_passwd = instance.access_ip_v4.apply(
            lambda local_ip: (
                "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR "
                "-i /home/ccloud/esxi_rsa root@{} 'echo VMware1!VMware1! | passwd --stdin root'"
            ).format(local_ip)
        )
        # config node
        command_config = instance.access_ip_v4.apply(
            lambda local_ip: "pwsh /home/ccloud/config.sh -LocalIP {} -IP {} -Gateway {} -Netmask {}".format(
                local_ip,
                node_ip,
                self.props.mgmt_network["subnet_gateway"],
                self.props.mgmt_network["subnet_mask"],
            )
        )
        # remove vmk0
        command_cleanup = "pwsh /home/ccloud/cleanup.sh -HostIP {}".format(node_ip)

        # connection
        conn_helper_args = ConnectionArgs(
            host=self.props.helper_vm["ip"],
            username="ccloud",
            private_key_file=self.props.private_key_file,
        )
        conn_esxi_args = ConnectionArgs(
            host=node_ip,
            username="root",
            private_key_file=self.props.private_key_file,
        )

        # execution
        step_1 = RemoteExec(
            "configure-" + node_name + "-step-1",
            host_id=instance.id,
            conn=conn_helper_args,
            commands=[command_set_passwd],
        )
        step_2 = RemoteExec(
            "configure-" + node_name + "-step-2",
            host_id=instance.id,
            conn=conn_helper_args,
            commands=[command_config],
            opts=ResourceOptions(depends_on=[step_1]),
        )
        step_3 = RemoteExec(
            "configure-" + node_name + "-step-3",
            host_id=instance.id,
            conn=conn_esxi_args,
            commands=[
                "/sbin/generate-certificates",
                "/etc/init.d/hostd restart",
                "/etc/init.d/vpxa restart",
            ],
            opts=ResourceOptions(depends_on=[step_2]),
        )
        step_4 = RemoteExec(
            "configure-" + node_name + "-step-4",
            host_id=instance.id,
            conn=conn_helper_args,
            commands=[command_cleanup],
            opts=ResourceOptions(depends_on=[step_3]),
        )
