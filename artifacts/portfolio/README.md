# 포트폴리오 검증 증거

2026-10-24, 동일 PC의 실제 Compose 서비스와 Chromium을 사용한 실험이다. 애플리케이션 기능 변경 없이 문서·측정·시연 도구를 정리한 범위다. [환경](environment.json), [기계 판정](final-results.json), [소스 동일성](application-source-check.json)을 함께 보존한다.

| 실행 | 판정 | 근거 |
|---|---|---|
| Go 단위·실제 MySQL/API/WS | PASS, 9개 | [통합 로그](final-verification.txt) |
| 프론트엔드·타입·빌드·정적 검사 | PASS, unit 5개 | 같은 통합 로그 및 browser-qa 단계별 파일 |
| Chromium 전체 E2E | PASS, 13개 / 59.3s, skip·flaky 0 | [JSON](../browser-qa/runs/portfolio-final-20261024/results.json) |
| Redis/Kafka/Go 실제 장애 | PASS, 68.19s | [fault 로그](../browser-qa/runs/portfolio-final-20261024/infrastructure-regression.txt) |
| npm audit / govulncheck | PASS, 해당 시점 발견 0건 | 통합 로그 |
| 20/50/100명·한 보드 | PASS, 620/1550/3100 ACK, 관측 실패 0건 | burst-final-*.json |
| 20명·5분·1Hz | PASS, 6000 ACK / 300.001s | [지속 부하](soak-20-5m.json) |
| 모든 load socket/API/Replay | PASS, 각 최종 state 일치 | 부하 JSON |
| 직접 SQL state/events/Outbox/audit | PASS, 미발행 0건·감사 수 일치 | [SQL 검증](sql-verification.json) |
| 두 서버 브라우저 시연 | PASS, 별도 1개 / 18.2s | [시연 결과](demo-results.json), [GIF](demo/syncframe-demo.gif), [WebM](demo/syncframe-demo.webm) |
| 공개 패턴 검사·문서 로컬 링크 | PASS, 검사 범위 내 발견/누락 0건 | [공개 검사](public-audit.json), [링크](document-links.json) |

주요 성능 표와 해석은 [performance.md](../../docs/performance.md), Revision·Outbox·Consumer 경계는 [review.md](../../docs/review.md), 핵심 요약과 면접 근거는 [portfolio.md](../../docs/portfolio.md)에 있다. 실패 Trace와 기존 원본은 이전 browser-qa 실행 폴더에 유지한다.

실행한 최종 검사에서 FAIL은 없었다. Consumer poison message 자동 복구는 구현되어 있지 않다. 정확한 Kafka ACK/SQL 사이 kill, 50/100명 5분, 장기 soak·실물 모바일·WAN·운영 HA·디스크 장애, Ticket-System과 동일 조건 비교는 UNVERIFIED다. 100명 ACK p95 약 1.74초를 매끄러운 100인 사용성으로 해석하지 않는다.

재현 명령은 [testing.md](../../docs/testing.md), [performance.md](../../docs/performance.md), [browser-qa.md](../../docs/browser-qa.md), [publication.md](../../docs/publication.md)에 있다. 기본 환경 변수 없이 실행해 SKIP한 인프라 테스트를 PASS로 세지 않는다.