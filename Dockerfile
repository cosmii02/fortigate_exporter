# Build using the minimum supported Golang version (match go.mod)
FROM golang:1.26 AS builder

WORKDIR /build

COPY . .
RUN go mod download
RUN make build

FROM scratch
WORKDIR /opt/fortigate_exporter

COPY --from=builder /build/target/fortigate-exporter .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt .
ENV SSL_CERT_DIR=/opt/fortigate_exporter

EXPOSE 9710
ENTRYPOINT ["./fortigate-exporter"]
CMD ["-auth-file", "/config/fortigate-key.yaml"]
