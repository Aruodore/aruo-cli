# Testing and completion

- MUST test externally meaningful changed behavior and relevant failure paths without coupling tests unnecessarily to private implementation details.
- Test depth MUST be proportional to risk: isolate logic, test important boundaries with real integrations where practical, and reserve end-to-end tests for critical cross-boundary journeys.
- A bug fix SHOULD include a stable regression test. If that is infeasible or disproportionate, document why and describe the alternative verification.
- Before completion, run repository-provided checks relevant to the changed area, such as formatting, linting, type checking, tests, builds, or dependency/security checks. Do not add tooling solely to satisfy this list.
- MUST inspect the final diff for unrelated changes, generated artifacts, secrets, weakened controls, and mismatched tests or documentation.
- MUST report commands run, results, checks unavailable or not run, known limitations, and any new operational responsibility.
