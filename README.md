# VCF Automation

Automation tool for deploying VMware Cloud Foundation (VCF) on SAP Converged
Cloud (CCloud). The tool can deploy and extend the infrastructures needed by a
VCF stack on CCloud, including networks, storages and vSphere servers.

## Setup

<!-- The project is provisioned in several separate stacks: `management` -->
<!-- and `workload`. Some resources are shared by management domain and workload -->
<!-- domain, e.g., management network and deployment network. Therefore they are -->
<!-- provisioned in the shared stack. -->

<!-- The configuration file are stored in the directory `$workdir/etc/`. The file -->
<!-- names follow the convention `{project-name}-{stack-name}.yaml`. -->

```
projectType: a string "vcf/management" or "vcf/workload"
stack: a unique name
props:
  openstack:
    region: ...
    domain: ...
    tenant: ...
  stack:
    externalNetwork:
      name: FloatingIP-external-...
      id: ...
```

- shared stack:

```
project: management
stack: shared
props:
  openstack:
    region: qa-de-1
    domain: monsoon3
    tenant: vcfonccloud-management-domain
  stack:
    shared:
      externalNetwork:
        name: FloatingIP-external-monsoon3-03
        id: 430991b3-da0d-41cb-ac54-d1d532841725
      managementNetwork:
        networkName: DapnIP-external-monsoon3-management-vcf01
        subnetName: DapnIP-sap-monsoon3-management-vcf01-01
        subnetGateway: 10.237.209.1
        subnetMask: 255.255.255.128
        esxiInterface: vmk20
        vlanID: 1007
      deploymentNetwork:
        networkName: mgmtdomain-deployment-network
        subnetName: mgmtdomain-deployment-subnet
        cidr: 10.180.6.0/24
        gatewayIP: 10.180.6.1
      publicRouter: mgmtdomain-public-router
      helperVM:
        imageName: ubuntu-20.04-amd64-vmware
        flavorID: 20
        ip: 10.237.209.100
      dnsZoneName: vcf01.qa-de-1.cloud.sap.
      reverseDnsZoneName: 209.237.10.in-addr.arpa.
```

- management domain: the management domain also needs the information of the
  shared resources, and thus the `dependsOn` field is set. The properties of the
  shared stack is merged into the management properties.

```
project: management
stack: dev
dependsOn:
  - management-shared.yaml
props:
  openstack:
    region: qa-de-1
    domain: monsoon3
    tenant: vcfonccloud-management-domain
  stack:
    esxiServerImage: vsphere-7.0.1-amd64-baremetal
    esxiServerFlavorID: 20
    esxiNodes:
      - name: node003-bb096
        id: 050c9481-c27a-4ac3-895a-3e2eaf345409
        ip: 10.237.209.20
    privateNetworks:
      - networkName: vmotion
        cidr: 10.180.0.0/24
        vlanID: 1000
        esxiInterface: vmk10
      - networkName: edgetep
        cidr: 10.180.1.0/24
        vlanID: 1001
        esxiInterface: vmk11
      - networkName: hosttep
        cidr: 10.180.2.0/24
        vlanID: 1002
        esxiInterface: vmk12
      - networkName: nfs
        cidr: 10.180.3.0/24
        vlanID: 1003
        esxiInterface: vmk13
      - networkName: vsan
        cidr: 10.180.4.0/24
        vlanID: 1004
        esxiInterface: vmk14
      - networkName: vsanwitness
        cidr: 10.180.5.0/24
        vlanID: 1005
        esxiInterface: vmk15
```

- workload domain:

```
project: management
stack: workload
dependsOn:
  - management-shared.yaml
props:
  openstack:
    region: qa-de-1
    domain: monsoon3
    tenant: vcfonccloud-management-domain
  stack:
    esxiServerImage: vsphere-7.0.1-amd64-baremetal
    esxiServerFlavorID: 20
    esxiNodes:
      - name: node002-bb096
        id: a5e68d40-4b1f-49aa-9a54-82283015b9b0
        ip: 10.237.209.14
      - name: node003-bb096
        id: 050c9481-c27a-4ac3-895a-3e2eaf345409
        ip: 10.237.209.15
      - name: node020-ap001
        id: 53bd7630-005c-4316-8617-d44c42ee4fe8
        ip: 10.237.209.16
      - name: node020-ap002
        id: 43d5cec9-4616-4229-9bb0-4cf866b93dbc
        ip: 10.237.209.17
    privateNetworks:
      - networkName: vmotion
        cidr: 10.180.10.0/24
        vlanID: 1000
        esxiInterface: vmk10
      - networkName: edgetep
        cidr: 10.180.11.0/24
        vlanID: 1001
        esxiInterface: vmk11
      - networkName: hosttep
        cidr: 10.180.12.0/24
        vlanID: 1002
        esxiInterface: vmk12
      - networkName: nfs
        cidr: 10.180.13.0/24
        vlanID: 1003
        esxiInterface: vmk13
      - networkName: vsan
        cidr: 10.180.14.0/24
        vlanID: 1004
        esxiInterface: vmk14
      - networkName: vsanwitness
        cidr: 10.180.15.0/24
        vlanID: 1005
        esxiInterface: vmk15
```

## API

- Visit `localhost:8080/stacks` to list stacks. The stack is in either `running`
  mode or `stopped` mode.

```
HTTP/1.1 200 OK
Date: Tue, 25 May 2021 17:55:13 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8

management-shared.yaml: running
management-dev.yaml: stopped
```

- The stack can be started or paused via
  `localhost:8080/{project-name}/{stack-name}/start` or
  `localhost:8080/{project-name}/{stack-name}/stop`

```
curl -X POST -i http://localhost:8080/management/dev/start
HTTP/1.1 200 OK
Date: Tue, 25 May 2021 18:00:05 GMT
Content-Length: 13
Content-Type: text/plain; charset=utf-8

stack started%
```

```
curl -X POST -i http://localhost:8080/management/dev/stop
HTTP/1.1 200 OK
Date: Tue, 25 May 2021 18:00:21 GMT
Content-Length: 13
Content-Type: text/plain; charset=utf-8

stack stopped%
```

- Show status of a stack via
  `localhost:8080/{project-name}/{stack-name}/state`

```
curl -X GET -i http://localhost:8080/management/dev/state
HTTP/1.1 200 OK
Date: Tue, 25 May 2021 18:01:34 GMT
Content-Length: 1498
Content-Type: text/plain; charset=utf-8

failed to run update: exit status 255
code: 255
stdout: Updating (dev):

    pulumi:pulumi:Stack management-dev running
    pulumi:providers:openstack ccadmin_master
    pulumi:providers:openstack cloud_admin
    openstack:compute:Keypair rsa-keypair
    pulumi:pulumi:Stack management-dev running error: Missing required configuration variable 'management:externalNetwork'
    pulumi:pulumi:Stack management-dev running error: an unhandled error occurred: Program exited with non-zero exit code: 1
    pulumi:pulumi:Stack management-dev **failed** 2 errors

Diagnostics:
  pulumi:pulumi:Stack (management-dev):
    error: Missing required configuration variable 'management:externalNetwork'
        please set a value using the command `pulumi config set management:externalNetwork <value>`
    error: an unhandled error occurred: Program exited with non-zero exit code: 1

Outputs:
  - DapnIPNetworkID         : "b30bea36-8746-4487-8589-2ea9e3ae9ea0"
  - DapnIPSubnetID          : "2f7f0b18-8d1b-4bc4-a1ca-79000b2dd4d8"
  - EsxiHostIP-node016-ap001: "10.237.209.10"
  - EsxiHostIP-node017-ap001: "10.237.209.11"
  - EsxiHostIP-node018-ap001: "10.237.209.12"
  - EsxiHostIP-node019-ap001: "10.237.209.13"
  - helperVMIP              : "10.237.209.100"

Resources:
    4 unchanged

Duration: 1s
```
