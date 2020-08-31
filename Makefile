#PROTOC_GEN_GO_GIT_TAG := "1.2.0"
PROTOC_GEN_GO_GIT_TAG := "1.4.2"

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

.PHONY: clean
clean:
	rm -rf ./tools

protoc-%:
	$(eval VERSION := $*)
	rm -rf tools/$(VERSION) || true
	mkdir -p tools/$(VERSION)
	wget https://github.com/protocolbuffers/protobuf/releases/download/v$(VERSION)/protoc-$(VERSION)-osx-x86_64.zip -qO tools/protoc.zip
	unzip -q tools/protoc.zip -d tools/$(VERSION)

.PHONY: protos
protos: guard-PROTOC_VERSION guard-PROTOC_GEN_GO
	echo "> BUILD protoc@v$(PROTOC_VERSION) protoc-gen-go@v$(PROTOC_GEN_GO)"
	$(MAKE) protoc-$(PROTOC_VERSION)
	mkdir -p protos/sample
	go get -u github.com/golang/protobuf/protoc-gen-go@v$(PROTOC_GEN_GO)
	tools/$(PROTOC_VERSION)/bin/protoc -I ./protos --go_out=plugins=grpc:./protos/sample sample.proto
	go mod vendor  &>/dev/null

.PHONY: test
test:
	go test -race -test.parallel=1 main_test.go -test.v

.PHONY: test-all
test-all: clean
	PROTOC_VERSION=3.5.1 PROTOC_GEN_GO=1.2.0 $(MAKE) protos test
	PROTOC_VERSION=3.5.1 PROTOC_GEN_GO=1.4.2 $(MAKE) protos test
	PROTOC_VERSION=3.12.4 PROTOC_GEN_GO=1.4.2 $(MAKE) protos test