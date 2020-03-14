FROM golang:1.14.0 AS builder
WORKDIR /go/src/github.com/yildizozan/conveyer-service-collector/

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./...

FROM scratch
COPY --from=builder /collector ./
CMD ["./collector"]
