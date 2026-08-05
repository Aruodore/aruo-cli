# Modern CLI Architecture Study

Status: completed  
Reviewed: 2026-08-04

## Scope and method

This study examined official documentation and repositories for Cobra, Hugo, Helm, Docker CLI, GitHub CLI, Terraform, uv, Bun, and Deno. The tools are not treated as interchangeable: Cobra is a Go command framework; Hugo, Helm, Docker CLI, GitHub CLI, and Terraform are mature Go applications; uv, Bun, and Deno demonstrate current expectations for fast cross-platform developer tools.

## Evidence extracted

| Project | Architectural lesson | Caution for Aruo |
| --- | --- | --- |
| Cobra | Mature command trees, generated help, suggestions, argument validation, completion, and `RunE` error propagation | Its generator's globals and `init()` registration are examples, not a scalable application architecture |
| Hugo | A single distributable executable, aggressive performance focus, modular extension model, and generated command documentation are valuable | A broad command surface and configuration compatibility become long-lived costs |
| Helm | Separating reusable action logic from command wiring allows behavior to be exercised without the CLI | Environment/settings objects can become broad mutable bags unless consumers receive narrow values |
| Docker CLI | Explicit streams and command-independent client/service layers enable testing and embedding | Compatibility with a remote API and historical flags creates substantial coupling; Aruo should version boundaries early |
| GitHub CLI | Constructor-built command trees, explicit I/O/config factories, consistent help, and testable command packages scale across many commands | Large CLIs can grow deeply nested command packages; application logic must not settle in handlers |
| Terraform | A thin command front, separated internal subsystems, diagnostics, machine-readable modes, and plugin protocols support a long-lived platform | Global metadata/configuration and plugin compatibility are expensive; isolation and version negotiation must be designed before ecosystem growth |
| uv | Fast startup, concise progress, strong defaults, global caching, cross-platform standalone delivery, and familiar compatibility surfaces set a modern UX baseline | One binary replacing many tools produces a very large surface; Aruo should add lifecycle areas incrementally |
| Bun | A single fast executable and integrated workflows reduce installation and context-switching costs | Product breadth and compatibility emulation can dominate maintenance |
| Deno | Secure defaults, explicit permissions, consolidated tooling, cross-platform releases, and a coherent help surface build trust | Consolidation must retain clear internal module boundaries and stable configuration semantics |

## Recurring patterns

1. Successful tools provide one obvious executable and verb-oriented command discovery.
2. Help is generated from the same command definitions that parse input.
3. Mature Go CLIs separate command adapters from reusable operations and external clients.
4. Streams and environment access must be injectable; process globals make command tests brittle.
5. Configuration is layered, but informational commands must not pay its startup or failure cost.
6. Human output and stable machine output are separate products.
7. Errors need one rendering boundary, actionable context, and predictable exit behavior.
8. Fast startup, caching, restrained dependencies, and standalone distribution are now baseline expectations, not differentiators.
9. Plugin and compatibility surfaces outlive implementations; they require explicit protocols and version policies.

## Decisions for Aruo

- Use Cobra v1.10.2 only in the CLI presentation layer.
- Construct commands with functions; do not use global commands or `init()` registration.
- Keep `cmd/aruo/main.go` as the process composition boundary.
- Inject streams, logger, build identity, and future ports explicitly.
- Use `log/slog`; do not add a logging framework before a measured need.
- Load typed configuration lazily per use case and retain value provenance.
- Render errors once at the top and keep domain packages silent.
- Start with `aruo`, `aruo help`, and `aruo version`; defer completion and global flags.
- Track startup latency and binary size as the dependency graph grows.

## Primary sources

- [Cobra repository and documentation](https://github.com/spf13/cobra)
- [Hugo repository](https://github.com/gohugoio/hugo)
- [Helm repository](https://github.com/helm/helm)
- [Docker CLI repository](https://github.com/docker/cli)
- [GitHub CLI repository](https://github.com/cli/cli)
- [Terraform repository](https://github.com/hashicorp/terraform)
- [uv repository](https://github.com/astral-sh/uv)
- [Bun repository](https://github.com/oven-sh/bun)
- [Deno repository](https://github.com/denoland/deno)

Repository structures are observations, not endorsements of every implementation choice. Conclusions above are inferences drawn across these primary sources.
