# AI agent entry point

Before changing this repository, read `.aruo/contract.yaml`, the rule files relevant to the affected capability, and application-owned guidance such as `AGENTS.local.md` when present.

In this contract, `MUST` is required. `SHOULD` is the default but may be declined with a recorded reason. A conditional rule applies only when the repository contains, or the task introduces or changes, the capability it names. Do not add infrastructure solely to satisfy an inapplicable rule.

Stay within the requested scope. Preserve existing architecture and project conventions unless the task requires changing them. Never weaken validation, authorization, security controls, meaningful tests, or quality gates merely to make a change pass. Never commit secrets.

Obtain explicit authorization before irreversible deletion, destructive data migration, credential or access-policy changes, or architectural expansion outside the task. If authorization cannot be requested, do not perform the action; report it as blocked.

Before completion, run repository-provided checks relevant to the change, inspect the diff, align tests and documentation with behavior, and report what was and was not verified. Update only `aruo.yaml` capability entries materially affected by the task; preserve application-specific and unknown fields.
