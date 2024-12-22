# FLAC2MP3

## Synopsis

Use go to convert FLAC files to MP3 files.

## Build and run

### Prerequisites

#### Linting

Install [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation) and start linting:

```shell
golangci-lint run --fix --verbose 
```

To verify the schema of the `.golangci.yml` config file run:

```shell
golangci-lint config verify --verbose --config .golangci.yml
```

#### Commit hooks

To ensure a consistent code style, apply the linting rules to new code and run tests, we use [pre-commit](https://pre-commit.com/). Cod
To install the commit hooks, run:

```shell
pre-commit install --hook-type pre-commit --hook-type pre-push
```

### Local

#### Start the application

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

