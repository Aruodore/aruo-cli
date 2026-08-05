# Benchmark assets and results

This directory will contain harness definitions, workloads, baseline metadata, and immutable raw results. `definitions/` states tasks and metrics; `results/` records environment-qualified observations. Generated binaries and scratch profiles are ignored. See [PERFORMANCE.md](../PERFORMANCE.md).

Template-rendering baselines are stored in [`results/`](results/). The engine benchmark is implemented beside its internal package so it exercises the full rendering boundary without exporting a compatibility API.
