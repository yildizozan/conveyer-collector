FROM golang:1.14.0 AS builder

LABEL maintainer="developer@yildizozan.com"

WORKDIR /go/src/conveyer-service-collector/

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GORCH amd64

COPY . ./

RUN go install -a -installsuffix nocgo ./...

FROM scratch
COPY --from=builder /go/bin/collector ./
CMD ["./collector"]
