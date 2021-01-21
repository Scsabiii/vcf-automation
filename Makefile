
.PHONY: build
build: ccmaas/Makefile
	$(MAKE) -C ccmaas build

.PHONY: docker
docker:
	DOCKER_BUILDKIT=1 docker build . -t keppel.eu-de-1.cloud.sap/ccloud/ccmaas-operator:latest
