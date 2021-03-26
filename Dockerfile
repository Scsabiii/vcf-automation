FROM pulumi/pulumi-python:2.18.1
LABEL source_repository="https://github.com/sapcc/avocado-automation"

ARG workdir=/pulumi/avocado
WORKDIR ${workdir} 

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
		"pulumi==2.18.1" \
		"pulumi-openstack==2.17.0" \
		"paramiko>=2.7.1" \
		"typing_extensions>=3.7.4"

RUN apt update && \
	apt install -yq --no-install-recommends \
	vim-tiny && \
	apt clean && \
	rm -rf /var/lib/apt/lists/*

# COPY test/etc ${workdir}/etc
COPY test/projects/management ${workdir}/projects/management
COPY bin/automation /pulumi/bin/automation

