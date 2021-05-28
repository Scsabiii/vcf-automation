FROM pulumi/pulumi-python:3.2.0
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
COPY test/projects/vcf ${workdir}/projects/vcf
COPY bin/automation /pulumi/bin/automation
COPY static ${workdir}/static
COPY templates ${workdir}/templates

