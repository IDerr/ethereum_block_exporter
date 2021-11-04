FROM golang:1.16 as builder

WORKDIR /ethereum_exporter
COPY . .

ARG VERSION=v0.0.1
RUN CGO_ENABLED=0 \
    go build -ldflags "-s -w -X main.version=$VERSION" github.com/iderr/ethereum_block_exporter/cmd/ethereum_block_exporter

FROM scratch

ENTRYPOINT ["/ethereum_block_exporter"]
EXPOSE 9368

COPY --from=builder /ethereum_exporter/ethereum_block_exporter /ethereum_block_exporter
