app=automation
# GOFILES := $(wildcard *.go cmd/*.go pkg/**/*.go pkg/**/**/*.go)
GOFILES := $(shell find . -name "*.go")

.PHONY: build
build: bin/${app} bin/${app}_linux_amd64 bin/${app}_darwin_amd64

.PHONY: docker
docker: bin/${app}_linux_amd64
	DOCKER_BUILDKIT=1 docker build . -t keppel.eu-de-1.cloud.sap/ccloud/vcf-${app}:latest
	docker push keppel.eu-de-1.cloud.sap/ccloud/vcf-${app}:latest

bin/${app}: bin/${app}_linux_amd64
	@cp $< $@

bin/${app}_linux_amd64: $(GOFILES)
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $@ -mod vendor
	@chmod +x $@

bin/${app}_darwin_amd64: $(GOFILES)
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o $@ -mod vendor
	@chmod +x $@

.PHONY: info
info: ; $(info $$GOFILES is [${GOFILES}])

