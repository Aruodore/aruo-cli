# ADR-0005: Standardize outcomes, preserve native layouts

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

Mature repositories share ownership, documentation, testing, release, security, and compatibility outcomes, while Go, Python, Rust, and JavaScript use different idiomatic layouts and tools.

## Considered options

- One universal tree/toolchain: consistent appearance, high ecosystem friction.
- No shared standard: native freedom, inconsistent quality and fleet visibility.
- Shared evidence policy with native adapters: consistent outcomes without cosmetic uniformity.

## Decision

Policy packs define required outcomes/evidence. Language/workload adapters select native structure and tooling. Exceptions are scoped, owned, justified, and expiring.

## Consequences

Cross-language reports require normalized observations and calibrated checks. Adapter maintainers must understand their ecosystems. File presence alone cannot satisfy substantive policies.

## Validation

Design partners review finding usefulness; high-severity false positives and exception rates are tracked.

