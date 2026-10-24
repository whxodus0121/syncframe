# 공개 범위와 증거 관리

## 공개 브랜치와 원본 이력

공개 브랜치는 main이다. 공개 전 전체 로컬 이력을 검사하면서 과거 테스트 로그·Trace에서 Windows 사용자 경로를 발견했다. 파일을 지우는 후속 커밋만으로는 과거 blob이 사라지지 않으므로, 검사한 현재 트리를 새 공개 루트 커밋으로 구성한다. 기존 개발 브랜치와 원본 이력은 로컬에 보존하며 공개 브랜치의 부모로 연결하지 않는다. 기존 커밋을 재작성하거나 원격에 강제 Push하지 않는다.

기존 실험의 e629601·61b7011 등은 측정 당시의 로컬 개발 식별자다. 공개 Git 이력의 커밋을 뜻하지 않는다. [소스 동일성 기록](../artifacts/portfolio/application-source-check.json)은 측정 코드와 공개 코드의 Git blob identity를 비교하는 근거다. 소스·설정·의존성은 공개 준비 과정에서 변경하지 않았다.

## 민감정보 검사

```powershell
./scripts/sanitize-browser-traces.ps1
python scripts/sanitize-public-artifacts.py
python scripts/audit-public.py --ref main --output artifacts/publication/security-audit.json
python scripts/check-doc-links.py
```

검사기는 지정한 ref에서 도달 가능한 Git blob, 추적·비무시 working file, ZIP 내부를 검사한다. ref를 생략하면 모든 로컬 ref를 검사한다. private key, GitHub/cloud token, 공유 capability가 포함된 링크·JSON·쿠키, JWT Bearer, 비샘플 DB 자격증명, 사용자 홈 경로와 로컬 프로필 이름을 찾는다. 일치한 비밀값을 출력하지 않고 분류와 위치만 남긴다. Git 작성자 주소도 공개용 noreply인지 별도 확인한다.

원본 개발 브랜치까지 포함한 로컬 전체 검사와 실제 공개 ref 검사는 범위가 다르다. Push는 main 하나만 대상으로 한다. 원본 보존용 로컬 ref·bundle·백업을 함께 전송하는 --all 또는 --mirror를 사용하지 않는다.

[공개 이력 검사 결과](../artifacts/publication/security-audit.json), [파일 크기 검사](../artifacts/publication/file-inventory.json), [코드·측정·영상 동일성](../artifacts/publication/integrity.json)에 공개 전 gate의 근거를 남긴다.

[개인정보 치환 기록](../artifacts/publication/privacy-redactions.json)에는 변경 파일과 ZIP 항목을 남긴다. 로컬 경로와 개인 이메일만 가리고 assertion·타이밍·문서 상태는 유지한다. 원본은 Git에서 제외한 private 백업에 보관한다. 이미지·영상·주요 부하 결과 362개 파일은 치환 전후 SHA-256이 일치했다. 기존 실패 재현을 성공 결과로 덮어쓰지 않는다.

## 공개 가능한 설정과 권한

Compose의 syncframe:syncframe 및 root:local-root-only는 의도적으로 공개된 로컬 데모 값이다. 포트는 127.0.0.1에 바인딩한다. 운영 비밀번호로 사용하지 않는다. 실제 .env, 인증 상태, HAR, private key, DB dump와 private 백업은 Git에서 제외한다.

보드 ID와 Operation ID는 권한 토큰이 아니다. 공유 링크의 64자리 capability와 sf_ 쿠키는 편집 권한이며 공개하지 않는다. 시연은 주소창·개발자 도구·공유 모달 없이 페이지 화면만 녹화했다. 쿠키나 브라우저 storage state 파일을 공개 산출물로 내보내지 않는다.

## 업로드 대상과 한계

검증에 필요한 원시 JSON, 실패 Trace·영상·스크린샷은 공개용 치환본을 유지한다. node_modules, 실행 파일, web/dist, 개인 설정과 DB dump는 업로드하지 않는다. PLAN/PROGRESS/AGENTS는 개발 관리 기록이며 README의 주요 기술 산출물로 소개하지 않는다. 기존 볼륨과 사용자 데이터는 초기화하지 않았다.

패턴 검사 발견 0건이 모든 형식의 비밀 부재를 증명하지는 않는다. 임의의 사설 토큰 형식, 영상 전체 OCR, 외부 복제본은 UNVERIFIED다. GitHub 공개 후에는 익명 HTTP 접근, 원격 tree/commit 일치와 README 이미지 로딩을 별도로 확인한다.

## GitHub 공개 접근 검증

[whxodus0121/syncframe](https://github.com/whxodus0121/syncframe)는 개인 계정의 Public 저장소이며 기본 브랜치는 main이다. 최초 원격 검증에서 636개 파일의 Git blob·mode와 tree가 로컬 커밋에 일치했다. 인증 헤더 없이 README·기술 문서·GIF·WebM·대표 화면을 조회해 HTTP 200과 내용 동일성을 확인했다. 로그인하지 않은 실제 브라우저에서도 GIF와 Mermaid 아키텍처가 표시됐다.

[익명 접근·무결성 원본](../artifacts/publication/remote-verification.json), [실제 GitHub README 화면](../artifacts/publication/github-readme.png)에 검사 당시의 커밋과 결과를 기록한다. 검증 기록을 추가한 후의 커밋은 별도로 다시 비교한다.

```powershell
python scripts/verify-github-publication.py --repo whxodus0121/syncframe
```

## 라이선스

루트 프로젝트 LICENSE는 지정되어 있지 않다. 공개 저장소를 만든다는 이유로 새 라이선스를 추가하지 않는다. 추후 선택 시 간결한 허가 조건의 MIT와 명시적인 특허 허가 조항을 갖는 Apache-2.0을 대안으로 검토할 수 있다. 프로젝트의 선택은 의존성 원문의 조건을 대체하지 않는다. [MIT 원문](https://opensource.org/license/mit), [Apache-2.0 원문](https://www.apache.org/licenses/LICENSE-2.0)

[로컬 의존성 LICENSE 확인 기록](../artifacts/publication/dependency-license-inventory.json)에 설치 버전과 원문 헤더를 보존한다. go-sql-driver/mysql은 MPL-2.0, kafka-go와 React는 MIT, lucide-react는 ISC이며 기타 Go 의존성에는 BSD·MIT·Apache 조건이 포함된다. 이 인벤토리는 로컬에서 원문을 확보한 패키지의 기록이며 전체 배포물의 법적 적합성 인증은 아니다.

서버 이미지의 조건도 클라이언트 라이브러리와 구분한다. Redis 7.4는 RSALv2/SSPLv1 이중 라이선스이고, MySQL Community 서버에는 GPLv2 조건이 있다. Kafka 4.0.1은 Apache-2.0이다. 저장소는 의존성 바이너리나 Docker 이미지를 재배포하지 않고 빌드·실행 정의를 제공한다. [Redis 7.4 원문](https://github.com/redis/redis/blob/7.4.0/LICENSE.txt), [MySQL 공식 안내](https://www.mysql.com/about/legal/licensing/oem/), [Kafka 4.0.1 원문](https://github.com/apache/kafka/blob/4.0.1/LICENSE)
