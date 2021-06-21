# build go binary
FROM keppel.eu-de-1.cloud.sap/ccloud-dockerhub-mirror/library/golang:1.16-alpine AS build
RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
WORKDIR $GOPATH/src/automation 
COPY . .
RUN go build -o automation .

# pulumi python
FROM keppel.eu-de-1.cloud.sap/ccloud-dockerhub-mirror/pulumi/pulumi-python:3.2.0
LABEL source_repository="https://github.com/sapcc/vcf-automation"

ARG workdir=/pulumi/automation
WORKDIR ${workdir} 

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
		"pulumi>=3.0.0<4.0.0" \
		"pulumi-openstack>=3.0.0<4.0.0" \
		"paramiko>=2.7.1" \
		"typing_extensions>=3.7.4" \
		"jinja2"

RUN apt update && \
	apt install -yq --no-install-recommends \
	iputils-ping \
	iputils-arping \
	iputils-tracepath \
	traceroute \
	jq \
	vim-tiny && \
	apt clean && \
	rm -rf /var/lib/apt/lists/*

# COPY test/etc ${workdir}/etc
COPY projects/vcf ${workdir}/projects/vcf
COPY --from=build /go/src/automation/automation /pulumi/bin/automation
COPY static ${workdir}/static
COPY templates ${workdir}/templates

