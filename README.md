# VCF Automation

Automation service for deploying VMware Cloud Foundation (VCF) on Openstack
powered SAP Converged Cloud (CCloud). The resources needed for VCF stack is
defined by yaml configuration. This automation tool provisions the resources in
CCloud accordingly and keep them in sync with configuration.

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

## Commands

- `automation server` starts automation server. It spawns a controller loop for
  each configuration and provision the stack.
- `automation configure` allows generate pulumi's config file in project
  directory on cli manually.

## API

- Endpoint `/vcf` returns json object which gives an overview of all running
  stacks. For example,

```
[
  {
    "name": "vcf-vcf-01-management",
    "config_file": "/pulumi/automation/etc/vcf-01-management.yaml",
    "status": "running",
    "has_error": true,
    "links": [
      {
        "name": "cloud-builder",
        "url": "http://localhost:8080/vcf/vcf-01-management/cloud-builder.json",
        "description": "payload for cloud builder"
      },
      {
        "name": "state",
        "url": "http://localhost:8080/vcf/vcf-01-management/state",
        "description": "resources deployed by automation"
      },
      {
        "name": "error",
        "url": "http://localhost:8080/vcf/vcf-01-management/error"
      },
      {
        "name": "start",
        "url": "http://localhost:8080/vcf/vcf-01-management/start",
        "description": "restart automation controller loop"
      },
      {
        "name": "stop",
        "url": "http://localhost:8080/vcf/vcf-01-management/stop",
        "description": "pause automation controller"
      },
      {
        "name": "reload",
        "url": "http://localhost:8080/vcf/vcf-01-management/reload",
        "description": "force controller to reload configuration"
      }
    ]
  }
]
```

- Endpoint `/vcf/reload` reloads all the configuration files in the configure
  directory and update the running controllers or spawns a new controller.

- Endpoint `/vcf/{stack-name}/[state,error,start,stop,reload]` shows stack
  details or gives stack specific control.
