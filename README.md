# GreenNode CLI

The GreenNode CLI (`grn`) is a unified tool to manage your GreenNode services from the command line.

- [Getting Started](#getting-started)
- [Getting Help](#getting-help)
- [More Resources](#more-resources)

## Getting Started

### Requirements

- Python 3.10 or later (3.10.x, 3.11.x, 3.12.x, 3.13.x)

### Installation

The safest way to install the GreenNode CLI is to use `pip` in a `virtualenv`:

```bash
python -m pip install grncli
```

or, if you are not installing in a `virtualenv`, to install globally:

```bash
sudo python -m pip install grncli
```

or for your user:

```bash
python -m pip install --user grncli
```

If you have the grncli package installed and want to upgrade to the latest version:

```bash
python -m pip install --upgrade grncli
```

On Linux and macOS, the GreenNode CLI can also be installed using a [bundled installer](https://vngcloud.github.io/greennode-cli/installation/#bundled-installer). For offline environments, see the [offline install](https://vngcloud.github.io/greennode-cli/installation/#offline-install) guide.

If you want to run the `develop` branch of the GreenNode CLI, see the [Contributing Guide](CONTRIBUTING.md).

### Configuration

Before using the GreenNode CLI, you need to configure your credentials. The quickest way is to run:

```bash
grn configure
```

```
GRN Client ID [None]: <your-client-id>
GRN Client Secret [None]: <your-client-secret>
Default region name [HCM-3]:
Default output format [json]:
```

Credentials are obtained from the [VNG Cloud IAM Portal](https://hcm-3.console.vngcloud.vn/iam/) under Service Accounts.

You can also configure credentials via environment variables:

```bash
export GRN_CLIENT_ID=<your-client-id>
export GRN_CLIENT_SECRET=<your-client-secret>
export GRN_DEFAULT_REGION=HCM-3
```

Or create the credential files directly:

```ini
# ~/.greenode/credentials
[default]
client_id = your-client-id
client_secret = your-client-secret
```

```ini
# ~/.greenode/config
[default]
region = HCM-3
output = json
```

To use multiple profiles:

```bash
grn configure --profile staging
grn --profile staging vks list-clusters
```

For more configuration options, see the [Configuration Guide](https://vngcloud.github.io/greennode-cli/configuration/).

### Basic Commands

The GreenNode CLI uses a multi-part command structure:

```bash
grn <service> <command> [options and parameters]
```

For example, to list your VKS clusters:

```bash
grn vks list-clusters
```

To get help on any command:

```bash
grn help
grn vks
grn vks create-cluster help
```

To check the version:

```bash
grn --version
```

## Getting Help

The best way to interact with our team is through GitHub:

- [Open an issue](https://github.com/vngcloud/greennode-cli/issues/new/choose) — Bug reports and feature requests
- Search [existing issues](https://github.com/vngcloud/greennode-cli/issues) before opening a new one

## More Resources

- [Documentation](https://vngcloud.github.io/greennode-cli/)
- [Changelog](CHANGELOG.md)
- [Contributing Guide](CONTRIBUTING.md)
- [PyPI Package](https://pypi.org/project/grncli/)
- [VNG Cloud Console](https://hcm-3.console.vngcloud.vn/)

## License

Apache License 2.0 — see [LICENSE](LICENSE).
