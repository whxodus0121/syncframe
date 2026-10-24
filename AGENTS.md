# SYNCFRAME engineering contract

- Go owns authoritative document state. MySQL is the only durable source of truth.
- Lock the board row before allocating its next revision. State, event, operation identity, snapshot and outbox commit atomically.
- Redis carries disposable presence/cursor state and wake-up hints. Always recover durable events from MySQL.
- Kafka is asynchronous audit delivery; editing never waits for it. Commit consumer offsets only after an idempotent SQL effect.
- Conflicts use server-ordered, field-level last-committed-write wins. Deleted IDs cannot be reused.
- Keep changes focused. Never fabricate test results or equate unit mocks with infrastructure fault evidence.
- Share capabilities are secrets: hash at rest, HttpOnly session cookies, same-origin WebSockets; never log credentials.
- Commands: `go test ./...`, `go vet ./...`, `cd web && npm test && npm run build`; integration: see docs/testing.md.
- Completion requires real multi-browser editing, restart recovery, rollback tests, Kafka/Redis fault tests, reproducible load evidence, and updated PROGRESS.md. Mark unavailable verification explicitly.
