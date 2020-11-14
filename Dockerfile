FROM golang:alpine as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=builder /go/bin/azure-key-vault-agent /azure-key-vault-agent
ENTRYPOINT ["/azure-key-vault-agent"]
CMD ["-config", "/akva.yaml"]
