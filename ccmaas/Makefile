app := ccmaas
GOFILES := main.go $(wildcard **/*.go)

.PHONY: build
build: ../bin/${app}_linux_amd64 ../bin/${app}_darwin_amd64

../bin/${app}_linux_amd64: $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -o $@ -mod vendor 

../bin/${app}_darwin_amd64: $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -o $@ -mod vendor
