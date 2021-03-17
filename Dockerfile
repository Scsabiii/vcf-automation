# pulumi base
############################################################
ARG pulumi_version=2.18.1
FROM pulumi/pulumi-base:${pulumi_version} as pulumi

# The runtime container
############################################################
FROM golang:1.15-alpine
LABEL source_repository="https://github.com/sapcc/avocado-automation"

ARG workdir=/pulumi/avocado
WORKDIR ${workdir} 

ENV PATH "/pulumi/bin:${PATH}"

COPY --from=pulumi /pulumi/bin/pulumi /pulumi/bin/pulumi
COPY --from=pulumi /pulumi/bin/pulumi-language-go /pulumi/bin/pulumi-language-go
COPY --from=pulumi /pulumi/bin/pulumi-analyzer-policy /pulumi/bin/pulumi-analyzer-policy

RUN apk add --no-cache git libc6-compat ca-certificates

COPY test/etc ${workdir}/etc
COPY projects ${workdir}/projects
COPY entrypoint.sh /

COPY bin/automation_linux_amd64 /pulumi/bin/automation

ENTRYPOINT [ "/pulumi/bin/automation" ]
