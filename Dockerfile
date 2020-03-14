FROM golang:1.14.0 AS builder

LABEL maintainer="developer@yildizozan.com"

WORKDIR /go/src/github.com/yildizozan/conveyer-service-collector/

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GORCH amd64

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . ./

RUN go install -a -installsuffix nocgo ./...

FROM scratch
COPY --from=builder /go/bin/collector ./
CMD ["./collector"]
