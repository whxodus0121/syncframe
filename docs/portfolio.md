# 포트폴리오 소개와 구현 근거

## 프로젝트 핵심 요약

SYNCFRAME은 Go 서버가 공동 편집의 순서를 결정하고 MySQL Revision 이력으로 연결 복구와 Time Machine을 제공하는 화이트보드다. 현재 state, event, Operation ID, Outbox를 한 SQL 트랜잭션에서 다루며 Redis와 Kafka는 각각 휘발성 알림과 비동기 감사로 분리한다. 두 API와 독립 Chromium 세션, DB 실패 주입, 브로커 중단·복구 및 부하를 통해 정합성의 범위를 확인한다.

포트폴리오의 초점은 사용 기술의 개수보다 **승인 시점, 재전송 identity, 순서 복원, 후처리 중복과 장애 경계**를 코드와 실험으로 설명하는 데 있다. [성능 측정](performance.md)의 수치는 단일 PC의 표본이며 운영 용량이 아니다.

## 기술적 도전 과제 세 가지

| 문제 상황 | 기술적 원인 | 선택한 해결 | 대안과 비용 | 검증 |
|---|---|---|---|---|
| 두 서버 동시 편집과 ACK 재전송 | 번호 발급 경합, 승인 결과의 모호성 | board row lock + 연속 Revision + Operation ID/fingerprint, 한 SQL commit | 낙관적 CAS는 재시도 증가, CRDT는 병합·전파 복잡도 | 동시 20 write, 중복 20회, ID 변조, 마지막 write 실패 rollback |
| 연결이 끊긴 클라이언트와 과거 화면 | Pub/Sub 유실, 역순·중복 수신, live/history 혼합 | SQL catch-up, 연속 reducer, snapshot+event 읽기 전용 preview | 공유 poller는 비용 감소 대신 소유권·fanout 관리 필요 | 두 서버/Redis 장애, Offline pending, 24순열 reducer, Replay==SQL |
| DB 저장 뒤 Kafka 전달 실패 | 서로 다른 저장소의 비원자적 이중 쓰기 | 트랜잭션 Outbox, SQL 멱등 감사 뒤 offset commit | CDC 또는 lease batch는 운영 복잡도; Kafka 없이 SQL 감사도 가능 | 실제 Kafka stop/start, drain, duplicate produce와 offset 확인 |

같은 텍스트 필드는 마지막 커밋으로 대체된다. 문자 병합, end-to-end exactly-once, 모든 crash 순간의 내구성을 주장하지 않는다.

## Ticket-System과 다른 검증 초점

Ticket-System에서는 대규모 동시 예매 요청에서 발생하는 재고 정합성 문제와 Kafka 기반 비동기 처리 및 장애 복구에 집중했다.

SYNCFRAME에서는 기존에 학습한 Go·Redis·Kafka·MySQL 기술을 실시간 협업 서비스에 적용했다. 특히 다중 서버 환경에서의 편집 순서 보장, WebSocket 재연결 이후 상태 동기화, 이벤트 이력을 이용한 과거 상태 복원이라는 새로운 문제를 다뤘다.

| SYNCFRAME에서 보여주는 역량 | 구체적인 구현·증거 |
|---|---|
| 서버 중심 상태 동기화 | SQL을 읽는 WS writer, 두 서버의 독립 세션 화면 수렴 |
| Revision 정합성 | 보드 행 잠금, 연속 이벤트 PK, 동시 커밋과 전체 rollback |
| 이벤트 Replay | 동일 이벤트에서 live와 과거 상태 재구성, snapshot 50 간격, 이력 중 live 편집 격리 |
| MySQL Outbox | 편집과 발행 의도를 원자 저장, 브로커 중단 중 ACK와 복구 후 감사 |
| 다중 서버 협업 | 8080/8081 실제 협업, Redis 없이 SQL polling fallback |
| 실패 경계 설명 | ACK/commit 모호성, 중복 발행, 순서 변경, poison message와 readiness 한계 |

프로젝트 간 수치 비교는 동일 환경·동일 워크로드를 갖추기 전까지 UNVERIFIED다.

## 면접 예상 질문

**왜 Redis나 Kafka를 문서의 원본으로 삼지 않았는가?**

Revision·현재 state·멱등 identity·발행 의도를 한 트랜잭션에 묶기 위해 MySQL을 기준으로 삼았다. Redis는 유실 가능한 힌트다. Kafka 감사 지연은 편집 승인과 분리된다. 비용은 SQL 읽기·잠금 부하다.

**FOR UPDATE가 멱등성을 자동으로 보장하는가?**

잠금은 같은 보드의 검사와 갱신을 직렬화한다. 멱등성은 Operation ID 유일 제약과 전체 Operation fingerprint 비교가 함께 만든다. actor/baseRevision을 바꾸고 같은 ID로 보내면 conflict다.

**ACK를 못 받았다면 새 ID로 재시도해도 되는가?**

커밋되었을 수 있으므로 같은 ID와 payload로 재시도해야 한다. 중복 요청에는 원래 이벤트를 돌려준다. ACK와 event의 도착 순서는 같다고 가정하지 않는다.

**서로 다른 필드를 동시에 수정하면? 삭제와 수정이 겹치면?**

patch의 명시된 필드만 바뀐다. 같은 필드는 더 높은 커밋 Revision이 이긴다. 삭제가 먼저 커밋되면 후속 수정과 ID 재사용은 tombstone으로 거부한다.

**Kafka key가 board ID라면 Revision 순서도 보장되는가?**

여러 publisher가 SKIP LOCKED로 서로 다른 행을 잡아 역순으로 보낼 수 있다. 파티션 내 Kafka 기록 순서와 SQL Revision 순서는 다른 계약이다. 감사는 순서 독립이고 Replay는 SQL Revision을 사용한다.

**Outbox는 exactly-once인가?**

아니다. Kafka ACK 뒤 SQL 표시 전 실패하면 재발행된다. Consumer도 SQL 반영 뒤 offset commit 전 실패하면 재처리한다. PK 멱등 SQL 효과와 at-least-once 전달을 구분해야 한다.

**Consumer가 잘못된 이벤트를 만나면?**

잘못된 JSON은 goroutine 종료, FK 실패는 반복 재시도다. 현재 자동 격리·감시와 readiness 연결이 없다. 운영화 시 우선 개선할 경계이며 장애 주입으로 검증한 범위를 넘겨 주장하지 않는다.

**100명에서 무엇이 병목인가?**

ACK 지연 증가는 측정했다. 보드 잠금·전체 JSON 쓰기·fanout 조회·공유 DB pool이 후보지만 프로파일링 없이 지배적 원인은 확정하지 않았다. 보드 분산/문서 크기별 실험과 SQL 대기 지표가 다음 진단이다.

**Time Machine이 라이브 보드를 바꾸지 않는 근거는?**

Store.At는 snapshot과 events를 읽기만 한다. 브라우저는 preview와 live state를 분리한다. 다른 세션이 편집하는 동안 과거 Revision을 유지하고 복귀 후 최신 내용을 표시하는 E2E가 있다.

## 코드를 따라가는 경로

| 주장 | 구현 | 검증 |
|---|---|---|
| 편집 원자성·Revision·중복 | [Store.Apply](../internal/store/store.go#L95), [schema](../internal/store/schema.sql) | [실제 MySQL 동시성·중복·rollback](../internal/store/store_integration_test.go) |
| 충돌·tombstone·순서 | [model.Apply / Replay](../internal/model/model.go) | [Go model tests](../internal/model/model_test.go), [클라이언트 reducer](../web/src/model.test.ts) |
| snapshot 복원 | [Store.At](../internal/store/store.go#L180) | MySQL replay 및 [Time Machine E2E](../web/e2e/browser-qa.spec.ts) |
| 단일 WS writer와 복구 | [server socket](../internal/httpapi/server.go), [useBoard](../web/src/useBoard.ts) | [실제 WS·보안·장애](../tests/system_test.go) |
| Outbox·Consumer | [PublishOne / RunConsumer](../internal/pipeline/pipeline.go), [Audit](../internal/store/store.go#L216) | TestMultiServerAndInfrastructureFaults 및 TestAuditIdempotence |
| 실제 화면·영속 상태 일치 | [브라우저 assertion helpers](../web/e2e/qa/helpers.ts) | 전체 replay == API == MySQL, 독립 context의 SVG |
| 부하·모든 socket 수렴 | [측정 도구](../cmd/soak/main.go) | [실측 JSON과 조건](performance.md) |

경계 검토는 [review](review.md), 최신 실행 판정은 [testing](testing.md)에 있다.