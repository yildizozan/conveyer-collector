FROM golang:alpine AS builder
WORKDIR /go/src/github.com/yildizozan/conveyer-service-collector/

# Dependencies
RUN go get -d -v golang.org/x/net/html

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/yildizozan/conveyer-service-collector/app .
CMD ["./app"]