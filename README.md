# GreenNode CLI

The GreenNode CLI (`grn`) is a unified tool to manage your GreenNode services from the command line.

- [Getting Started](#getting-started)
- [Getting Help](#getting-help)
- [More Resources](#more-resources)

## Getting Started

### Requirements

- No dependencies required — `grn` is a single binary

### Installation

Download the latest binary for your platform from [GitHub Releases](https://github.com/vngcloud/greennode-cli/releases):

**macOS / Linux:**

```bash
# Download (replace OS and ARCH as needed)
curl -L -o grn https://github.com/vngcloud/greennode-cli/releases/latest/download/grn-darwin-arm64
chmod +x grn
sudo mv grn /usr/local/bin/
```

**Or build from source:**

```bash
git clone https://github.com/vngcloud/greennode-cli.git
cd greennode-cli/go
go build -o grn .
sudo mv grn /usr/local/bin/
```

**Verify installation:**

```bash
grn --version
# grn-cli/0.1.0 Go/1.22.2 darwin/arm64
```

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
grn vks create-cluster --help
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
- [VNG Cloud Console](https://hcm-3.console.vngcloud.vn/)

## License

Apache License 2.0 — see [LICENSE](LICENSE).
