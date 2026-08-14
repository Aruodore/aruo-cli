# Architecture and dependencies

- Understand the existing boundaries and data flow before changing them. Put code beside the behavior it owns.
- Create an abstraction only when it removes demonstrated complexity. Prefer direct code and established project conventions.
- Before adding a dependency, document the concrete capability it provides, why existing code or platform APIs are insufficient, its maintenance/security posture, and removal cost.
- Keep the application a modular monolith unless measured operational needs justify another boundary.
- Architecture changes require a stated problem, alternatives considered, migration impact, and updated documentation.
