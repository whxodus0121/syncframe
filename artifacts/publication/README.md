# GitHub 공개 검증

저장소: [whxodus0121/syncframe](https://github.com/whxodus0121/syncframe), Public, 기본 브랜치 main.

| 검사 | 결과와 증거 |
|---|---|
| 실제 회귀 | Go 9개, frontend 5개, Chromium 13개(65.8s), Redis/Kafka/Go 장애(70.80s) PASS. [원본](validation.json) |
| 공개 대상 Git 이력 | main에서 도달 가능한 이력과 working file·ZIP 검사 PASS. [보안 검사](security-audit.json) |
| 개인정보 | 로그·Trace 41개 파일의 로컬 경로 치환, 원본은 비공개 보관. [치환 목록](privacy-redactions.json) |
| 코드·측정·시연 보존 | 앱/부하 구현 불변, 측정·미디어 362개 파일 SHA-256 동일. [동일성](integrity.json) |
| 파일 크기 | 초기 대상 약 64.19 MiB, 최대 파일 약 4.3 MiB. [목록](file-inventory.json) |
| 원격 무결성 | 최초 636개 파일 blob·mode 및 tree가 로컬과 일치. [익명 HTTP·Git 비교](remote-verification.json) |
| README 화면 | 비로그인 실제 브라우저에서 GIF와 Mermaid 표시. [캡처](github-readme.png) |
| 라이선스 | 루트 LICENSE 임의 추가 없음. [의존성 원문 확인](dependency-license-inventory.json), [선택 대안](../../docs/publication.md#라이선스) |

검증 도구의 첫 실행에서는 Git Trees API에 commit SHA를 전달한 응답의 sha가 commit ID여서 tree identity assertion이 실패했다. commit 응답의 실제 tree SHA로 조회하도록 수정한 뒤 동일 assertion과 전체 파일 비교를 통과했다. 제품 결함이나 assertion 제거로 처리하지 않았다.

remote-verification.json은 검사 당시의 불변 커밋을 지정한다. 그 후 추가한 검증 기록의 커밋도 Push 후 별도 비교한다. 원본 로컬 개발 이력은 공개 main의 조상이 아니며 원격에는 main만 전송한다. 공개 파일 경로의 사용자 정보·권한 토큰을 제외하면서 기존 실패 장면과 측정값을 보존했다.

성능 수치는 [기존 실험](../../docs/performance.md)을 사용하며 장시간 부하를 반복하지 않았다. Kafka ACK/SQL 경계의 정확한 kill, 운영 HA·실물 모바일·WAN·영상 전체 OCR 등 기존 미검증 경계는 그대로 유지한다.