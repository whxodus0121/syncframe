# 정합성과 장애 경계 검토

## 검토 범위와 방법

MySQL 트랜잭션, Operation ID, Outbox, Kafka Consumer를 불변식 관점에서 재검토했다. 근거는 소스, 실제 MySQL 실패 주입, 중복 메시지 및 컨테이너 장애 테스트다. 외부 기관의 감사나 보안 인증을 의미하지 않는다. 이 검토에서는 애플리케이션 코드를 변경하지 않았다.

| 불변식 | 구현 근거 | 실행 근거 / 판정 |
|---|---|---|
| 보드의 승인 Revision은 중복·공백 없이 증가 | Store.Apply: board FOR UPDATE, event/state/outbox/snapshot 한 트랜잭션 | TestMySQLConcurrencyAndReplay, TestMySQLRollbackInjectedFailure: PASS |
| 동일 요청 재전송은 새 편집을 만들지 않음 | (board_id, operation_id) UNIQUE, 잠금 후 fingerprint 비교 | TestMySQLDuplicateAndDeletion: 동시 20회 중복 및 ID 변조 거부 PASS |
| SQL 커밋 전에 ACK하지 않음 | Store.Apply 성공 후 WS ack 채널에 전달 | 실제 WS 회귀 PASS. 커밋 응답이 끊긴 정확한 순간의 kill은 UNVERIFIED |
| Kafka 장애에도 승인 편집의 Outbox 유지 | 편집 트랜잭션의 Outbox INSERT, 별도 publisher | 실제 Kafka stop/start 및 pending→audit drain PASS |
| Kafka 재전달 시 감사 행 중복 방지 | Audit PK no-op, SQL 성공 뒤 CommitMessages | 실제 중복 produce와 consumer offset 확인 PASS |
| Redis 알림 없이 커밋 이력 복구 | WS writer의 SQL Revision 조회와 1초 polling | 두 API에서 Redis stop/start 회귀 PASS |

함수와 테스트 위치는 [구현 근거](portfolio.md), 최신 실행 로그는 [검증 방법](testing.md), 측정 원본은 [성능 실험](performance.md)에 연결했다.

## 놓치기 쉬운 장애 경계

**MySQL 커밋 결과의 모호성.** 클라이언트 timeout은 rollback의 증거가 아니다. 서버가 커밋하고 ACK 전 연결이 끊길 수 있다. 같은 Operation ID와 같은 actor/baseRevision/payload로 재시도해야 기존 이벤트를 식별한다. 새 ID는 새로운 요청이다. DB 전체 장애·전원 손실에 대한 내구성 인증은 하지 않았다.

**Outbox의 중복과 순서.** Kafka ACK 다음 SQL published 표시 전 실패하면 재발행된다. SKIP LOCKED를 쓰는 두 publisher는 같은 보드의 Revision도 역순으로 발행할 수 있다. 같은 Kafka key만으로 여러 producer의 SQL 순서까지 보장되지는 않는다. 감사 테이블은 순서 독립이며 라이브 복원은 SQL을 읽는다. 정확한 ACK/SQL 사이 프로세스 kill은 UNVERIFIED이며 중복 주입 검증과 구분한다.

**Consumer의 poison message.** 잘못된 JSON이면 해당 Consumer goroutine이 종료된다. JSON은 유효하지만 FK가 없는 이벤트는 SQL 재시도를 계속해 해당 파티션의 진행을 막을 수 있다. readiness는 MySQL만 확인하므로 이 정지를 감지하지 않는다. 신뢰된 내부 publisher만 사용하는 로컬 범위다. 운영 환경에서는 토픽 쓰기 권한, 실패 지표·경보, 검증 후 격리/DLQ와 재처리 정책이 필요하다. 공유 토픽에 악성 payload를 주입하는 실험은 하지 않았다.

**감사의 신뢰 경계.** Audit는 canonical event와 payload 전체를 비교하지 않는다. FK는 실제 (board, revision) 존재를 확인할 뿐이며 먼저 저장된 payload의 진실성까지 증명하지 않는다. 내부 생산자를 전제로 한다. 외부 producer를 허용하려면 canonical event 조회·내용 검증이나 별도 신뢰 모델이 필요하다.

**보드별 직렬화 비용.** DB pool은 프로세스당 최대 16개다. 인기 보드의 잠금 대기, full JSON 쓰기, 연결별 catch-up 조회와 감사 작업이 pool을 공유한다. 경쟁 사용자 증가에 따른 ACK 지연은 실측했지만 특정 쿼리가 지배적인 병목인지는 프로파일링하지 않았다. 객체별 행 또는 보드별 fanout/poller가 대안이며 트랜잭션·캐시 소유권 복잡도가 늘어난다.

## 판정

테스트한 범위에서 SQL 승인·연속 Revision·재전송 멱등성·장애 후 수렴이 유지됐다. 이를 exactly-once 전달, 무중단 감사 처리, 운영 HA나 모든 crash window의 증명으로 확대하지 않는다. 제약은 [장애 복구](failure-recovery.md)와 [대안](tradeoffs.md)에 명시했다.