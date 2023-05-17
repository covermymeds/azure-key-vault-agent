# Changelog

# [v1.7.0] - 2023-05-19

### New minor release

- Update Go version to 1.17.x
- Update goreleaser github action to v3 and limit build OS types to darwin and linux
- Add template helper `expandFullChain` - which returns a map of all secret objects from a keyvault, including pem and key files when a secret is a certificate
