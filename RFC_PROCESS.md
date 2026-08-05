# RFC process

RFCs propose significant product, compatibility, governance, security, or cross-cutting implementation changes before work commits the project to them.

## When an RFC is required

Use an RFC for new top-level commands, public APIs/schemas/protocols, persistence or network services, plugin permissions, supported platforms/languages, release policy, governance, or work likely to span multiple modules. Small reversible fixes and implementation details do not require one.

## Lifecycle

`draft → proposed → final-comment-period → accepted | rejected | withdrawn → implemented | superseded`

1. Copy [`rfcs/0000-template.md`](rfcs/0000-template.md), choose the next number, and open a draft PR.
2. Identify owner, decision deadline if real, stakeholders, alternatives, compatibility, security, operations, docs, tests, and rollback.
3. Maintainers assign reviewers and request evidence/prototypes where uncertainty is material.
4. After substantive review, an authorized maintainer starts a minimum seven-day final-comment period unless a documented security emergency requires shorter handling.
5. Decision makers record acceptance/rejection rationale. Acceptance authorizes the direction, not automatic merge of implementation.
6. Implementation PRs link the RFC; completion updates status and verifies success measures.

## Acceptance criteria

The problem is evidenced; scope/non-goals are clear; alternatives are represented fairly; affected contracts and migration are explicit; security/privacy/accessibility/performance/maintenance costs are addressed; ownership exists; testing and rollout/rollback are credible; the proposal aligns with vision and accepted ADRs or explicitly supersedes them.

Consensus is preferred, but documented maintainer decision resolves stalemate under [GOVERNANCE.md](GOVERNANCE.md). RFC history is never rewritten to hide rejected ideas.

