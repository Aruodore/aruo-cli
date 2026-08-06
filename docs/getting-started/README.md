# Getting started

Aruo has no production release yet. Contributors can run the current checkout
directly:

```sh
go run ./cmd/aruo help
go run ./cmd/aruo version
go run ./cmd/aruo create --help
go run ./cmd/aruo doctor .
```

For repeated local use, install a development binary:

```sh
go build -o "$HOME/.local/bin/aruo" ./cmd/aruo
aruo version
```

`$HOME/.local/bin` must be on `PATH`. The binary reports version `dev` because
it is built from an unreleased checkout; rebuild it after pulling changes.

The future production quick start will cover verified release installation, an
initial workflow on a disposable repository, expected output, cleanup, and the
next two workflows. It must finish in five minutes on each tier-1 platform and
run in release CI.

For now, start with the [project vision](../../VISION.md) and [architecture](../architecture/README.md). Do not install unofficial packages using the Aruo name.
