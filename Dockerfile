FROM pulumi/pulumi-python:3.2.0
LABEL source_repository="https://github.com/sapcc/avocado-automation"

ARG workdir=/pulumi/avocado
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
	vim-tiny && \
	apt clean && \
	rm -rf /var/lib/apt/lists/*

# COPY test/etc ${workdir}/etc
COPY test/projects/management ${workdir}/projects/management
COPY bin/automation /pulumi/bin/automation

