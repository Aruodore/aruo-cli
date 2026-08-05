# Repository Auditing Tools Study

Status: completed  
Reviewed: 2026-08-05

GitHub's community profile measures recommended community files and exposes the exact evidence behind its percentage. TODO Group Repo Linter separates filesystem axioms, configurable rules, fixes, and multiple renderers. OpenSSF Scorecard uses individually documented security heuristics with remediation and explicitly warns that its aggregate is opinionated rather than definitive. OpenSSF Best Practices complements automation with attestations that cannot always be inferred from a checkout.

Aruo adopts stable rule IDs, typed evidence, category transparency, local and machine renderers, and actionable remediation. It does not copy a universal “quality” claim, silently award points for unavailable remote evidence, or mix fixes into the read-only command.

Primary sources:

- [GitHub community profiles](https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/about-community-profiles-for-public-repositories)
- [GitHub community metrics API](https://docs.github.com/en/rest/metrics/community)
- [TODO Group Repo Linter](https://todogroup.github.io/repolinter/)
- [OpenSSF Scorecard](https://github.com/ossf/scorecard)
- [OpenSSF Scorecard checks](https://github.com/ossf/scorecard/blob/main/docs/checks.md)
- [GitHub Actions workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax)
