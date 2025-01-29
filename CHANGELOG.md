# Changelog

# [v1.8.0] - 2025-01-29

### New Minor Release
- Added support for Cyberark-sourced secrets

# [v1.7.1] - 2023-05-30

### New bugfix release

- Fixes bug in template helper `expandFullChain` introduced in https://github.com/covermymeds/azure-key-vault-agent/pull/86 where `expandFullChain` would throw a fatal error if secrets did not have a ContentType set


# [v1.7.0] - 2023-05-19

### New minor release

- Update Go version to 1.17.x
- Update goreleaser github action to v3 and limit build OS types to darwin and linux
- Add template helper `expandFullChain` - which returns a map of all secret objects from a keyvault, including pem and key files when a secret is a certificate
