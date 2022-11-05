FROM golang:1.19 as builder

COPY . /src
WORKDIR /src

RUN make clean && make BUILD_IN_CONTAINER=false mimir-gateway


# executable image
FROM alpine:3

RUN apk add --no-cache ca-certificates

COPY --from=builder /src/cmd/mimir-gateway/mimir-gateway /usr/bin/mimir-gateway

ENTRYPOINT ["/usr/bin/mimir-gateway"]
