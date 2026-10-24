# 대안과 의도적인 범위

선택 기준은 로컬에서 일관성과 장애 경계를 재현할 수 있는 구현 복잡도이다. 작은 서비스에 Kafka까지 필요한 것은 아니며 SQL 감사만으로도 제품 기능을 구성할 수 있다. 이 구성은 Outbox·중복 발행·offset 처리의 경계를 드러내는 대신 운영 의존성을 늘린다. [코드 검토](review.md)와 [실측 결과](performance.md)를 함께 읽어야 비용을 판단할 수 있다.

| 선택 | 장점 | 비용 / 대안 |
|---|---|---|
| 보드별 SQL 직렬화 | Revision/멱등성/상태를 작은 트랜잭션으로 설명·검증 가능 | 한 인기 보드가 병목. CRDT는 오프라인 병합에 유리하지만 tombstone/compaction/전파 검증 복잡 |
| 필드 patch LWW | 서로 다른 필드를 불필요하게 덮지 않고 구현이 작음 | 같은 필드의 사용자 의도는 유실 가능. 조건부 object version은 명시적 재시도 UX가 필요 |
| JSON 현재 상태 | 로컬 크기에서 일관성·snapshot 단순 | 문서 크기에 비례한 write amplification. 객체별 행/projection이 확장에 유리 |
| Redis hint + SQL polling | Pub/Sub 유실과 Redis 재시작에도 커밋 복구 | 연결당 1 query/sec, idle 비용. 보드별 공유 poller는 비용을 줄이지만 fanout 관리 추가 |
| 단일 Go 배포 | 낮은 운영 비용·수평 복제 가능 | 감사 부하도 같은 프로세스/DB pool 사용. 규모가 커지면 별도 worker가 격리에 유리 |
| Outbox row lock 후 동기 publish | 별도 lease 스케줄러 없이 다중 worker 안전 | Kafka timeout 동안 worker SQL 연결 유지. READ COMMITTED로 삽입 gap lock을 피하지만 대규모에는 claim/lease batch가 적합 |
| Kafka 감사 흐름 | 편집 승인과 후처리 지연 분리, 재처리 가능한 스트림 | 작은 서비스의 기능만 보면 SQL 감사 테이블로도 충분. Kafka는 장애·중복 경계를 검증하는 포트폴리오 목적도 있음 |
| snapshot 50간격 | 과거 복원 suffix 길이 제한 | 스냅샷 저장량·트랜잭션 추가. retention/compaction은 현재 없음 |
| 익명 공유 capability | 계정 서비스 없이 두 브라우저 즉시 시연 | 키 유출 시 동등 편집 권한, 철회/역할 없음. 공개 배포에는 계정+ACL 필요 |
| SVG + React | 도형 조작·커서·키보드 선택을 브라우저 기본 기능으로 구현 | 수천 객체/복잡한 stroke에는 canvas/WebGL 또는 편집 엔진이 적합 |

현재 제외: 문자 단위 CRDT/OT, freehand drawing, 연결선·그룹·이미지 업로드, undo/redo, 임의 history 분기, 사용자 계정·RBAC·링크 철회, 보드 이름 변경/삭제, 대규모 문서 virtual rendering, 자동 스토리지 compaction, 운영 모니터링/HA/TLS. 타임머신은 read-only이며 과거 상태를 live로 덮어쓰는 기능이 아니다.

부하 수치는 단일 호스트의 제한된 실험이다. 장애 없이 성공한 작은 표본을 가용성 SLO, 운영 용량 또는 “데이터 유실 0%”로 해석하지 않는다. 구현된 보장과 재현한 실험 조건만 포트폴리오에서 주장한다.
