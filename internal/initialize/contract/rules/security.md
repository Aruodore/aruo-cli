# Security and privacy

- MUST treat data from users, networks, files, environment, dependencies, and external systems as untrusted until validated at the relevant trust boundary.
- Validation and safe use are separate controls. MUST use parameterized data access, context-appropriate output encoding, structured process APIs, safe path handling, and constrained outbound destinations wherever untrusted data reaches an interpreter, browser, command, path, redirect, or network request.
- Protected work MUST authenticate identity and authorize the requested action on the selected resource at a trusted enforcement boundary. Deny by default when identity, policy, or ownership cannot be established.
- MUST NOT commit secrets or expose them in client artifacts, logs, errors, fixtures, or examples. Example environment files, when present, must contain placeholders only and remain current.
- MUST minimize collection, storage, logging, and retention of personal or confidential data to what the capability requires.
- MUST use maintained platform or reviewed libraries for cryptography and credential handling. Never invent cryptographic protocols or use reversible encoding as protection.
- Web applications using ambient credentials MUST address CSRF for state changes and use context-appropriate cookie, header, redirect, and browser-output protections.
- File or content ingestion MUST constrain resource use, avoid trusting supplied names or types, isolate access, and record scanning and retention decisions proportionate to risk.
- Publicly reachable or repeatedly callable sensitive operations MUST consider abuse, resource exhaustion, enumeration, and rate-limit or quota behavior. Record the decision when the risk is material.
