# Performance

## Philosophy and budgets

Performance is part of CLI usability and CI cost. Budgets are measured on published reference hardware and revised through ADR/RFC evidence.

- `aruo version` and top-level help: p50 under 50 ms, no network/plugins/repository scan.
- cold read-only local command on a small repository: p50 under 150 ms before native tools.
- idle resident memory: target below 30 MiB for simple commands.
- compressed release archive: target below 25 MiB.
- identical no-op check should benefit from cache without changing semantics.

## Benchmarking

Benchmarks separate startup, repository discovery, parse, policy, planning, serialization, subprocess, and network time. Suites cover small/medium/large repositories, many files, monorepos, conflict-heavy upgrades, and Windows filesystems. Reports include commit, Go/tool versions, OS, CPU, memory, storage, warmup, repetitions, p50/p95/p99, allocations, raw output, and reproduction command.

## Profiling and allocations

Use Go benchmarks, `pprof`, execution traces, escape analysis, and platform profilers. Optimize only a measured material path. Allocation work prioritizes repeated parsing, repository walks, report generation, and large plans; avoid pooling until profiles show benefit.

## Concurrency

Concurrency is bounded by resource class: CPU work, filesystem scans, subprocesses, and network requests have separate limits. Preserve deterministic output ordering. All work accepts cancellation and deadlines. Avoid goroutine-per-file designs and concurrent writes to user files.

## Caching

Cache keys include content digest, relevant config, adapter/tool version, and policy/artifact digest. Never cache secrets, failed external effects, or nondeterministic checks. Cache corruption degrades to recomputation; `aruo cache` can explain and safely prune entries.

Performance regressions beyond an accepted budget fail release qualification or require an explicit ADR-backed exception.

