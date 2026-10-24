# Development stages

1. [x] Inspect tools; define consistency, security and verification contracts.
2. [x] Implement transactional MySQL event store, snapshots and conflict rules.
3. [x] Implement REST, authenticated WebSockets, presence and reconnect catch-up.
4. [x] Implement Redis hints and Kafka outbox/idempotent audit consumer.
5. [x] Build React canvas, collaboration and read-only history playback.
6. [x] Add unit, real integration, browser, fault and load tests.
7. [x] Rebuild from clean project resources, review and document measured results.

Acceptance: all twelve requested scenarios have executable checks; reports distinguish tested, failed and unverified. No durable edit is acknowledged before SQL commit. History replay equals live state. Redis/Kafka outages cannot delete committed edits.
