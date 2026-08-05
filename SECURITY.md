# Security policy

## Supported versions

No production release exists. Design documents receive corrections on `main`; they do not constitute a supported executable. Before 1.0, each release will publish its support status. After 1.0, the latest minor of the current major receives security fixes, with older lines supported only when explicitly listed here.

## Reporting a vulnerability

Use GitHub private vulnerability reporting when enabled. A dedicated `security@` address and PGP key must be verified before the first binary release and added here. Do not open a public issue for suspected command execution, traversal, plugin, credential, artifact-signing, or supply-chain vulnerabilities.

Include affected version/commit, platform, reproduction, impact, prerequisites, and suggested mitigation if known. Remove real secrets and personal data. Maintainers target acknowledgement within three business days, an initial assessment within seven, and status updates at least every fourteen days; these targets may change with published support capacity.

## Disclosure

We coordinate remediation, advisory, CVE request where appropriate, release, and credit with the reporter. Disclosure timing considers exploitability and downstream readiness. We do not demand secrecy after a reasonable remediation period or threaten good-faith research.

## Security design

- planning is read-only and does not execute repository code;
- templates cannot hide executable hooks;
- plugins are pinned, verified, permissioned child processes;
- paths are checked against traversal and symlink escape;
- credentials use environment/keychain/provider references and are redacted;
- releases are reproducible where feasible, checksummed, signed, attested, and accompanied by an SBOM;
- CI actions and dependencies are pinned, reviewed, and updated through PRs;
- automation uses least-privilege tokens and protected environments.

Dependencies are minimized, licensed, scanned, and updated in bounded batches. Critical reachable vulnerabilities block release; exceptions require owner, rationale, compensating control, and expiry.

