FROM golang:alpine as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM alpine
COPY --from=builder /go/bin/azure-key-vault-agent /azure-key-vault-agent
CMD ["/azure-key-vault-agent", "-config", "/akva.yaml"]
