FROM golang:latest as builder

WORKDIR /app
COPY . .
ENV GO111MODULE=on

RUN go get -d -v ./...
RUN go install -v ./...

RUN go install github.com/stellar/go/tools/...
RUN go install github.com/stellar/go/services/...
RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o server_stellar_gallery main.go

#FROM scratch
FROM centos:latest
COPY --from=builder /app/server_stellar_gallery /

# Copy CA certificates to prevent x509: certificate signed by unknown authority errors
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 5000
ENTRYPOINT ["/server_stellar_gallery"]
CMD ["--account_public_key"] 