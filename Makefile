.PHONY: build clean deploy

build:
	dep ensure -v
#	env go build -ldflags="-s -w" -o bin/collector    		cmd/collector/main.go
	env go build -o bin/collector    		cmd/collector/main.go
	env go build -o bin/sensor    			cmd/sensor/main.go
	env go build -o bin/server    			cmd/server/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose