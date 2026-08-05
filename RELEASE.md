# Release process

## Version policy

Aruo uses Semantic Versioning for its documented CLI, exit codes, configuration, machine-output schemas, plugin protocol, artifact formats, and supported public Go APIs. Before 1.0, breaking changes are permitted only with explicit changelog and migration guidance; “pre-1.0” is not permission for needless churn.

## Flow

1. Contributors add reviewed release intent with user-visible changes.
2. Automation opens a release PR updating version metadata, changelog, docs, compatibility tables, and migration guides.
3. Maintainers run the full cross-platform, security, performance, artifact, quick-start, and upgrade qualification suite.
4. An authorized maintainer approves and merges the release PR.
5. CI creates a protected signed tag, builds once in trusted runners, and produces platform archives, checksums, SBOMs, signatures, provenance attestations, and source archive.
6. Verification installs and exercises those exact artifacts before publishing package-manager metadata and an immutable GitHub Release.
7. Release notes summarize outcomes, compatibility, known issues, migrations, contributors, and verification instructions.

Publishing uses short-lived trusted identity/OIDC where supported, protected environments, least privilege, and separation of approval from execution. A failed publish is resumed from verified artifacts; released tags/artifacts are never replaced.

## Deprecation and migration

Deprecations state replacement, reason, first warning version, planned removal boundary, and detection. Stable contracts receive at least one minor release of warning and are removed only in an appropriate major release unless security makes that unsafe. Machine-readable schema/protocol support windows are published. Breaking releases include tested migration guides and, where feasible, `aruo migrate` plans.

## Rollback

Package versions are immutable. Faulty releases are deprecated/yanked only when ecosystem rules and user safety justify it, followed by a new corrective version. Security releases follow coordinated disclosure and may use an abbreviated public process without skipping verification.

