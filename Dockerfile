ARG pulumi_version=2.18.1
FROM pulumi/pulumi-python:${pulumi_version}
LABEL source_repository="https://github.com/sapcc/avocado-automation"

ARG workdir=/pulumi/avocado
WORKDIR ${workdir} 

RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
		"pulumi>=2.0.0,<3.0.0" \
		"pulumi-openstack>=2.0.0,<3.0.0" \
		"paramiko>=2.7.1" \
		"typing_extensions>=3.7.4"

COPY test/etc ${workdir}/etc
COPY test/projects/management ${workdir}/projects/management
COPY bin/automation /pulumi/bin/automation

