FROM golang:1.19 as builder

COPY . /src
WORKDIR /src

RUN make clean && make BUILD_IN_CONTAINER=false loki-gateway


# executable image
FROM alpine:3

RUN apk add --no-cache ca-certificates

COPY --from=builder /src/cmd/loki-gateway/loki-gateway /usr/bin/loki-gateway

ENTRYPOINT ["/usr/bin/loki-gateway"]
