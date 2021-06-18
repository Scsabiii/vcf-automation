{
  "skipEsxThumbprintValidation": true,
  "managementPoolName": "sddc-a-np01",
  "sddcManagerSpec": {
    "secondUserCredentials": {
      "username": "vcf",
      "password": "{{ vmware_password }}"
    },
    "ipAddress": "{{ sddc_manager.ip }}",
    "netmask": "255.255.255.128",
    "hostname": "{{ sddc_manager.hostname }}",
    "rootUserCredentials": {
      "username": "root",
      "password": "{{ vmware_password }}"
    },
    "restApiCredentials": {
      "username": "admin",
      "password": "{{ vmware_password }}"
    },
    "localUserPassword": "{{ vmware_password }}",
    "vcenterId": "vcenter-1"
  },
  "sddcId": "{{ sddc_manager.id }}",
  "esxLicense": "{{ sddc_manager.esx_license }}",
  "taskName": "workflowconfig/workflowspec-ems.json",
  "ceipEnabled": false,
  "ntpServers": ["147.204.9.202", "147.204.9.203"],
  "dnsSpec": {
    "nameserver": "147.204.9.200",
    "secondaryNameserver": "147.204.9.201",
    "domain": "{{ sddc_manager.domain }}",
    "subdomain": "{{ sddc_manager.domain }}"
  },
  "networkSpecs": [
    {
      "networkType": "MANAGEMENT",
      "subnet": "{{ management_network.subnet_cidr }}",
      "gateway": "{{ management_network.subnet_gateway }}",
      "vlanId": "1007",
      "mtu": "8950",
      "portGroupKey": "vc-a-mgmt-vds01-pg-mgmt",
      "standbyUplinks": ["uplink2"],
      "activeUplinks": ["uplink1"]
    },
    {
      "networkType": "VMOTION",
      "subnet": "10.180.0.0/24",
      "gateway": "10.180.0.1",
      "vlanId": "1000",
      "mtu": "8950",
      "portGroupKey": "vc-a-mgmt-vds01-pg-vmotion",
      "association": "vc-a-mgmt-dc01",
      "includeIpAddressRanges": [
        { "endIpAddress": "10.180.0.132", "startIpAddress": "10.180.0.101" }
      ],
      "standbyUplinks": ["uplink2"],
      "activeUplinks": ["uplink1"]
    },
    {
      "networkType": "VSAN",
      "subnet": "10.180.4.0/24",
      "gateway": "10.180.4.1",
      "vlanId": "1004",
      "mtu": "8950",
      "portGroupKey": "vc-a-mgmt-vds01-pg-vsan",
      "includeIpAddressRanges": [
        { "endIpAddress": "10.180.4.132", "startIpAddress": "10.180.4.101" }
      ],
      "standbyUplinks": ["uplink2"],
      "activeUplinks": ["uplink1"]
    }
  ],
  "nsxtSpec": {
    "nsxtManagerSize": "medium",
    "nsxtManagers": [
      {%- for m in nsxt_managers %}
      {
          "hostname": "{{ m.hostname }}",
          "ip": "{{ m.ip }}"
      {%- if loop.last %}
      }
      {%- else %}
      },
      {%- endif %}
      {%- endfor %}
    ],
    "rootNsxtManagerPassword": "{{ vmware_password }}",
    "nsxtAdminPassword": "{{ vmware_password }}",
    "nsxtAuditPassword": "{{ vmware_password }}",
    "rootLoginEnabledForNsxtManager": "true",
    "sshEnabledForNsxtManager": "true",
    "overLayTransportZone": {
      "zoneName": "{{ region }}-tz-overlay01",
      "networkName": "netName-overlay"
    },
    "vlanTransportZone": {
      "zoneName": "{{ region }}-tz-vlan01",
      "networkName": "netName-vlan"
    },
    "vip": "{{ nsxt.ip }}",
    "vipFqdn": "{{ nsxt.hostname }}",
    "nsxtLicense": "{{ nsxt.license }}",
    "transportVlanId": 1002,
    "ipAddressPoolSpec": {
      "name": "vc-a-mgmt-vds01-pg-tep01",
      "description": "ESXi Host Overlay TEP IP Pool",
      "subnets": [
        {
          "ipAddressPoolRanges": [
            {
              "start": "10.180.2.101",
              "end": "10.180.2.164"
            }
          ],
          "cidr": "10.180.2.0/24",
          "gateway": "10.180.2.1"
        }
      ]
    }
  },
  "vsanSpec": {
    "vsanName": "vsan-1",
    "vsanDedup": "false",
    "datastoreName": "vc-{{ region }}-cl01-ds-vsan01"
  },
  "dvSwitchVersion": "7.0.0",
  "dvsSpecs": [
    {
      "dvsName": "vc-a-cl01-vds01",
      "vcenterId": "vcenter-1",
      "vmnics": ["vmnic0", "vmnic2"],
      "mtu": 8950,
      "networks": ["MANAGEMENT", "VMOTION", "VSAN"],
      "niocSpecs": [
        {
          "trafficType": "VSAN",
          "value": "HIGH"
        },
        {
          "trafficType": "VMOTION",
          "value": "LOW"
        },
        {
          "trafficType": "VDP",
          "value": "LOW"
        },
        {
          "trafficType": "VIRTUALMACHINE",
          "value": "HIGH"
        },
        {
          "trafficType": "MANAGEMENT",
          "value": "NORMAL"
        },
        {
          "trafficType": "NFS",
          "value": "LOW"
        },
        {
          "trafficType": "HBR",
          "value": "LOW"
        },
        {
          "trafficType": "FAULTTOLERANCE",
          "value": "LOW"
        },
        {
          "trafficType": "ISCSI",
          "value": "LOW"
        }
      ],
      "isUsedByNsxt": true
    }
  ],
  "clusterSpec": {
    "clusterName": "vc-a-mgmt-cl01",
    "vcenterName": "vcenter-1",
    "clusterEvcMode": "",
    "vmFolders": {
      "MANAGEMENT": "qa-vcf01-fd-mgmt",
      "NETWORKING": "qa-vcf01-fd-nsx",
      "EDGENODES": "qa-vcf01-fd-edge"
    }
  },
  "pscSpecs": [
    {
      "pscId": "psc-1",
      "vcenterId": "vcenter-1",
      "adminUserSsoPassword": "{{ vmware_password }}",
      "pscSsoSpec": {
        "ssoDomain": "vsphere.local"
      }
    }
  ],
  "vcenterSpec": {
    "vcenterIp": "{{ vcenter.ip }}",
    "vcenterHostname": "{{ vcenter.hostname }}",
    "vcenterId": "vcenter-1",
    "vmSize": "small",
    "storageSize": "",
    "rootVcenterPassword": "{{ vmware_password }}"
  },
  "hostSpecs": [
    {%- for server in esxi_servers %}
    {
      "association": "vc-a-mgmt-dc01",
      "ipAddressPrivate": {
        "ipAddress": "{{ server.ip }}",
        "cidr": "{{ management_network.subnet_cidr }}",
        "gateway": "{{ management_network.subnet_gateway }}"
      },
      "hostname": "esxi-{{ server.name }}",
      "credentials": {
        "username": "root",
        "password": "{{ vmware_password }}"
      },
      "vSwitch": "vSwitch0",
      "serverId": "host-{{ loop.index }}"
    {%- if loop.last %}
    }
    {%- else %}
    },
    {%- endif %}
    {%- endfor %}
  ],
  "excludedComponents": ["NSX-V", "AVN", "EBGP"]
}
