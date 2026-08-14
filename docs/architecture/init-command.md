# Init command architecture

`aruo init [repository]` adopts an existing application by installing an engineering contract for AI agents and humans. It does not scaffold application features, add dependencies, execute repository code, or claim production readiness.

## Ownership

- Aruo-managed: `AGENTS.md` and `.aruo/**`. `.aruo/managed.json` records contract version and SHA-256 hashes for every managed file except itself.
- Application-owned: `aruo.yaml` and all application code. Aruo creates the initial intent file but never lists it as managed.

This boundary enables a future `aruo update` to replace unchanged managed rules while preserving application decisions. Update and merge behavior are not implemented yet.

## Workflow

1. Resolve and open the existing repository without running it.
2. Detect Go, Python, or Node from local manifests; for Node, detect supported framework dependencies and lockfile-selected package manager.
3. Render the embedded, versioned core contract and stack observation entirely in memory.
4. Report every proposed file. `--dry-run` stops here.
5. Refuse if any destination already exists. There is no force flag.
6. Materialize files in same-filesystem staging, create `.aruo` exclusively, and link top-level files with no-replace semantics. Roll back paths created by this operation if commit fails.

The first version deliberately avoids merging an existing `AGENTS.md`. Safe composition and managed updates require an explicit design; silently appending or overwriting instructions would make ownership ambiguous.

## Installed contract

The core contract defines intent principles and completion requirements. Focused rules cover architecture and dependencies, security, APIs and errors, data and migrations, testing, observability, and delivery. Stack detection is descriptive in v1: it gives agents context but does not create a growing matrix of framework-specific application templates.

Doctor validates managed hashes separately from application intent and the repository-health score. A modified or missing managed rule is a blocking contract finding; an application-owned `aruo.yaml` is expected to evolve.
