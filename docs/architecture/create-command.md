# Create Command Architecture

## Contract

`aruo create [directory]` creates a new repository from one qualified catalog entry. It supports guided input and complete flag-driven automation. It never modifies an existing destination, downloads an unpinned template, runs generated code, installs dependencies, initializes Git, or publishes anything.

```text
Cobra adapter → input resolver → catalog selection → template renderer
                                              │
                                              ▼ file plan
                                    atomic repository writer
                                              │
                                              ▼
                                      success + next steps
```

The command package owns flags and prompts only. `internal/create` owns validation, catalog-neutral planning, and application orchestration. The built-in catalog owns language and workload assumptions. The template engine remains a pure renderer. The filesystem writer stages a complete tree in a sibling temporary directory, applies portable modes, then renames it to the requested destination.

## UX decisions

- The destination is the single positional argument and project name defaults to its base name.
- Interactive mode asks only for unresolved required inputs, shows defaults, and confirms the exact destination/template/file count.
- `--non-interactive` never reads stdin. Missing required input is an actionable error naming the flag.
- `--yes` accepts the final creation confirmation; it does not invent identity or publishing decisions.
- `--template` is the stable automation selector. `--language` and `--kind` filter catalog discovery without embedding language branches in the command.
- Existing destinations fail, even if empty. A future `init` command will own adoption of existing directories.
- Success output is concise and includes `cd`, native validation, and Git initialization as explicit next steps.

## Catalog contract

A catalog entry provides stable ID, display name, language, kind, description, template filesystem, blueprint, defaults, input specifications, validation, and qualification evidence. The create workflow loops over declared inputs; it does not know that Go uses modules or Python uses package names.

Built-in, organization, registry, and future plugin catalogs will implement the same lookup contract. Plugin contributions remain data bundles validated and rendered by core; they do not execute inside the command.

## Production readiness

“Production-ready” means a catalog entry has a declared standard profile and qualification tests. The first `go-library` entry emits working source and tests plus README, license, changelog, roadmap, security, contributing, code of conduct, editor policy, developer tasks, CI, issue forms, pull-request template, documentation structure, and Aruo metadata. Generation tests assert this required inventory and run native validation where the toolchain is available.

Generated content is still user-owned. `.aruo.yaml` records template identity and project intent for future audit and upgrade workflows; it does not grant Aruo blanket ownership.

## Failure and safety

All input and rendering validation happens before filesystem mutation. Writes reject absolute/traversal paths and unsupported modes. The destination parent must exist. If staging fails or cancellation occurs, the temporary tree is removed. The final rename is same-parent and atomic where the platform filesystem supports atomic directory rename. No partial destination is intentionally exposed.

Crash-left staging directories are recognizable by the `.aruo-create-*` prefix and may be safely inspected or removed. Symlink creation is unsupported in the first version.

## Extension rules

Add a language or workload by registering a qualified catalog entry, not by branching in `create`. Add an input through its entry schema. Add post-generation work as visible future plan operations, never renderer hooks. Add an overwrite/adopt behavior only through a separate reviewed workflow with conflict and recovery semantics.
