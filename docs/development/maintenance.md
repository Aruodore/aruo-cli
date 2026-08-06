# Maintenance

## Cadence

- Weekly: triage new issues/security alerts and failed scheduled workflows.
- Monthly: dependency/tool updates in bounded reviewable groups; stale exceptions and flaky tests.
- Quarterly: quick start/docs/search review, template qualification, support matrix, CI permissions/actions, benchmark baselines, maintainer/access audit.
- Per release: migrations, provenance, compatibility, artifact verification, support declarations.
- Annually: roadmap, governance, threat model, dependency/license posture, archived platforms/templates, succession plan.

## Dependencies and templates

Dependencies need a clear role, active maintenance, compatible license, security posture, and removal path. Prefer the standard library and native tools. Automated update PRs never auto-merge major or security-sensitive changes solely because tests pass.

Every supported blueprint/template has an owner, compatibility range, fixture matrix, and retirement policy. Ecosystem default changes are evaluated and rolled out through versioned upgrades—not edited silently in old artifacts. Unmaintained variants move through deprecated to archived with migration guidance.

## CI and issue care

CI is product infrastructure: pin actions, minimize permissions, monitor runner deprecations, bound caches/artifacts, test recovery, and remove redundant checks. Triage classifies impact, reproduction, area, status, and urgency without closing valid reports for imperfect formatting. Security and conduct reports use private paths.

## Support and stewardship

Support channels, response expectations, supported versions/platforms, and maintainer capacity are published honestly. Critical responsibilities have two people or an explicit succession risk. Maintainers take rotations and breaks; automation reduces toil but does not conceal unowned systems.

Archival is a responsible lifecycle state. An archived project gets a final release/status notice, read-only repository where appropriate, dependency/security warning, migration/fork guidance, and preserved history.

