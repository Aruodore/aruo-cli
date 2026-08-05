# Project Creation Experience Study

Status: completed  
Reviewed: 2026-08-04

Cargo demonstrates a compact positional path, explicit library/binary intent, VCS control, safe failure, and ecosystem-native output. npm's `create` convention demonstrates delegating to specialized initializers and forwarding arguments, but also shows the reproducibility risk of implicit latest/global packages. create-next-app pairs a recommended default path with customization and `--yes`; create-t3-app asks capability questions while supporting CI flags; Vite uses a short template selector and explicit non-interactive mode. uv and Poetry distinguish application/library/package shapes, infer safe metadata, and refuse an already initialized target.

Recurring strengths are a fast happy path, few high-value questions, discoverable template names, deterministic non-interactive flags, immediate runnable output, and concise next steps. Recurring weaknesses are hidden network/package installation, ambiguous defaults, enormous option matrices, partial directories after failure, outdated cached initializers, and scaffolds that are never upgradeable.

Aruo adopts the simple path-plus-template model, recommended defaults, interactive/non-interactive parity, refusal to overwrite, and native next steps. It adds a staged atomic write, versioned template identity, engineering-standard qualification, and a catalog boundary that can later reconcile upgrades.

Primary sources:

- [Cargo `new`](https://doc.rust-lang.org/cargo/commands/cargo-new.html)
- [npm `init` / `create`](https://docs.npmjs.com/cli/commands/npm-init)
- [create-next-app](https://nextjs.org/docs/app/getting-started/installation)
- [Create T3 App installation](https://create.t3.gg/en/installation)
- [uv project creation](https://docs.astral.sh/uv/concepts/projects/init/)
- [Poetry CLI](https://python-poetry.org/docs/cli/)
- [Vite scaffolding](https://vite.dev/guide/#scaffolding-your-first-vite-project)
