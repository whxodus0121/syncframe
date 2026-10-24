# SYNCFRAME — 완료 상태와 재개 안내

## 현재 상태

GitHub Public 생성 및 최초 Push 성공: https://github.com/whxodus0121/syncframe, main, 06e2603. Description/Topics 9개 설정 확인. 익명 API/HTTP에서 main commit/tree와 636개 파일 SHA·mode 일치, README/docs/GIF/WebM/대표 화면 내용 일치 확인. 실제 비로그인 브라우저 GIF/Mermaid 렌더링 확인 및 github-readme.png 저장. 원본 master=61b7011 보존, force push 없음. 공개 검증 기록을 커밋한 뒤 main만 다시 Push하고 최종 HEAD 일치를 확인한다.

GitHub 공개 준비: publication-final-20261024 회귀 PASS(Go 9, frontend 5, Chromium 13/65.8s, fault 70.80s). 기존 측정 재실행 없이 원본 일치 확인. 과거 로그/Trace의 사용자 경로를 발견해 공개용 치환본 41개 파일을 생성하고 원본은 ignored private 폴더 및 bundle에 보존했다. 기존 master/개발 이력을 재작성하지 않고, 검사한 tree로 별도 main 루트 커밋을 만들어 main만 공개한다. LICENSE는 임의 추가하지 않았다. GitHub 생성/Push 및 원격 검증은 보안 gate 통과 뒤 진행한다.

2026-10-24 20:10 KST: 공개 포트폴리오 검증·문서 편집. 애플리케이션/설정/lock 21개 파일은 e629601과 Git blob 동일(artifacts/portfolio/application-source-check.json). 제품 기능 변경 없음. 측정 도구 cmd/soak, SQL/공개/문서 검사와 Chromium 시연 자동화만 추가했다.

최신 통합 회귀: portfolio-final-20261024 exit 0. Go 9개, frontend 5개, Chromium 13개(59.3s, skip/flaky 0), 실제 Redis/Kafka/Go 장애 회귀(68.19s), vet/build/npm audit/govulncheck PASS. 예상 Console/Network 6건, 예상 외 0건. 시연 별도 1개 PASS(18.2s).

20/50/100명 단기 각 620/1550/3100건, 처리량 162.69/132.31/124.32 edits/s, p95 207.72/763.13/1735.73ms. 20명·5분(1Hz)은 6000건, p95 169.16ms, 실패 0건. 모든 socket, full Replay, API 일치 및 별도 SQL 상태/이벤트/Outbox/audit 비교 PASS. 주요 결과와 이전 반복 원본을 artifacts/portfolio에 보존한다. 초기 성능 도구와 변경 도구를 서버 최적화 전후로 해석하지 않는다.

README와 docs 전체 문체/설계 근거/한계를 편집하고 docs/portfolio.md, publication.md를 추가했다. 실제 두 서버 공동 편집·Offline pending/복구·Time Machine WebM 및 8장 GIF와 화면 검토를 남겼다. 기존 E2E가 덮어쓴 legacy screenshot은 이번 캡처를 portfolio/regression-screenshots에 복사한 뒤 Git의 기존 원본을 복구했다. 실패 원본은 삭제하지 않았다.

잔여: poison JSON consumer 종료, FK 오류 재시도 정체, readiness의 감사 worker 미감시, audit canonical payload 미대조는 코드상 한계로 문서화했다. 50/100명 5분, Kafka ACK/SQL 정확한 kill, 실물 모바일/WAN/HA/디스크 장애는 UNVERIFIED. Ticket-System 소스와 동일 조건 수치 비교도 UNVERIFIED. GitHub push는 수행 범위에 포함하지 않는다.

2026-10-24 19:21 KST: 실제 브라우저 종합 QA 및 자동 수정 완료. `./scripts/browser-qa.ps1 -RunId final-20261024` exit 0. 전체 Chromium E2E 13개 PASS(59.8s), Go 실제 DB/API/WS 테스트 9개 PASS, frontend 5개 PASS, 실제 인프라 장애 회귀 PASS(67.06s). Chromium 156.0.8078.4 / Playwright 1.64.1입니다.

새 탭 공유 실패, 미전송 초안 덮어쓰기, 속성 패널의 드래그 가림, 모바일 초기 History 가림, Canvas 복귀 무동작, 이력 로딩 중 live 오표시, 빈 안내/활동 loading 표시를 수정했습니다. 데스크톱·태블릿·모바일 및 history/offline/error 화면을 이미지 도구로 직접 검토했습니다. 예상 외 Console/Network 오류는 없고 의도적 Offline/401은 별도 기록했습니다.

최신 결과/수정/환경/한계: docs/browser-qa.md 및 artifacts/browser-qa/runs/final-20261024/. 실패 Trace/영상/스크린샷은 이전 실행 폴더에 보존하며 테스트 capability는 마스킹했습니다. 새 실행은 `./scripts/browser-qa.ps1`로 재현합니다. 기존 데이터·볼륨은 이번 QA에서 삭제하지 않았습니다.

2026-10-24 KST: 로컬 협업 서비스와 자동 검증 경로 구현 완료. 빈 디렉터리에서 Git 초기화 후 단계별 commit을 남겼습니다. 새 MySQL/Kafka 볼륨에서 `./scripts/verify.ps1 -Fresh` exit 0을 확인했습니다. 아래 미수행 검증/제품 한계는 별도로 유지합니다.

현재 실행: http://localhost:8080 (API 1), http://localhost:8081 (API 2). Compose의 5개 컨테이너가 healthy입니다. 앱 옆 브라우저에는 실제 스타터 보드를 열어 두었습니다. 초기 개발용 native Go 프로세스와 오래된 실행 파일은 정리했습니다.

## 구현 완료

- React/TypeScript 랜딩, 보드 생성/공유, 도형·노트·텍스트 생성, 드래그 이동/핸들 크기 조절, 속성 편집, 삭제, 키보드 이동.
- 실제 WebSocket 협업, Redis 커서/Presence TTL, 연결·degraded·pending 표시.
- MySQL 보드별 Revision 직렬화, 필드별 last-committed-write wins, Operation ID 멱등성/fingerprint 충돌, 삭제 tombstone.
- snapshot + event suffix 복원, 이력 목록, revision 이동, Play/Pause/속도/이전·다음, 현재 위치 overlay 비교, 읽기 전용 preview.
- 재연결 catch-up, 순서 buffer, 같은 ID의 pending 재전송, offline 처리/수신 watchdog.
- SQL Outbox, Kafka 후처리, 감사 consumer PK 멱등성과 처리 후 offset commit.
- Docker Compose, schema 초기화, 환경 변수 예시, 검증·부하 도구, 설계/검증/성능/대안/리뷰 문서와 실제 스크린샷.

## 최초 구현 시점의 검증 기록

| 검증 | 결과 | 증거 |
|---|---|---|
| clean Compose build, 새 DB, 5개 서비스 | PASS | artifacts/verification.txt, services-final.json |
| Go unit/실제 DB/in-process WS/API 보안·복구 9개 | PASS | artifacts/verification.txt |
| Redis/Kafka 실제 중단·재시작, Go restart, 다중 서버 | PASS, 66.78s | 별도 fault test, artifacts/verification.txt |
| frontend reducer | PASS, 5개 | artifacts/verification.txt |
| TypeScript strict/unused, Vite build, go vet | PASS | artifacts/verification.txt |
| Chromium 협업/오프라인/이력, 반응형, axe | PASS, 3개 | artifacts/e2e-results.json, e2e-final.txt (13.6s) |
| 커서, 핸들 resize, 키보드 이동 | PASS | 최종 E2E와 screenshots |
| Linux Go race + 실제 MySQL/WS/종료 | PASS | artifacts/race.txt; 외부 API/fault는 SKIP, 위 별도 suite로 검증 |
| npm audit / Go 1.26.10 govulncheck | PASS, 발견 0건 | verification.txt / security-go.txt |
| 1/5/20 사용자 부하 | PASS, 806 편집 / 실패 0 / 최종 상태 일치 | artifacts/load-results.json |

20명 단계: 620/620 성공, 4.331초, 143.17 events/s, ACK p95 226.77ms. 짧은 단일 보드 실험이며 운영 용량이 아닙니다. 최초 full clean run 후 앱 소스 변경 없이 커서/핸들/키보드 E2E assertion을 보강하고 전체 브라우저 suite를 재실행했습니다.

## 발견하고 수정한 문제

- ServeMux fallback 충돌로 초기 서버 시작 실패.
- offline 뒤 기존 소켓이 남아 Live 표시 유지.
- SQL 내부 오류 문자열 의존/노출 가능성 → typed 오류 및 generic 503.
- Outbox next-key lock 가능성 → publisher READ COMMITTED.
- INSERT IGNORE의 FK 실패 은폐 가능성 → 중복 PK만 no-op.
- hijacked socket 종료 → socket tracking + WaitGroup.
- trigger 권한 실패 → 앱 권한 변경 없이 테스트 관리자 연결.
- Vitest의 Playwright 수집 → unit include 경계.
- 초기 의존성과 Go 1.26.5 advisory → 패치 버전/Go 1.26.10.
- 낮은 텍스트 대비 → 색상 수정과 axe 재검증.

## 한계 / UNVERIFIED

- 문자 CRDT/OT, 계정/RBAC/초대 철회, undo/redo, 대형 문서/retention/HA는 제외했습니다.
- 보드 직렬화, full JSON 갱신, 연결별 polling, Presence/활동 이력 조회 비용이 확장 한계입니다. 500 lifetime object IDs / 100,000 events / 인스턴스당 200 WS 제한.
- 실제 모바일·Firefox/WebKit·스크린리더 평가, 5분 초과 soak, 디스크 full/전원 장애, HA, producer ack/SQL 경계 kill은 미수행입니다.
- 네트워크 순서 변경은 production reducer의 결정적 permutation 테스트이며 packet proxy 실험이 아닙니다.
- 인터넷 공개용 TLS/인증/키 철회·서비스 인증·quota는 별도 범위입니다. 현재는 loopback 데모입니다.

## 재개 및 실행

1. README.md, docs/testing.md, docs/tradeoffs.md를 읽습니다.
2. `docker compose --profile multi up -d --build --wait`로 시작합니다.
3. 회귀는 `./scripts/verify.ps1`; 새 DB가 필요할 때만 `-Fresh`를 명시합니다. 다른 프로젝트/사용자 볼륨을 정리하지 않습니다.
4. 장애 suite는 다른 DB/브라우저 테스트와 분리합니다.
5. 필수 잔여 구현은 없습니다. 후속 확장은 사용자 요구에 따라 계정/키 철회 또는 큰 보드의 저장·fanout 비용부터 별도 범위를 정합니다.
