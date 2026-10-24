# 단일 보드 부하 실험

## 목적과 측정 환경

한 보드에 사용자가 집중할 때 SQL Revision 직렬화와 WebSocket fanout이 어떤 지연을 만드는지 측정한다. 브라우저 렌더링 벤치마크나 운영 최대 수용 인원 측정은 아니다. 애플리케이션 소스는 기준 commit e629601에서 변경하지 않았으며, 측정·검증 도구만 추가했다. 따라서 실행 간 차이를 코드 최적화 효과로 해석하지 않는다.

2026-10-24 KST, Windows amd64, Intel Core i5-10400F / 12 logical CPUs, Docker VM 12 CPU / 약 7.69 GiB. Go 1.26.10 클라이언트의 GOMAXPROCS는 4다. API 2개, MySQL 8.4, Redis 7.4, Kafka 4.0.1이 같은 PC에서 동작한다. 대상은 loopback 8080 한 인스턴스이며 API 8081의 감사 worker도 켜져 있다. 부하 자체를 두 API에 분산한 실험은 UNVERIFIED다.

컨테이너 메모리 제한은 MySQL 512 MiB, Kafka 640 MiB, Redis 96 MiB, API 각각 192 MiB다. 합계 약 1.6 GiB이며 Docker·호스트 클라이언트·빌드 자원은 별도다. 새 보드를 생성하지만 DB와 Kafka 볼륨은 초기화하지 않는다. 기존 실험 데이터와 캐시가 남은 로컬 조건이다.

## 부하와 지표 정의

[cmd/soak](../cmd/soak/main.go)는 실제 HTTP로 보드를 생성한 뒤 독립 WebSocket을 연결한다. 각 socket은 계속 메시지를 읽고 snapshot/event를 실제 reducer로 적용한다. 사용자마다 미승인 Operation은 최대 하나다. 재시도·재접속으로 실패를 숨기지 않는다.

- 단기: 20·50·100명 각각 새 보드 하나, 사용자당 객체 1개 생성 + x 수정 30회. ACK 다음 즉시 전송하며 think time이 없다.
- 지속: 20명, 300초, 사용자당 1초 간격. 지연으로 놓친 tick을 무한 큐에 쌓지 않는 closed-loop 방식이다. 이는 목표 약 20 edits/s의 안정성 확인이며 포화 처리량이 아니다.
- ACK 지연: WriteJSON 직전부터 해당 Operation ID의 ACK 수신까지. SQL commit과 loopback·클라이언트 스케줄링을 포함하며 event의 타 클라이언트 도착·화면 렌더링 지연은 포함하지 않는다.
- 처리량: 성공 ACK 수 / 편집 구간 시간. 보드 생성·연결·최종 수렴·감사 drain은 시간에서 제외한다.
- p50/p95/p99: 성공 ACK 표본을 정렬해 floor((N-1)×p) 위치를 사용한다. 실패율은 실제 시도 중 실패 수다.
- 종료 판정: 모든 socket state == state API, full ordered Replay == state, Operation ID 유일성, 예상 객체와 최종 x, Revision == 성공 수. 별도 [SQL 검증](../artifacts/portfolio/sql-verification.json)은 직접 MySQL state/events와 Replay, Outbox·audit 수를 비교한다.

## 결과

다음은 동일 측정 도구로 재실행한 최종 표본이다. 각 행은 한 실행이며 신뢰구간은 산출하지 않았다.

| 조건 | 승인 / 시도 | 실패율 | 시간(s) | edits/s | ACK p50(ms) | p95(ms) | p99(ms) | 최종 수렴 |
|---|---:|---:|---:|---:|---:|---:|---:|---|
| 20명 단기 | 620 / 620 | 0% | 3.811 | 162.69 | 105.32 | 207.72 | 267.00 | PASS |
| 50명 단기 | 1,550 / 1,550 | 0% | 11.715 | 132.31 | 302.84 | 763.13 | 1,024.73 | PASS |
| 100명 단기 | 3,100 / 3,100 | 0% | 24.935 | 124.32 | 615.06 | 1,735.73 | 2,569.57 | PASS |
| 20명·5분·1Hz | 6,000 / 6,000 | 0% | 300.001 | 20.00 | 69.99 | 169.16 | 269.08 | PASS |

원본: [20명](../artifacts/portfolio/burst-final-20.json) · [50명](../artifacts/portfolio/burst-final-50.json) · [100명](../artifacts/portfolio/burst-final-100.json) · [5분 지속](../artifacts/portfolio/soak-20-5m.json).

동일 도구의 앞선 반복에서는 20/50/100명 처리량이 174.22/135.12/104.53 edits/s, p95는 171.83/714.19/2,201.21ms였다. [반복 20](../artifacts/portfolio/burst-20.json) · [반복 50](../artifacts/portfolio/burst-50.json) · [반복 100](../artifacts/portfolio/burst-100.json). 이 실행 중 호스트에서 공개 파일 패턴 검사도 수행했다. 최종 반복은 그 검사가 끝난 뒤 실행했다. 부하 측정 중에는 브라우저 suite·서비스 재시작·다른 부하를 실행하지 않았다.

기존 cmd/load로 실행한 [초기 20/50/100 표본](../artifacts/portfolio/burst-original.json)도 보존한다. 이 도구는 편집 중 ACK와 함께 메시지를 읽지만 종료 후 모든 socket 수렴은 확인하지 않는다. 실패 시 미전송 잔여 요청까지 failed에 포함하는 기존 계산 경계도 있어, 주요 표는 실제 시도와 지속 수신을 구분하는 새 도구를 사용했다. 도구 변경에 따른 차이를 서버 성능 개선으로 주장하지 않는다. 최초 1/5/20 결과는 [이전 기록](../artifacts/load-results.json)이다.

## 자원 관찰과 한계

지속 부하 시작 직후부터 종료 약 50초 후까지 docker stats를 21회 기록했다(10:50:06~10:55:47 UTC). 각 샘플 수집 뒤 15초를 쉬었으며 명령 소요 시간을 포함한 실제 간격은 약 17초다. 부하 시작은 10:49:57 UTC다. [원본](../artifacts/portfolio/resources.json) · [요약](../artifacts/portfolio/resource-summary.json) · [종료 상태](../artifacts/portfolio/container-health.txt). docker stats의 CPU 100%는 대략 코어 하나 기준이며 VM 전체 100%와 다르다. 순간 피크를 모두 포착하는 계측은 아니다.

표본상 최대 CPU는 API 73.06%, API2 6.66%, MySQL 57.48%, Kafka 119.70%, Redis 3.52%다. 마지막 MySQL 메모리는 약 497.7 MiB / 512 MiB로 여유가 적었다. 측정 직후 OOMKilled=false, RestartCount=0을 확인했다. 50·100명의 5분 지속 부하는 수행하지 않았다. 이 결과로 메모리 누수 부재나 장기 안정성을 주장할 수 없다.

사용자가 늘수록 throughput은 줄고 ACK tail이 길어진다. 보드 행 잠금, full JSON write, 사용자별 SQL catch-up, 감사 작업과 공유하는 16개 DB pool이 병목 후보이다. 프로파일·SQL wait 측정 없이 한 원인을 확정하지 않는다. 느린 네트워크, 큰 텍스트·500개 객체, 수일의 이력, 여러 보드 분산, WAN은 UNVERIFIED다.

0%는 이 표본의 관측 실패율이다. 운영 SLA나 데이터 유실 확률이 아니다. 100명에서 승인에 성공하더라도 p95 약 1.74초·p99 약 2.57초는 지연에 민감한 협업의 제약이다. 단기 수치만으로 매끄러운 100인 사용성을 주장하지 않는다.

## 재현 명령

서비스가 준비된 상태에서 다른 테스트와 겹치지 않게 실행한다. 기본 포트, PowerShell 기준이다.

```powershell
$env:GOMAXPROCS='4'
foreach ($n in @(20,50,100)) {
  go run ./cmd/soak -users $n -updates 30 -output "artifacts/portfolio/repeat-$n.json"
  if ($LASTEXITCODE -ne 0) { throw 'load failed' }
}
go run ./cmd/soak -users 20 -duration 5m -interval 1s -output artifacts/portfolio/repeat-soak.json
python scripts/verify-load-state.py artifacts/portfolio/repeat-20.json artifacts/portfolio/repeat-50.json artifacts/portfolio/repeat-100.json artifacts/portfolio/repeat-soak.json
```

verify-load-state.py는 인수로 준 report를 검사하며 생략하면 문서의 burst-final-*.json과 soak-20-5m.json을 대상으로 한다. SQL 검사 결과 파일도 새 실험에서는 별도 복사해 이전 증거를 보존한다. 부하 도구는 사용자 100명, 지속 편집 구간 5분, 추정 90,000건 한도로 제한한다. 단기는 사용자당 생성 1회와 수정 최대 50회다. 시작 준비와 진행 중 ACK 대기·종료 검증 시간은 편집 구간 제한과 별도다.