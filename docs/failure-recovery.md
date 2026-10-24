# 장애와 복구

| 장애 | 사용자/저장소 동작 | 자동 검증 |
|---|---|---|
| 중복 Operation | 원래 이벤트 ack, Revision 불변 | 동일 ID 동시 20회, payload 변조 거부 |
| 누락/역순 event | 미래 Revision 보류, SQL catch-up | reducer 24개 순열, WS 누락 이력 |
| 브라우저 연결 종료 | 지수 backoff, 마지막 Revision 복구, pending 같은 ID 재전송 | Playwright 두 context 중 하나 offline/online |
| Go 재시작 | 메모리 상태에 의존하지 않고 SQL 복원 | 실제 `docker compose restart api`, state/권한/이력 비교 |
| Redis 중단 | 커서·전체 presence 제한, DB polling으로 영속 편집 전파 | 실제 stop/start, 별도 서버의 Revision 수신과 degraded/restored presence |
| Kafka 중단 | SQL 편집 성공, outbox 대기; 감사 지연 | 실제 stop/start, pending 확인, 전량 감사 반영 및 broker 중복 재전달 |
| MySQL 트랜잭션 실패 | ack 없음, event/state/outbox 모두 rollback | outbox INSERT trigger가 앞선 event/state write 뒤에 실패 |
| MySQL 전체 중단 | readiness 503, 편집은 retryable 거부/연결 종료 | 트랜잭션 실패는 검증; 장시간 전체 DB 중단의 브라우저 복구는 별도 경계 |
| 느린 소켓 | bounded ack/ephemeral queue, write deadline 후 닫음 | shutdown/race 검사; 장시간 느린 독자 부하 실험은 미수행 |

브라우저는 `offline` 이벤트에서 소켓을 닫고, 18초간 서버 메시지가 없으면 재연결한다. 서버는 5초마다 ping/Presence, pong 기반 40초 read deadline을 사용한다. DB 이벤트 조회는 1초이며 장애 시 재접속을 유도한다. 커서는 40ms 서버 throttle/60ms 브라우저 throttle을 적용하고 저장하지 않는다.

Presence payload는 Redis 15초 TTL, membership sorted set은 갱신 시 만료 항목을 제거하고 자체 30초 만료를 둔다. 서버 재시작/소켓 비정상 종료 후에는 TTL 안에 접속 표시가 잠시 남을 수 있다. Redis 복구 뒤 살아 있는 연결은 다음 5초 갱신에서 재등록한다. Redis가 없으면 해당 Go 프로세스의 접속자만 보여주며 degraded 표시를 보낸다.

## 장애 테스트의 안전 범위

Consumer의 잘못된 JSON은 해당 goroutine 종료로 이어지고, 존재하지 않는 이벤트 FK는 SQL 재시도를 계속한다. readiness는 MySQL 연결만 검사하므로 감사 정지를 감지하지 않는다. 자동 격리·DLQ·worker 감독은 현재 보장하지 않는다. 또한 audit payload와 원본 event의 전체 내용 동일성은 검사하지 않아 내부 producer를 신뢰한다. [경계 검토](review.md)에 코드 근거와 미검증 항목을 정리했다.

`TEST_DOCKER_FAULTS=1`과 `TEST_COMPOSE_PROJECT=syncframe`(또는 syncframe-verify)을 명시해야 stop/restart가 실행된다. 이 테스트는 현재 서비스에 실제 장애를 발생시키므로 개인 데모에서만 실행한다. 테스트 종료 cleanup은 Redis/Kafka를 다시 start한다. 다른 프로젝트 컨테이너나 시스템 전체 Docker 상태를 정리하지 않는다. `verify.ps1 -Fresh`만 지정된 데모 볼륨을 삭제한다.

SQL 실패 주입은 임의 운영 endpoint 대신 테스트 전용 trigger를 사용한다. 앱 계정은 그대로 두고 `TEST_ADMIN_DSN` 계정으로 trigger만 생성·삭제한다. 전역 `log_bin_trust_function_creators`를 변경하거나 앱에 SUPER를 부여하지 않는다.

## 후속 개선

Kafka poison message 처리와 Consumer 정지 감지는 후속 개선 사항이다. 현재는 잘못된 JSON에서 goroutine이 종료되고 FK 오류에서는 재시도가 계속된다. 격리·DLQ, worker 상태 지표와 경보, 명시적인 재처리 정책을 우선 검토한다. 이 정리에서는 핵심 백엔드나 해당 장애 동작을 변경하지 않았으며, 자동 복구가 구현되었다고 주장하지 않는다.

## 복구 상태 조회

```sh
docker compose logs --tail 50 api api2
docker compose exec kafka /opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server kafka:9092 --group syncframe-audit-v1 --describe
```

SQL의 `outbox WHERE published_at IS NULL`과 `audit_events` 수를 함께 본다. publish 표시가 완료되어도 consumer가 아직 처리하지 않았을 수 있다. Kafka 전달 완료와 감사 반영 완료를 구분한다. DB 백업, MySQL 복제/HA, Kafka 다중 broker, 디스크 손상·볼륨 삭제 복구는 구현·검증 범위가 아니다.
