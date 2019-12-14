# Usage

Create `.env` and populate it as follows:

```bash
AZURE_TENANT_ID=<tenant id>
AZURE_SUBSCRIPTION_ID=<subscription id>
AZURE_CLIENT_ID=<SPN name including http>
AZURE_CLIENT_SECRET=<SPN password>
```

Then run:

`go run main/main.go`

# Notes for later

To track tool deps explicitly:

https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
