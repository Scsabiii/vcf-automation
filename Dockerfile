# pulumi base
############################################################
ARG pulumi_version=2.18.1
FROM pulumi/pulumi-base:${pulumi_version} as pulumi

# The runtime container
############################################################

FROM golang:1.15-alpine
LABEL source_repository="https://github.com/sapcc/ccmaas-operator"

WORKDIR /pulumi/ccmaas

ENV PATH "/pulumi/bin:${PATH}"

COPY --from=pulumi /pulumi/bin/pulumi /pulumi/bin/pulumi
COPY --from=pulumi /pulumi/bin/pulumi-language-go /pulumi/bin/pulumi-language-go
COPY --from=pulumi /pulumi/bin/pulumi-analyzer-policy /pulumi/bin/pulumi-analyzer-policy

RUN apk add --no-cache git libc6-compat ca-certificates

COPY etc /pulumi/ccmaas/etc
COPY projects /pulumi/ccmaas/projects
COPY entrypoint.sh /

COPY bin/ccmaas_linux_amd64 /pulumi/bin/ccmaas
# RUN cd src/ && go build -mod vendor -o /pulumi/bin/ccmaas

ENTRYPOINT [ "/pulumi/bin/ccmaas"]
