FROM golang:alpine as builder
WORKDIR /go/src/app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o /go/bin/azure-key-vault-agent

# Start a new stage and copy the built binary
FROM alpine
COPY --from=builder /go/bin/azure-key-vault-agent /azure-key-vault-agent
ENTRYPOINT ["/azure-key-vault-agent"]
