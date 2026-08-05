# Risk analysis

| Risk | Failure mode | Mitigation / trigger |
|---|---|---|
| Scope expansion | shallow support for everything | strict milestone gates; capability maturity labels; non-goals |
| Template combinatorics | untestable option matrix | composable capabilities, constraints, pairwise plus boundary fixtures, retire variants |
| Upgrade corruption | user changes overwritten | clean-tree recommendation, plan/diff, provenance granularity, staging, rollback journal, conflict refusal |
| False confidence | badge/file presence called quality | evidence levels and risk-based checks; qualify claims |
| Plugin compromise | credential/source exfiltration | process isolation, permissions, pin/sign, trust prompt, allowlist, revocation; no v1 marketplace |
| Policy centralization | Aruo fights native ecosystems | outcome-level contracts and adapter ownership; public exceptions |
| Configuration sprawl | unclear source of truth | documented precedence and `config explain`; one project intent file |
| Cross-platform drift | shell/path/permission failures | native filesystem/process APIs; Windows CI; case/symlink/line-ending fixtures |
| Startup regression | CLI feels heavy | Rust binary, lazy adapters/plugins/network, budgets and benchmark gate |
| Maintenance bus factor | templates become stale | named owners, support window, automated weekly qualification, lifecycle/archive states |
| Ecosystem churn | formatter/framework defaults change | pinned artifacts, compatibility matrix, adapter releases independent of core |
| Forge dependence | GitHub assumptions block adoption | forge port; local workflows first; GitHub adapter initially, GitLab later by demand |
| AI hallucination | unsafe or fabricated changes | AI proposes only; planner validates; evidence citations; opt-in; no secret access by default |
| Hosted-service economics | cloud cost or lock-in | local-first core and exportable open formats; hosted system optional |
| Adoption | “another generator” perception | lead with adopt/check/upgrade; publish migration evidence and long-term maintenance results |
| Naming | collision, pronunciation, trademark/domain risk | formal trademark/package/domain search before 0.1; reserve package/CLI names; documented pronunciation |
| License/supply chain | incompatible templates or compromised artifacts | SPDX metadata, dependency/license scan, signatures, provenance, reproducible catalog builds |
| Privacy | repository metadata leaves machine | offline default; explicit export preview/redaction; minimal optional analytics |

## Highest-risk assumptions to test

1. Semantic repository upgrades can achieve acceptable conflict rates without owning whole files.
2. A cross-language policy IR remains useful without collapsing ecosystem nuance.
3. Users will accept explicit plan/approval steps in return for safety.
4. Maintainers can support a curated artifact matrix sustainably.

Run design-partner experiments before broad implementation. Success criteria: >90% conflict-free upgrades for Aruo-managed regions in representative edited fixtures, <5% disputed high-severity audit findings, and a default create-to-green-check path under five minutes without hidden accounts.
