# Engineering philosophy

1. **Engineering over convenience.** Convenience is valuable only when it preserves contracts, ownership, and reversibility. A fast wrong repository is debt delivered early.
2. **Outcomes over file presence.** A `SECURITY.md` is evidence of discoverability, not proof of secure engineering. Checks distinguish documents, configuration, execution, and verified behavior.
3. **Documentation is an interface.** User-facing behavior, examples, migrations, and limitations change with code and are tested accordingly.
4. **Tests follow risk.** Mandatory testing means relevant evidence: unit, integration, contract, end-to-end, compatibility, security, performance, or scientific reproducibility—not one universal coverage number.
5. **Performance has a budget.** Startup, steady-state runtime, memory, artifact size, and benchmark methodology are explicit product constraints.
6. **Automation beats repetition; explanation beats magic.** Repeated work should become deterministic automation. Every planned mutation states what, why, source policy, and rollback.
7. **Opinionated defaults, explicit escape hatches.** Defaults minimize early choices; deviations are supported through time-bounded, owned exceptions rather than hidden forks.
8. **Native ecosystems deserve respect.** Go tests remain colocated, Python uses `pyproject.toml`, Rust uses Cargo conventions, and JS packages define exports. Consistency is in quality outcomes and command semantics.
9. **Generated code becomes user code.** Aruo never assumes continuing ownership of arbitrary files. Provenance is granular and upgrades preserve local decisions.
10. **Trust is a feature.** Least privilege, offline operation, dry runs, signed artifacts, secret redaction, plugin isolation, and auditable actions are core UX.
11. **Every project should teach.** The generated repository explains its structure, commands, tradeoffs, examples, and supported contribution path.
12. **Maintenance is part of creation.** A project without ownership, lifecycle status, release policy, dependency strategy, or archival path is incomplete.
13. **Evidence before claims.** “Fast,” “secure,” and “compatible” require reproducible measurements and scoped definitions.
14. **Small reversible steps.** Plans are atomic, idempotent where possible, checkpointed, and safe to resume.
15. **Open core means open governance.** Formats, schemas, conformance suites, and core workflows remain inspectable and portable.

When principles conflict, safety and user ownership outrank automation; correctness outranks speed; ecosystem convention outranks cosmetic uniformity; a smaller maintained feature outranks a broad unowned one.
