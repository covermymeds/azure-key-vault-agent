# Usage

`akva --config=akva.yaml`

# Authentication

You can specify named credentials in the config file under the top-level key `credentials`.

Additionally, you can specify the following environment variables (or specify these key=value pairs in a file called `.env`):

```bash
AZURE_TENANT_ID=<tenant id>
AZURE_CLIENT_ID=<SPN name including http>
AZURE_CLIENT_SECRET=<SPN password>
```

They will be loaded as the credential name `default`.

# Config

Create a yaml file which holds configuration for one or more credentials and one or more workers.

Each worker can pull one or more resources and use those resources to write to one or more file sinks.

A simple example is given below:

```yaml
credentials:
  -
    name: default
    tenantID: test-id
    clientID: https://test-client-id
    clientSecret: test-secret

workers:
  -
    resources:
      - kind: secret
        name: password
        vaultBaseURL: https://test-kv.vault.azure.net/
        # No credential specified, so "default" will be used
    sinks:
      - path: ./password
        template: "{{ .Secrets.password.Value }}"
```

## Credentials
The `credentials` section is a list of one or more named credentials used for fetching resources. Each
credential has a `tenantID`, `clientID`, `clientSecret`.

The ENV vars (or .env file) will be injected
as a credential with the name `default` if you don't override `default` within your config file.

## Resources

The `resources` section is a list of one or more resources to fetch. Each resource has a `kind`, `vaultBaseURL`,
and optional `credential` field.

Valid kinds are: `cert`, `secret`, `all-secrets`, and `key`.

Note: The `all-secrets` resource cannot be used with any `secret` resources.

Unless a resource has a `kind` of `all-secrets`, there is also a required `name` field for the resource.

If you don't specify `credential`, a credential with the name `default` will be used (you can either
specify the `default` credential in the `credentials` array, or as ENV vars / .env file)

## Sinks

The `sinks` section is a list of one or more files to write to. Each sink has a `path` and either `template` (inline template) or `templatePath` (path to template on the filesystem). The template syntax is golang's [text/template](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) library (with [sprig](https://github.com/Masterminds/sprig) helpers).

`sinks` also support configuring file ownership and permission bits via the `owner`, `group`, and `mode` settings.
- `owner` and `group` are the names of the respective entity and must both be present.  If omitted the executing user and group will be applied.
- `mode` accepts file modes in either 3 or 4 digit notation `777`, `1644`, `0600` are all valid examples.  If omitted a default of `0644` will be used.


Each template has access to all of the resources specified in the `resources` section above, separated by kind and resource name. The fields available to you for any given resource can be found by looking at the corresponding source structs:

* [Secret](https://github.com/Azure/azure-sdk-for-go/blob/119eb84d/services/keyvault/2016-10-01/keyvault/models.go#L2287-L2304)
* [Cert](https://github.com/Azure/azure-sdk-for-go/blob/119eb84d/services/keyvault/2016-10-01/keyvault/models.go#L275-L296)
* [Key](https://github.com/Azure/azure-sdk-for-go/blob/119eb84d/services/keyvault/2016-10-01/keyvault/models.go#L1649-L1660)

For example, if you wanted to read the `Value` attribute of a `Secret` whose name was `test`, the template for that would be: `{{ .Secrets.test.Value }}`

## Other fields

Other worker-level fields that you can specify are:

* `frequency`: How often the worker should poll its resources and see if there are any changes. Defaults to 60s
* `preChange`: If the newly rendered sink contents differ from the file contents already on disk, the command specified here will be executed before the file is written
* `postChange`: If the newly rendered sink contents differ from the file contents already on disk, the command specified here will be executed after the file is written

# Examples

### SSL Cert + Private key

When you create a Cert in azure key vault, it automatically creates a Secret and Key with the same name. In the associated Secret, the value will be a blob that contains both the private key and cert.

To fetch the private key, you'll need to ensure that the Secret is in your resources section. You will also need to use the built-in `privateKey` and `cert` helpers to parse the blob into its respective pieces.

Note: `cert` helper will only return the leaf certificate

In the example below, it is assumed you have created a PEM format certificate with the name `pem-test`:

```yaml
workers:
  -
    resources:
      - kind: secret
        name: pem-test
        vaultBaseURL: https://test-kv.vault.azure.net/
    frequency: 60s
    postChange: service nginx restart
    sinks:
      - path: ./pem-test.key
        template: '{{ index .Secrets "pem-test" | privateKey }}'
        owner: myuser
        group: mygroup
        mode: 0600
      - path: ./pem-test.cert
        template: '{{ index .Secrets "pem-test" | cert }}'
```

Complete List of Cert Helpers:

`cert` - returns PEM formatted leaf certificate.

`privateKey` - returns PEM formatted private key.

`issuers` - returns sorted issuers in PEM format.

`fullChain` - returns full certificate chain including leaf cert in PEM format.

Note:
- The resource type `cert` does not contain any chain information due to the way Azure stores the data.  If you wish to use `issuers` or `fullChain` helpers, you must do so on a `secret` resource.
- The `issuers` and `fullChain` helpers will do their best to reconstruct the chain, but can only work with the data
given.  So if you did not store your certificate with its chain an empty string will be returned.
### Multiple secrets in a file

Let's suppose you had 4 secrets in a given key vault, `dbHost`, `dbName`, `dbUser`, `dbPass`.

Here's some sample configs:

#### Using individual secret lookups
```yaml
workers:
  -
    resources:
      - kind: secret
        name: dbHost
        vaultBaseURL: https://test-kv.vault.azure.net/
      - kind: secret
        name: dbName
        vaultBaseURL: https://test-kv.vault.azure.net/
      - kind: secret
        name: dbUser
        vaultBaseURL: https://test-kv.vault.azure.net/
      - kind: secret
        name: dbPass
        vaultBaseURL: https://test-kv.vault.azure.net/
    frequency: 60s
    postChange: docker restart webapp
    sinks:
      - path: ./config.yml
        template: "databaseUrl: psql://{{ .Secrets.dbUser.Value }}:{{ .Secrets.dbPass.Value }}@{{ .Secrets.dbHost.Value }}/{{ .Secrets.dbName.Value }}"
```

#### Using all-secrets
```yaml
workers:
  -
    resources:
      - kind: all-secrets
        vaultBaseURL: https://test-kv.vault.azure.net/
    frequency: 60s
    postChange: docker restart webapp
    sinks:
      - path: ./config.yml
        template: "databaseUrl: psql://{{ .Secrets.dbUser.Value }}:{{ .Secrets.dbPass.Value }}@{{ .Secrets.dbHost.Value }}/{{ .Secrets.dbName.Value }}"
```

You can also use the built-in `toValues` helper to get key/value pairs of all of your secrets.
```yaml
workers:
  -
    resources:
      - kind: all-secrets
        vaultBaseURL: https://test-kv.vault.azure.net/
    frequency: 60s
    postChange: docker restart webapp
    sinks:
      - path: ./config.json
        template: "{{ index .Secrets | toValues | toJson }""
```
will output the following in the `config.json` file
```json
{ "dbHost": "my-host", "dbName": "my-db", "dbUser": "my-user", "dbPass": "my-pass" }
```

### Resources with special characters in their name

Go's text/template syntax cannot handle reading fields with special characters (including hyphens)
in the name directly. If you have a resource with a hyphen (or other funky character) you will have
to use the built-in [index](https://golang.org/pkg/text/template/#hdr-Functions) function for fetching
the appropriate value:

```yaml
workers:
  -
    resources:
      - kind: secret
        name: my-test
        vaultBaseURL: https://test-kv.vault.azure.net/
    frequency: 60s
    postChange: service nginx restart
    sinks:
      - path: ./my-test
        template: '{{ index .Secrets "pem-test.Value" }}'
```

### Different credentials for different resources

Using two sets of credentials, one named `default` and one named `shared`, fetch multiple resources.
Don't specify any credential for the resources using `default`.

```yaml
credentials:
  -
    name: default
    tenantID: my-tenant-id
    clientID: http://cjohnson-test-spn
    clientSecret: cjohnson-test-secret
  -
    name: shared
    tenantID: my-tenant-id
    clientID: http://shared-test-spn
    clientSecret: shared-test-secret

workers:
  -
    resources:
      - kind: secret
        name: thing1
        vaultBaseURL: https://test-kv.vault.azure.net/
        # No credential specified, so "default" will be used
      - kind: secret
        name: thing2
        vaultBaseURL: https://test-kv.vault.azure.net/
        # No credential specified, so "default" will be used
      - kind: secret
        name: thing3
        vaultBaseURL: https://test-kv.vault.azure.net/
        credential: shared # Refers to credentials.name == "shared" above
    frequency: 60s
    sinks:
      - path: ./secret.txt
        template: "{{ .Secrets.thing1.Value }}{{ .Secrets.thing2.Value }}{{ .Secrets.thing3.Value }}"
```

# Workers

Workers work in a loop, whose frequency is controlled by the `frequency` field in your config. Each iteration of the loop, the worker performs the following:

* Fetch all of the specified resources
* If any errors occur, fail the iteration and:
  * For high-frequency workers (<60s) just wait for the next iteration and try again
  * For low-frequency workers (>60s), enter a retry/backoff cycle, with jitter to avoid the [thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem)
* If no errors occurred, then for each sink specified:
  * Load and/or parse the specified template, and render it using the fetched resources
  * Compare the results of the template to the contents of the destination path
  * If the contents differ, trigger any `preChange` hook, write the contents to the `path`, and trigger any `postChange` hook

# Config watcher

A filesystem watch is placed on the specified config file, and if the file is changed, the config will be re-parsed and all of the workers will be killed and recreated based on the new config

# Known Issues

- Using a 4 digit `mode` on MacOS will only support `sticky` (i.e. `1644`). `setuid` and `setgid` do not work.
