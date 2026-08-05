# Commands

Go executable entry points live here. `cmd/aruo` is deliberately thin: it owns process streams, the injected logger and build metadata, then delegates to `internal/cli`. It must not contain command or domain behavior. The current executable implements only root help, `help`, and `version`.
