# Terminal UX Baseline — 2026-08-06

## Purpose

Establish reproducible starting measurements for the terminal interaction
stack (`internal/tux`) introduced by
[docs/cli/terminal-ux.md](../../docs/cli/terminal-ux.md) and
[ADR-0010](../../decisions/0010-charm-v2-terminal-ux-stack.md), against the
budgets in that specification and in
[research/technology/2026-go-terminal-ux-stack.md](../../research/technology/2026-go-terminal-ux-stack.md).

## Workloads

- `BenchmarkRunVersion` (`internal/cli`) times `cli.Run` for `aruo version`
  with fresh in-memory streams per iteration. It measures Aruo's own command
  construction and execution overhead, excluding process startup (`exec`,
  dynamic linking, OS scheduling).
- `BenchmarkInterruptCancellationLatency` (`internal/tux/lifecycle`) times
  the delay between an interrupt arriving on the signal channel and the
  manager's context observing cancellation, using a fake signal channel so
  the sample excludes real OS signal delivery latency.
- A real compiled-binary measurement (`/usr/bin/env python3` timing
  `subprocess.run` around the release binary) captures full process startup,
  which the in-process benchmark above cannot.
- Binary size delta isolates the cost of `internal/tux`'s Charm dependencies
  by comparing the current stripped binary against a stripped build of
  `57631ee` (the pre-terminal-UX baseline commit, Cobra only).

## Environment

- OS/architecture: Linux amd64
- CPU: Intel Core i5-7300U @ 2.60 GHz
- Go: repository-pinned Go 1.26 toolchain (`go1.26.5`)
- Commands:
  - `go test -run '^$' -bench '^BenchmarkRunVersion$' -benchmem -count=5 ./internal/cli`
  - `go test -run '^$' -bench '^BenchmarkInterruptCancellationLatency$' -benchmem -count=5 ./internal/tux/lifecycle`
  - `go build -ldflags="-s -w" -o /tmp/aruo-release ./cmd/aruo`
  - 20 samples of `subprocess.run(["/tmp/aruo-release", "version"])`, sorted, p95 taken at rank 19

## Results

### `aruo version` — in-process command execution (`BenchmarkRunVersion`)

| Sample | Time | Bytes allocated | Allocations |
| ---: | ---: | ---: | ---: |
| 1 | 30,067 ns/op | 18,570 B/op | 143 allocs/op |
| 2 | 29,762 ns/op | 18,483 B/op | 143 allocs/op |
| 3 | 29,867 ns/op | 18,671 B/op | 143 allocs/op |
| 4 | 32,712 ns/op | 18,481 B/op | 143 allocs/op |
| 5 | 32,137 ns/op | 18,991 B/op | 143 allocs/op |

Median is approximately **30 µs**, roughly 1,600x under the 50 ms budget for
cached help/version. Allocation count is stable across samples.

### `aruo version` — full process wall time (20 samples, stripped release binary)

Median **8.03 ms**, p95 **12.47 ms**. Comfortably under the 50 ms budget;
the gap versus the in-process benchmark above is process exec/link/scheduler
overhead outside Aruo's control.

### Ctrl+C acknowledgement latency (`BenchmarkInterruptCancellationLatency`)

| Sample | Time | Bytes allocated | Allocations |
| ---: | ---: | ---: | ---: |
| 1 | 2,742 ns/op | 618 B/op | 8 allocs/op |
| 2 | 2,739 ns/op | 619 B/op | 8 allocs/op |
| 3 | 2,602 ns/op | 619 B/op | 8 allocs/op |
| 4 | 2,880 ns/op | 619 B/op | 8 allocs/op |
| 5 | 2,894 ns/op | 620 B/op | 8 allocs/op |

Median is approximately **2.8 µs** from signal-channel receipt to context
cancellation, well under the 100 ms p95 acknowledgement budget. This excludes
real OS signal delivery latency, measured instead by the subprocess fixtures
in `cmd/aruo/main_test.go`, which observed sub-second delivery and reaction
in all runs without a numeric budget assertion.

### Resident memory

`aruo version`'s child-process max RSS was **8,280 KiB** (~8.1 MiB),
measured via `resource.getrusage(RUSAGE_CHILDREN)` around a single
`fork`/`exec`. Well under the 50 MiB ordinary-CLI-path target.

### Binary size delta attributable to `internal/tux`

| Build | Stripped size (`-ldflags="-s -w"`) |
| --- | ---: |
| `57631ee` (pre-terminal-UX, Cobra only) | 4,714,761 bytes (4.50 MiB) |
| Current (`internal/tux` + Charm v2 adapters) | 7,786,761 bytes (7.43 MiB) |
| **Delta** | **3,072,000 bytes (2.93 MiB)** |

Under the research report's 5 MiB review threshold for TUX-attributable
binary growth.

## Dependency versions measured

| Module | Version |
| --- | --- |
| `charm.land/huh/v2` | v2.0.3 |
| `charm.land/bubbletea/v2` | v2.0.2 |
| `charm.land/bubbles/v2` | v2.0.0 |
| `charm.land/lipgloss/v2` | v2.0.5 |
| `github.com/spf13/cobra` | v1.10.2 |
| `github.com/spf13/pflag` | v1.0.9 (indirect) |
| `golang.org/x/term` | v0.45.0 |
| Go toolchain | go1.26.5 (module floor go1.26.0) |

## Interpretation

Help/version and Ctrl+C acknowledgement both sit far inside their budgets on
this hardware; neither is a bottleneck worth optimizing before release.
Binary size growth from the Charm v2 family is real (2.93 MiB) but under the
review threshold set before adoption. These are not release budgets and
should not be compared across machines without normalization.

Not yet measured on this baseline: first interactive paint latency for a
real Huh form under a PTY, filtering large option lists, concurrent progress
rendering with many tasks, and Windows/macOS timings. These require either a
real PTY harness or cross-platform hardware not available in this
environment; qualification coverage for the underlying behaviors (adapter
selection, signal handling, narrow-terminal degradation) exists in
`internal/tux/session`, `internal/tux/charm`, and `internal/cli`, but their
*performance* characteristics on those paths are not yet benchmarked.

## Reproduction and regression policy

Run the commands above on an otherwise idle system and retain all samples.
Investigate sustained regressions greater than 20% on normalized hardware;
do not fail pull requests on noisy shared runners.
