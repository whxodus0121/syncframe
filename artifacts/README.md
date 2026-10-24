# 공개 검증 근거

현재 main에는 대표 시연과 판정에 필요한 원본만 둔다. 전체 반복 실행 자료는 로컬 검증 백업과 [정리 전 공개 커밋](https://github.com/whxodus0121/syncframe/tree/678c58e0a0c9bce89390bbe836b4514b3d1198e4/artifacts)에 보존한다. Git 이력을 재작성하지 않으므로 이전 자료는 과거 커밋에서 계속 조회할 수 있다.

| 자료 | 목적 |
|---|---|
| [GIF](portfolio/demo/syncframe-demo.gif) · [WebM](portfolio/demo/syncframe-demo.webm) | 두 서버 공동 편집 → Offline pending/복구 → 읽기 전용 Time Machine |
| [모바일 전](evidence/screenshots/mobile-before.png) · [후](evidence/screenshots/mobile.png) | 초기 History가 도구를 가리던 결함과 수정 화면 |
| [데스크톱](evidence/screenshots/desktop.png) · [태블릿](evidence/screenshots/tablet.png) | 실제 반응형 속성 편집 |
| [Offline](evidence/screenshots/offline-pending.png) · [Overlay](evidence/screenshots/history-overlay.png) | 연결 상태와 과거/현재 위치 비교 |
| [대표 실패 Trace](evidence/failure-draft-overwrite/trace.zip) | 원격 커밋이 미전송 초안을 덮던 결함의 수정 전 재현 |
| [명령 재실행](evidence/readme-command.json) | 인수 없는 QA 명령의 자동 RunId와 실제 전체 실행 결과 |
| [대표 E2E JSON](evidence/browser/results.json) · [실행 로그](evidence/browser/browser-tests.txt) | Chromium 13개, 별도 실제 서비스 검증 |
| [Go](evidence/browser/go-tests.txt) · [frontend](evidence/browser/frontend-tests.txt) · [빌드](evidence/browser/frontend-build.txt) · [장애 회귀](evidence/browser/infrastructure-regression.txt) | 실제 실행 계층별 근거 |
| [20명](evidence/loads/burst-final-20.json) · [50명](evidence/loads/burst-final-50.json) · [100명](evidence/loads/burst-final-100.json) · [5분](evidence/loads/soak-20-5m.json) | ACK 지연·처리량·실패·모든 socket 수렴 |
| [직접 SQL 검증](evidence/loads/sql-verification.json) · [자원 표본](evidence/loads/resources.json) | Replay·Outbox·audit 비교와 자원 조건 |
| [보존 목록과 해시](evidence/curation.json) | 대표 자료의 원래 경로·Git blob·SHA-256 및 백업 검증 |
| [보안 검사](evidence/security-audit.json) · [문서 링크](evidence/document-links.json) | 공개 대상 검사 범위와 판정 |

대표 로그는 publication-final-20261024 시점의 상세 원본이다. 최신 인수 없는 명령 실행은 별도 JSON에 기록하며 실행 시간을 합산하지 않는다. E2E JSON 내부 attachment 경로는 원래 실행의 메타데이터다. 전체 첨부 파일은 정리 전 커밋과 로컬 백업에서 보존한다. 서로 다른 실행의 자료를 같은 실행의 원본처럼 합치지 않는다.

GIF는 실제 캡처 8장의 상태 전환이며 WebM은 연속 브라우저 녹화다. 실패 Trace는 capability와 개인 경로를 가린 공개본이다. 보드 상태·입력·실패 assertion·타이밍을 성공 결과로 바꾸지 않았다.

전체 CI 결과는 [Actions](https://github.com/whxodus0121/syncframe/actions/workflows/ci.yml)에 보존한다. 경량 CI는 MySQL/API 통합·Chromium·인프라 장애를 실행하지 않으며 Go SKIP을 PASS로 세지 않는다. Docker 검증과 성능 수치의 해석은 [testing](../docs/testing.md), [browser QA](../docs/browser-qa.md), [performance](../docs/performance.md)에 있다.