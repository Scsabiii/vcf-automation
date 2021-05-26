# VCF deployment automation

This repository contains codes for provisioning VSphere on physical servers in SAP Converged Cloud.

## Setup

### shared stack

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

### management domain

```
project: management
stack: dev2
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

### workload domain

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
