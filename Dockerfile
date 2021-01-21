# pulumi base
############################################################
ARG pulumi_version=latest
FROM pulumi/pulumi-base:${pulumi_version} as pulumi

# The runtime container
############################################################

FROM golang:1.15-alpine
LABEL source_repository="https://github.com/sapcc/ccmaas-operator"

WORKDIR /pulumi/ccmaas/src

ENV PATH "/pulumi/bin:${PATH}"

COPY --from=pulumi /pulumi/bin/pulumi /pulumi/bin/pulumi
COPY --from=pulumi /pulumi/bin/pulumi-language-go /pulumi/bin/pulumi-language-go
COPY --from=pulumi /pulumi/bin/pulumi-analyzer-policy /pulumi/bin/pulumi-analyzer-policy

RUN apk add --no-cache git libc6-compat ca-certificates

COPY ccmaas /pulumi/ccmaas/src

COPY etc /pulumi/ccmaas/etc
COPY projects /pulumi/ccmaas/projects
# RUN cd src/ && go build -mod vendor -o /pulumi/bin/ccmaas

ENTRYPOINT [ "/pulumi/bin/ccmaas"]
