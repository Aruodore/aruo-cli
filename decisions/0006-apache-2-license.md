# ADR-0006: License Aruo under Apache-2.0

- Status: Accepted
- Date: 2026-08-04
- Decision makers: Aruo project lead

## Context

Aruo is developer infrastructure intended for adoption by individuals, companies, and other open-source projects. The license must permit integration and distribution while setting a clear contributor patent framework.

## Considered options

- MIT: short and permissive, but lacks Apache-2.0’s express patent terms.
- Apache-2.0: permissive with explicit copyright/patent grants, patent termination, and notice obligations.
- MPL-2.0/GPL: stronger reciprocity, but may limit embedding/adoption in developer platforms.

## Decision

License Aruo source and documentation under Apache License 2.0 unless a third-party asset or generated-project choice is explicitly identified otherwise. Generated projects select their own reviewed license; using Aruo does not force Apache-2.0 on them.

## Consequences

Distributions must carry the license and applicable notices. Dependency/template licensing must be tracked. Trademark rights are separate. Legal counsel should review material licensing changes; this ADR is not legal advice.

## Validation

CI will identify license metadata, scan dependencies/artifacts, and verify release archives include the license.

