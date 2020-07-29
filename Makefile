.PHONY: build clean deploy

.PHONY: compile
PROTOS = position measurement
compile:
	for proto in $(PROTOS); do protoc -I pkg/proto/$$proto/ pkg/proto/$$proto/$$proto.proto --go_out=plugins=grpc:pkg/proto/$$proto; done

build:
	go build -race -ldflags="-s -w" -o bin/collector cmd/collector/main.go

clean:
	rm -rf ./bin
