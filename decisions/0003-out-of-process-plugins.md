# ADR-0003: Run plugins out of process

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

Plugins need language independence and fault containment. In-process packages inherit the CLI’s memory, credentials, and filesystem authority and couple compatibility to Go APIs.

## Considered options

- In-process Go plugins: fast calls, but poor portability/versioning and unsafe trust boundary.
- WASM sandbox: promising isolation, but capability/runtime complexity and ecosystem constraints.
- Child processes over JSON Lines: universal and inspectable, with serialization/process overhead.

## Decision

Use child processes, version-negotiated JSON Lines, typed messages, deadlines/cancellation, manifests, and explicit capabilities. Evaluate WASM later for stricter portable sandboxing.

## Consequences

Third parties can use any language. Core owns all writes/output and can terminate failures. Protocol schemas, permission brokerage, signing, and conformance testing are required.

## Validation

Test malicious/invalid messages, timeouts, cancellation, oversized output, permission denial, protocol skew, and secret isolation.

