# Security

- Treat all external input as untrusted. Validate at the server boundary and reject invalid data with stable, non-sensitive errors.
- Authenticate identity before protected work and enforce authorization server-side at every protected resource boundary.
- Never commit secrets or expose them to client bundles, logs, errors, fixtures, or examples. Keep `.env.example` non-secret and current.
- Use secure defaults for cookies, headers, redirects, file paths, and cryptographic operations. Consider CSRF whenever ambient credentials authorize state changes.
- For uploads, constrain size and type, generate storage names, isolate access, and document scanning and retention decisions.
- Record threat and abuse considerations for public or high-risk endpoints, including rate-limit topology and failure behavior.
