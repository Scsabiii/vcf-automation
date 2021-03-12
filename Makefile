app=ccmaas
GOFILES := main.go $(wildcard cmd/**/*.go pkg/**/*.go)

.PHONY: build
build: bin/${app} bin/${app}_linux_amd64 bin/${app}_darwin_amd64

.PHONY: docker
docker:
	DOCKER_BUILDKIT=1 docker build . -t keppel.eu-de-1.cloud.sap/ccloud/ccmaas-operator:latest

bin/${app}: bin/${app}_darwin_amd64
	@cp $< $@

bin/${app}_linux_amd64: $(GOFILES)
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o $@ -mod vendor

bin/${app}_darwin_amd64: $(GOFILES)
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o $@ -mod vendor

