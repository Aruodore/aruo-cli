# Documentation standard

Documentation is source-controlled product behavior. Every page should let its intended reader answer:

1. **What is this?** Scope, audience, and outcome.
2. **Why does it exist?** Problem, constraints, and tradeoffs.
3. **How do I use or change it?** Prerequisites and ordered steps.
4. **What is a real example?** Runnable, versioned, tested input and expected output.
5. **Where is exact reference?** Types, defaults, errors, compatibility, and limits.
6. **How does it fit the architecture?** Ownership, dependencies, trust boundary, and ADRs.
7. **What commonly goes wrong?** Symptoms, causes, diagnosis, recovery, and cleanup.

Use the Diátaxis separation of tutorials, how-to guides, reference, and explanation. Cross-link instead of duplicating. Commands and config reference are generated from authoritative models and enriched manually. Examples and quick starts run in CI. Broken links, stale versions, inaccessible UI, and invalid snippets are release failures.

Write direct inclusive English, sentence-case headings, explicit units, stable terminology, and honest limitations. Avoid “easy,” “obvious,” and claims without evidence. Essential diagrams include prose equivalents; media has alternative text; the site targets WCAG 2.2 AA.

The source lives in [`docs/`](docs/README.md). [`website/`](website/README.md) owns rendering only. Every docs PR receives a preview, and every section has a maintainer. Versioned documentation follows the support policy and visibly marks archived releases.

