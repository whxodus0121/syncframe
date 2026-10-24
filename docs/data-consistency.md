# 데이터 정합성

## 보장하는 단위

보드 하나에서 성공한 영속 편집에는 연속된 Revision이 부여된다. 서로 다른 보드 사이의 전역 순서는 없다. Revision은 서버 도착 시각이나 클라이언트 timestamp가 아닌 **MySQL 보드 행 잠금 획득 및 커밋 순서**이다. [MySQL locking reads](https://dev.mysql.com/doc/refman/8.4/en/innodb-locking-reads.html)의 `FOR UPDATE` 방식으로 번호 발급과 기록을 보호한다.

```text
BEGIN (READ COMMITTED)
  SELECT board state FOR UPDATE
  IF operation ID exists:
    same fingerprint → return original event (no new revision)
    different fingerprint → reject conflict
  validate operation against current state
  reject future baseRevision; accept older baseRevision under explicit merge rule
  apply field patch or create tombstone
  INSERT event at revision + 1
  UPDATE board revision and full state
  INSERT outbox row
  INSERT snapshot if revision % 50 == 0
COMMIT
ACK
```

중간 실패는 전체 rollback한다. 응답이 사라진 뒤 재전송하면 UNIQUE operation identity로 이전 이벤트를 반환한다. 커밋 결과가 네트워크 장애로 모호해져도 **동일 ID와 동일 payload**로 재시도해야 한다. 이미 승인된 Operation의 baseRevision·actor까지 바꾸면 다른 요청으로 판정한다.

## 충돌 규칙

브라우저 속성의 미전송 text/숫자 초안은 원격 state 수신으로 덮지 않는다. blur하면 별도 Operation을 전송한 뒤 아래 서버 규칙을 적용한다. 입력 소실을 막는 보호이며 문자 단위 병합은 아니다.

| 상황 | 규칙 |
|---|---|
| 서로 다른 ID 생성 | 각각 커밋, 두 객체 모두 존재 |
| 같은 객체의 x와 color 수정 | patch가 명시한 필드만 갱신해 둘 다 유지 |
| 같은 필드 동시 수정 | 더 높은 커밋 Revision의 값 유지 |
| 전체 텍스트 동시 수정 | text 한 필드에 마지막 커밋 우선; 문자 단위 병합 없음 |
| 삭제 이후 오래된 update | conflict, 새 Revision 없음 |
| 삭제된 ID로 create | conflict; tombstone ID 재사용 금지 |
| 이미 성공한 create를 삭제 후 재전송 | 이전 event를 ack; 현재 tombstone은 변경하지 않음 |
| 동일 ID 다른 payload | conflict; 기존 event 보존 |
| 미래 baseRevision | conflict; 클라이언트 상태 오류 |

이 모델은 오프라인 의도나 사용자별 인과관계까지 보존하지 않는다. 오래된 patch도 생존한 객체에는 적용하므로, 다른 사용자가 먼저 이동한 위치를 재연결 직후 다시 덮어쓸 수 있다. 더 엄격한 도메인이라면 필드 버전 조건 또는 stale-write 거부가 대안이다.

## 메시지 순서와 복구

예를 들어 Revision 10에서 x=100, color=#3765fd인 객체에 A가 x=200을, B가 color=#16a57a을 보내면 11과 12의 커밋 순서와 무관하게 x=200, color=#16a57a을 유지한다. B도 x=300을 보낸다면 x patch 중 더 높은 Revision의 값이 남는다. Revision 13에서 삭제한 뒤 baseRevision=10인 update가 도착하면 conflict이며 Revision 14를 만들지 않는다. 이미 승인된 create의 동일 ID 재전송은 과거 이벤트를 반환할 뿐 삭제를 취소하지 않는다.

클라이언트 reducer는 현재 Revision 이하를 무시하고 미래 이벤트를 Revision map에 보관한다. `current+1`이 있어야 적용한다. 빈 구간을 만나면 `GET events?after=current`를 반복해 복구한다. 4개 이벤트의 24개 순열과 중복·삭제 이후 stale 수신을 단위 테스트한다. 정상 WS도 SQL을 순서대로 읽으므로 알림의 순서와 무관하다. 마지막 알림 하나가 유실된 경우 gap 알림이 추가로 오지 않아도 서버의 1초 polling이 복구한다.

초기 snapshot은 이미 커밋한 최신 상태이다. 재연결 시 클라이언트의 마지막 Revision을 보내 누락 이력을 받는다. ack 수신은 pending 큐에서 ID를 제거하고, 문서 변경은 snapshot/event 경로만 담당한다. ack가 event보다 먼저 와도 최종 상태는 event를 통해 일치한다.

## 스냅샷과 타임머신

Revision N 이하의 가장 최근 snapshot을 읽고 순서대로 suffix를 적용한다. 50 간격이므로 정상 데이터에서 최대 49개 이벤트 적용으로 과거 상태를 재구성한다. 내부 Store.At 조회는 최대 1,000개를 읽지만 N을 넘는 이벤트는 적용하지 않는다. 외부 events API의 page 크기는 500개이다. 따라서 적용 개수 제한이 SQL 조회량까지 49개로 제한한다는 뜻은 아니다. 과거 조회는 SQL read만 수행한다. React는 live state와 history preview를 별도로 유지하며, history 도구는 영속 Operation을 만들지 않는다.

## Kafka 전달

Outbox와 Kafka ack는 원자적이지 않다. Kafka가 승인한 직후 SQL publish 표시 전에 crash하면 중복 발행된다. consumer는 감사 행을 먼저 저장하고 offset을 커밋하며 PK로 효과를 한 번만 남긴다. 이는 at-least-once delivery + idempotent SQL effect이며 exactly-once가 아니다. [Kafka 공식 설계](https://kafka.apache.org/40/design/design/)를 참고했다.

두 publisher가 `SKIP LOCKED`로 작업하므로 Kafka 도착 순서가 보드 Revision과 다를 수 있다. 감사 처리는 순서에 의존하지 않고, 편집의 정렬 기준은 MySQL이다. 후속 분석에서 순서가 필요하면 소비 측에서 Revision을 정렬/버퍼링하거나 보드별 publisher 소유권을 추가해야 한다. 현재는 그 복잡도를 도입하지 않았다.
