# FLAC2MP3

## Synopsis

Golang cli program to convert FLAC files to MP3 files. Example project to demonstrate the usage of the semaphores
implemented by buffered channels.

## Build and run

### Prerequisites

#### Install ffmpeg

To convert audio files, you need to have `ffmpeg` installed. On macOS, you can install it via `brew` or use the package
manager of your system:

```shell
brew install ffmpeg
```

#### Linting

Install [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation) and start linting:

```shell
make lint
```

#### Commit hooks

To ensure a consistent code style, apply the linting rules to new code and run tests, we
use [pre-commit](https://pre-commit.com/). Cod
To install the commit hooks, run:

```shell
make setup
```

### Build

Build go binary:

```shell
make build
```

### Run

Run the program:

Start the program with the `-input` flag to specify the directory with the FLAC files:

```shell
./bin/flac2mp3 -input /tmp/music 
```

The mp3 files will be saved in the same base directory as the FLAC files with the extension `_320`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
