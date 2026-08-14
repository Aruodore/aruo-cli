# Architecture and dependencies

- MUST understand and preserve existing boundaries, data flow, and project conventions unless the task requires changing them.
- MUST stay within the requested scope. Material architecture changes require a stated problem, considered alternatives, migration and compatibility impact, and documentation proportional to the decision.
- SHOULD use direct code until an abstraction removes demonstrated duplication, coupling, or complexity.
- MUST follow language, framework, operating-system, and repository naming conventions. For new otherwise-unconstrained files, SHOULD prefer lowercase kebab-case. Do not rename unrelated files solely for consistency.
- MUST NOT introduce a distributed or deployment boundary without a stated reliability, security, scaling, deployment, or ownership need that outweighs its operational cost.
- Dependency additions MUST have a concrete need, come from a trusted source, and include review of manifest and lockfile changes. Runtime, privileged, native, abandoned, or security-sensitive dependencies require proportionally deeper maintenance, integrity, and removal-risk review.
