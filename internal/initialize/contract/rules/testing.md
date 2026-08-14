# Testing and completion

- Test externally meaningful behavior and failure paths, not private implementation details.
- Match test depth to risk: unit tests for logic, integration tests for boundaries, and end-to-end tests only for critical journeys.
- Every bug fix needs a regression test unless the limitation is explicitly documented.
- Before completion run the repository's formatter, linter, type checker, tests, production build, and relevant dependency/security checks.
- Do not delete, skip, loosen, or rewrite a meaningful test simply to make a change pass.
- Report commands run, results, checks not run, known limitations, and any new production responsibility.
