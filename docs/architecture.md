# 아키텍처

## 서비스 경계

한 Go 프로세스에 HTTP API, WebSocket, Redis subscriber, SQL Outbox publisher, Kafka consumer를 둔다. 별도의 API gateway·편집 서버·감사 마이크로서비스를 만들지 않아 로컬 자원과 배포 복잡도를 줄였다. `multi` Compose profile은 같은 코드를 두 번 실행한다. 모든 서버는 같은 MySQL을 사용하고 프로세스 메모리에 영속 상태를 보유하지 않는다.

| 기술 | 책임 | 장애 시 제한 |
|---|---|---|
| Go | 입력 검증·권한, Revision 트랜잭션, WS 연결 및 catch-up | 연결이 끊기며 재접속 필요 |
| MySQL/InnoDB | 현재 상태, 순서 있는 이벤트, ID 중복 방지, snapshots, outbox, audit | 새 편집 승인·정확한 상태 조회 불가; 503/재시도 |
| Redis | 커서, 접속 상태 TTL, 서버 간 wake-up 힌트 | 다른 서버의 커서/전체 Presence 제한, 편집은 DB polling으로 복구 |
| Kafka | 편집의 비동기 감사 이벤트 흐름 | 감사 지연, outbox 증가; 편집·이력은 동작 |
| React/TypeScript | 조작, 연결·pending 표시, 순서 복원 reducer, 이력 재생 | localStorage pending을 같은 ID로 재전송 |

Redis Pub/Sub는 전달 보장 저장소가 아니다. [공식 Redis 문서](https://redis.io/docs/latest/develop/pubsub/)의 at-most-once 특성 때문에 알림 자체로 보드 상태를 갱신하지 않는다. 클라이언트에게 전송할 순서는 MySQL의 `(board_id, revision)`에서 읽는다.

## 영속 데이터

MySQL을 기준 저장소로 선택한 이유는 한 보드의 다음 Revision, Operation ID 중복 검사, 현재 JSON과 이력, 후처리 발행 의도를 같은 트랜잭션에서 결정할 수 있기 때문이다. Redis와 Kafka를 별도 원본으로 두면 이 상태들을 맞추기 위한 추가 합의·보상 경계가 생긴다. 대신 MySQL 장애 시 편집 승인을 중단하고 인기 보드의 쓰기 직렬화 비용을 받아들인다.

- `boards`: ID, 제목, capability SHA-256, 현재 Revision, JSON state, 생성 시각.
- `events`: `(board_id, revision)` PK, `(board_id, operation_id)` UNIQUE, canonical fingerprint와 전체 이벤트 JSON. 이벤트는 적용 결과 객체를 담는다.
- `snapshots`: `(board_id, revision)` PK. Revision 0 및 50의 배수에 저장한다.
- `outbox`: 이벤트 FK, `(board_id, revision)` UNIQUE, Kafka 발행 완료 시각.
- `audit_events`: `(board_id, revision)` PK, 감사 payload, 수신 시각. FK 위반을 IGNORE하지 않는다.

객체는 현재 상태 JSON 내부에 저장한다. 500 ID 제한에서 구현·스냅샷 일관성을 단순하게 만들지만, 매 편집마다 문서 JSON을 복사하고 다시 쓴다. 큰 문서에는 객체별 행과 별도 projection 설계가 더 적합하다.

## REST / WebSocket 계약

| 경로 | 용도 |
|---|---|
| `POST /api/boards {title}` | 보드 생성, HttpOnly 쿠키와 공유 키 반환 |
| `POST /api/boards/:id/session {token}` | 공유 키 검증 후 쿠키 교환 |
| `POST /api/boards/:id/share` | 쿠키 권한 검증 후 기존 공유 키 반환; 새 탭 재공유 |
| `GET /api/boards/:id` | 현재 state |
| `GET /api/boards/:id?revision=N` | snapshot + event로 읽기 전용 복원 |
| `GET /api/boards/:id/events?after=N` | N 이후 이벤트 최대 500개; 마지막 Revision부터 반복 조회 |
| `POST /api/boards/:id/operations` | 동일 트랜잭션 경로의 편집 API; 시드·자동 검증용 |
| `GET /api/boards/:id/ws?after=N&name=...` | 쿠키 검증, origin 검사 후 실시간 연결 |
| `GET /healthz`, `GET /readyz` | 프로세스 liveness / MySQL readiness |

클라이언트 → 서버 메시지: `operation` + Operation 또는 `cursor` + x/y. 서버 → 클라이언트: `hello`, `snapshot`, `event`, `ack`, `rejected`, `presence`, `cursor`. `after=0`에는 현재 snapshot, 그 외에는 누락 event를 전송한다. 미래 Revision의 resume 요청은 거부한다. ack와 event의 상대적 도착 순서는 보장하지 않는다.

Go의 WS writer goroutine 하나만 데이터 프레임을 쓴다. reader는 편집을 저장하고 bounded channel로 ack를 전달한다. 영속 전송 큐를 메모리에 무한히 쌓지 않고 SQL catch-up으로 대체한다. 느린 소켓은 write deadline으로 닫고 재연결하도록 한다. 종료 시 활성 연결을 닫고 reader/writer 정리가 끝난 뒤 저장소를 닫는다.

## 보안 경계

새 탭의 재공유는 기존 보드 쿠키를 검증한 `POST /share`에서 같은 capability를 반환한다. 새 권한이나 키를 만들지 않으며 same-origin POST와 no-store 정책을 적용한다. 명시적 공유를 위해 JS가 키를 다시 읽을 수 있으므로 HttpOnly만으로 XSS 시 키 유출을 막는다고 주장하지 않는다. 무권한 share 요청은 401로 검증한다.

로컬 데모의 익명 capability 모델이다. 공유 키는 256비트 난수이며 DB에는 해시만 저장한다. URL fragment의 키를 쿠키로 교환한 뒤 주소창에서 제거한다. 쿠키는 HttpOnly, SameSite=Strict, 보드별 path이며 직접 TLS 접속이면 Secure이다. 브라우저가 공유를 위해 키를 sessionStorage에도 보유하므로 XSS가 발생하면 보호되지 않는다. React 텍스트 렌더링과 CSP를 사용하고 `dangerouslySetInnerHTML`을 사용하지 않는다.

동일 출처 변경 요청 및 WS origin, 16 KiB 메시지/본문, 필드·좌표·텍스트 길이, 60 edits/sec/WS 연결, 서버별 200 연결 한도를 적용한다. SQL은 파라미터 바인딩하며 내부 DB 오류를 외부에 노출하지 않는다. REST의 공개 보드 생성에는 전역/사용자별 quota가 없다. TLS termination proxy, 키 철회, RBAC, 사용자 계정, 서비스 간 TLS/SASL, DDoS 방어는 운영 배포 전에 필요하다. Compose의 고정 비밀번호는 loopback 데모용이다.
