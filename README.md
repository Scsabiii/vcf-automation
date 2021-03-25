# Automation for Hypervisor Deployment (Project Avocado)

This repository contains the codes for provisioning physical servers in SAP Converged Cloud.

## Usage

- Create management domain

```
  {
    project: management,
    stack: "a string that can be anything",
    props: {
      "openstack": {
        "domain": "openstack domain name",
        "tenant": "openstack project name",
        "region": "openstack region",
      },
      "stack": {

      }
    }

  }
```
