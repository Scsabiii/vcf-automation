# syntax = docker/dockerfile:experimental
# Interim container so we can copy pulumi binaries
# Must be defined first
ARG PULUMI_VERSION=latest
ARG PULUMI_IMAGE=keppel.eu-de-1.cloud.sap/ccloud/pulumi-go
FROM ${PULUMI_IMAGE}:${PULUMI_VERSION}

LABEL source_repository="https://github.com/sapcc/esxi-operator"

COPY projects /pulumi/projects
