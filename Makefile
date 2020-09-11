all: compile build
.PHONY: all

.PHONY: compile
PROTOS = conveyor
compile:
	for proto in $(PROTOS); do protoc -I pkg/proto/$$proto/ pkg/proto/$$proto/$$proto.proto --go_out=plugins=grpc:pkg/proto/$$proto; done

.PHONY: build
build:
	go build -race -ldflags="-s -w" -o bin/collector cmd/collector/main.go

.PHONY: clean
clean:
	rm -rf ./bin
