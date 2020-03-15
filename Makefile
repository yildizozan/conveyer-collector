.PHONY: build clean deploy

build:
#	env go build -ldflags="-s -w" -o bin/collector    		cmd/collector/main.go
	go build -o bin/collector	cmd/collector/main.go
	go build -o bin/sensor    	cmd/sensor/main.go
	go build -o bin/server    	cmd/server/main.go

clean:
	rm -rf ./bin