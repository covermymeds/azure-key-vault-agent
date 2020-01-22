# Usage

`akva --config=akva.yaml`

## Single execution

To run once then exit, you can use:

`akva --config=akva.yaml --once`

# Authentication

You can either specify the following as environment variables, or put them into a file called `.env`:

```bash
AZURE_TENANT_ID=<tenant id>
AZURE_SUBSCRIPTION_ID=<subscription id>
AZURE_CLIENT_ID=<SPN name including http>
AZURE_CLIENT_SECRET=<SPN password>
```

# Config

Create a yaml file which holds configuration for one or more workers. Each worker can pull one or more resources and use those resources to write to one or more file sinks.

A simple example is given below:

```yaml
workers:
  -
    resources:
      - kind: secret
        name: password
        vaultBaseURL: https://test-kv.vault.azure.net/
    sinks:
      - path: ./password
        template: "{{ .Secrets.password.Value }}"
```

## Resources

The `resources` section is a list of one or more resources to fetch. Each resource has a `kind`, `name`, and `vaultBaseURL` field. Valid kinds are: `cert`, `secret`, `key`.

## Sinks

The `sinks` section is a list of one or more files to write to. Each sink has a `path` and either `template` (inline template) or `templatePath` (path to template on the filesystem). The template syntax is golang's [text/template](https://golang.org/pkg/text/template/#hdr-Text_and_spaces) library (with [sprig](https://github.com/Masterminds/sprig) helpers).

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
        template: '{{ privateKey "pem-test" }}'
      - path: ./pem-test.cert
        template: '{{ cert "pem-test" }}'
```

### Multiple secrets in a file

Let's suppose you had 4 secrets in a given key vault, `dbHost`, `dbName`, `dbUser`, `dbPass`.

Here's a sample config:

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
      - path: ./config.json
        template: "databaseUrl: psql://{{ .Secrets.dbUser.Value }}:{{ .Secrets.dbPass.Value }}@{{ .Secrets.dbHost.Value }}/{{ .Secrets.dbName.Value }}"
```

### Resources with special characters in their name

Go's text/template syntax cannot handle reading fields with special characters (including hyphens) in the name directly. If you have a resource with a hyphen (or other funky character) you will have to use the built-in [index](https://golang.org/pkg/text/template/#hdr-Functions) function for fetching the appropriate value:

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
