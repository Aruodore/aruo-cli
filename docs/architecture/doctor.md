# Repository Doctor Architecture

## Purpose

`aruo doctor [repository]` performs a deterministic, read-only assessment of repository engineering health. It evaluates local evidence only: structure, documentation, workflows, tests, licensing, security policy, and GitHub configuration. It neither executes repository code nor contacts a forge.

The name intentionally broadens the earlier placeholder meaning of “installation diagnosis.” Repository health is the first mode requested by the product; future environment checks must be an explicit mode rather than silently mixed into this score.

```text
CLI → read-only repository snapshot → check registry → typed assessments
                                                   │
                                                   ▼
                                     score + recommendations
                                      ├── human renderer
                                      └── versioned JSON
```

## Evidence model

Each check has a stable ID, category, maximum points, explanation, observed evidence, and actionable recommendations. A check returns earned points from zero to its declared maximum. The engine validates plugin/check output, orders results deterministically, aggregates categories, and never lets a check print or mutate files.

Checks are heuristics. Presence is not proof of effectiveness, and absence is not always negligence. The report says what was observed and avoids claims requiring GitHub APIs, execution, or organizational context. Branch protection, actual CI history, unresolved vulnerabilities, code review practice, secret scanning, and maintainer activity are explicitly outside the local score.

## Score

The first policy is versioned as `aruo.repository-health/v1` and totals 100 points:

| Category | Weight | Intent |
| --- | ---: | --- |
| Completeness | 20 | Essential lifecycle and community artifacts exist |
| Documentation | 20 | Users and contributors can understand and operate the project |
| CI | 15 | Workflows validate changes with safe dependency references and permissions |
| Tests | 15 | Tests exist, are documented, and are invoked by CI |
| License | 10 | Reuse permission is present and locally recognizable |
| Security | 10 | Private reporting and dependency maintenance are documented/configured |
| GitHub | 10 | Contribution, dependency, ownership, and release automation are represented |

Scores are rounded only for display; category points remain explicit. Grades are `A` (90–100), `B` (80–89), `C` (70–79), `D` (60–69), and `F` below 60. A default minimum of 80 makes `doctor` exit with findings code 3 while still producing the complete report. Operational failures remain exit 1.

The score is not comparable across future policy versions without migration notes. New rules cannot silently change v1 results; a new policy version or explicitly documented reweighting is required.

## Recommendations

Recommendations identify the missing or weak outcome, why it matters, and a concrete next action. They are ordered by lost points, then stable check ID. Aruo does not generate blanket “add file” advice when substance is the problem. Human output emphasizes deductions; JSON includes passing and failing evidence for automation.

## Security and filesystem boundary

The scanner opens the target with `os.OpenRoot` and consumes its `fs.FS`, preventing traversal through symlinks outside the repository. It skips `.git`, dependency caches, and generated vendor trees. Reads are size-bounded. No template, workflow, hook, test, package manager, Git command, or plugin executable is run.

## Plugin extension

Core owns the `Check` contract and score validation. Built-in checks implement it in process. Future out-of-process audit plugins receive a bounded, permissioned observation request through the versioned JSON Lines protocol and return typed assessments. They do not receive Go interfaces, arbitrary filesystem access, scoring authority beyond manifest-declared check IDs/weights, or terminal output.

Plugin results must declare policy compatibility, category, maximum points, evidence, and remediation. Core rejects duplicate IDs, unknown categories, out-of-range points, oversized evidence, protocol skew, and undeclared checks. Organization policy packs may select or reweight checks only under a new named score profile.

## Extension rules

1. Prefer observable outcomes over tool-name checks.
2. State false-positive and false-negative limits.
3. Separate local evidence from forge/runtime evidence.
4. Keep language-specific checks behind detector/adaptor registration.
5. Add fixture tests for pass, fail, partial, malformed, and non-applicable cases.
6. Version any scoring change that affects longitudinal comparison.
