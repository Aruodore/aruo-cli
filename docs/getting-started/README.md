# Getting started

This page is for contributors working from a checkout. To install a release,
see the [README](../../README.md#installation).

Run commands directly from the checkout:

```sh
go run ./cmd/aruo help
go run ./cmd/aruo version
go run ./cmd/aruo init --help
go run ./cmd/aruo create --help
go run ./cmd/aruo doctor .
```

For repeated local use, install a development binary:

```sh
go build -o "$HOME/.local/bin/aruo" ./cmd/aruo
aruo version
```

To install Aruo's contract into an application without touching its code or dependencies:

```sh
cd /path/to/application
aruo init --dry-run
aruo init --yes
aruo doctor
```

`$HOME/.local/bin` must be on `PATH`. A binary built this way always reports
version `dev`; rebuild it after pulling changes. See the README's
[Quick start](../../README.md#quick-start) for the same workflow against a
release binary instead.

Start with the [project vision](../design-principles/vision.md) and
[architecture](../architecture/README.md). Do not install unofficial packages
using the Aruo name.
