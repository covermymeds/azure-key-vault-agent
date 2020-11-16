FROM alpine
RUN apk add --no-cache libc6-compat
ADD https://github.com/covermymeds/azure-key-vault-agent/releases/download/v1.4.0/azure-key-vault-agent_1.4.0_linux_amd64.tar.gz /opt/akva/
RUN tar -xvf /opt/akva/azure-key-vault-agent_1.4.0_linux_amd64.tar.gz -C /opt/akva/
CMD ["/opt/akva/azure-key-vault-agent", "-config", "/opt/akva/akva.yaml"]
