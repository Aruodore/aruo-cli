# Architecture and dependencies

- Understand the existing boundaries and data flow before changing them. Put code beside the behavior it owns.
- Name project-owned files with lowercase kebab-case. Keep a non-kebab filename only when a language, framework, operating system, or established tool requires that exact name (for example `README.md`, `AGENTS.md`, `Dockerfile`, Go `_test.go` files, or Python modules that must be importable).
- Create an abstraction only when it removes demonstrated complexity. Prefer direct code and established project conventions.
- Before adding a dependency, document the concrete capability it provides, why existing code or platform APIs are insufficient, its maintenance/security posture, and removal cost.
- Keep the application a modular monolith unless measured operational needs justify another boundary.
- Architecture changes require a stated problem, alternatives considered, migration impact, and updated documentation.
