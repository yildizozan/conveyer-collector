FROM golang:1.14.0 AS builder
WORKDIR /go/src/github.com/yildizozan/conveyer-service-collector/

# Dependencies
#RUN go get -d -v golang.org/x/net/html
RUN go get -d -v github.com/streadway/amqp

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./...

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/yildizozan/conveyer-service-collector/collector .
CMD ["./collector"]
