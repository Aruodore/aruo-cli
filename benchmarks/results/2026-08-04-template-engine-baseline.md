# Template Engine Baseline — 2026-08-04

## Purpose

Establish a reproducible starting measurement for Aruo's pure repository template renderer before caching or catalog-scale optimization.

## Workload

`BenchmarkRenderRepository` renders 40 text files from `fstest.MapFS`. Every file performs project-name and module substitution plus a conditional. Each iteration includes source reads, strict parsing, execution, destination validation, collision detection, byte-limit enforcement, and deterministic sorting. It intentionally measures the complete current engine rather than template execution alone.

## Environment

- OS/architecture: Linux amd64
- CPU: Intel Core i5-7300U @ 2.60 GHz
- Go: repository-pinned Go 1.26 toolchain compatibility
- Command: `go test -run '^$' -bench '^BenchmarkRenderRepository$' -benchmem -count=5 ./internal/templateengine`

## Results

| Sample | Time | Bytes allocated | Allocations |
| ---: | ---: | ---: | ---: |
| 1 | 827,500 ns/op | 404,535 B/op | 4,690 allocs/op |
| 2 | 780,998 ns/op | 404,539 B/op | 4,690 allocs/op |
| 3 | 783,465 ns/op | 404,535 B/op | 4,690 allocs/op |
| 4 | 781,799 ns/op | 404,538 B/op | 4,690 allocs/op |
| 5 | 841,296 ns/op | 404,542 B/op | 4,690 allocs/op |

Median time is approximately **783 µs for 40 files**, or about **19.6 µs per file** for this synthetic workload. Allocation count is stable across samples.

## Interpretation

This is not a release budget and should not be compared across machines without normalization. It demonstrates that correctness-first parsing on each render is already sub-millisecond for a small repository bundle on older laptop hardware. Parsing and reflective template execution dominate allocations; a compiled-bundle cache could reduce them, but would require immutable bundle identity and safe invalidation.

No cache is introduced from this result. Aruo should first benchmark realistic composed repositories, large documentation files, raw assets, and concurrent generation. Optimization is justified only if end-to-end profiling shows rendering is material relative to artifact verification, semantic planning, filesystem application, and native tool validation.

## Reproduction and regression policy

Run the command above on an otherwise idle system and retain all samples. Track distributions in CI only after stable benchmark hardware exists. Investigate sustained regressions greater than 20% on normalized hardware; do not fail pull requests on noisy shared runners.
