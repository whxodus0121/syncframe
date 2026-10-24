# 공개 자료와 검증 범위

## 현재 main과 보존 정책

공개 저장소는 [whxodus0121/syncframe](https://github.com/whxodus0121/syncframe)이며 기본 브랜치는 main이다. 현재 main에서는 내부 개발 관리 문서와 반복 실행 폴더를 제외하고 [대표 증거](../artifacts/README.md)만 유지한다. 내부 기록은 로컬에 보존하고 Git에서 제외한다. 테스트 결과가 새로 생성되어도 반복 실행 폴더는 자동으로 공개 대상에 포함되지 않는다.

정리 전에 내부 기록과 추적된 산출물 579개를 ZIP으로 백업하고 각 항목의 SHA-256을 검증했으며 전체 Git bundle도 보존했다. 공개 대표 자료는 정리 전 커밋의 Git blob을 그대로 사용한다. [보존 목록](../artifacts/evidence/curation.json)에 원래 경로·커밋·blob·SHA-256을 기록했다. 최신 자동 RunId 실행 원본도 로컬에 별도 보관한다.

성능 기록의 e629601 등은 원래 로컬 개발 식별자다. [소스 동일성](../artifacts/evidence/application-source-check.json)은 공개 코드와 해당 실험 코드의 Git blob을 비교하는 근거다.

이번 정리는 일반적인 삭제·이동 커밋이다. 기존 이력을 재작성하거나 Force Push하지 않는다. 따라서 내부 문서와 이전 자료는 [정리 전 커밋](https://github.com/whxodus0121/syncframe/tree/678c58e0a0c9bce89390bbe836b4514b3d1198e4)에서 계속 보인다. 현재 main에서 제외하는 것과 과거 이력을 지우는 것은 다르다. 전체 clone 용량도 같은 비율로 줄지는 않는다.

## 검사와 공개 가능한 설정

```powershell
python scripts/audit-public.py --ref main
python scripts/check-doc-links.py
python scripts/verify-github-publication.py --repo whxodus0121/syncframe
```

[보안 검사](../artifacts/evidence/security-audit.json)는 공개 main에서 도달 가능한 Git blob, 추적·비무시 working file, ZIP 내부를 검사한다. ref를 생략하면 비공개 로컬 개발 브랜치까지 포함하므로 공개 범위 검사는 main을 지정한다. private key, token, 공유 capability·쿠키, 비샘플 DB 연결 문자열, 사용자 경로와 로컬 프로필 이름을 찾는다. 일치한 비밀값을 출력하지 않는다.

[문서 검사](../artifacts/evidence/document-links.json)는 파일의 존재뿐 아니라 Git 추적 여부도 확인한다. 로컬에만 남은 무시 파일을 공개 링크의 정상 대상으로 세지 않는다. [원격 검사 도구](../scripts/verify-github-publication.py)는 익명 HTTP와 정확한 commit/tree/blob·mode 비교를 수행한다. 과거 대규모 업로드의 검증 원본은 정리 전 커밋에 보존한다.

Compose의 syncframe:syncframe 및 root:local-root-only는 공개된 로컬 데모 값이다. 서비스 포트는 127.0.0.1에 바인딩하며 운영 비밀번호로 사용하지 않는다. 실제 .env, 인증 상태, HAR, private key, DB dump·백업은 Git에서 제외한다. 공유 링크의 64자리 capability와 sf_ 쿠키는 편집 권한이므로 공개하지 않는다.

처음 공개할 때 생성한 main은 개인 경로가 있는 비공개 개발 이력과 분리되어 있다. 그 이후의 공개 커밋은 그대로 유지한다. 정리 전에 공개된 Trace는 이미 capability와 개인 경로를 가린 버전이며 원래 실패·입력·타이밍을 보존한다.

## CI와 실제 서비스 검증

[경량 Actions](https://github.com/whxodus0121/syncframe/actions/workflows/ci.yml)는 Go 테스트·vet, 프론트엔드 테스트·TypeScript/Vite 빌드·문서 링크를 검사한다. 인프라 변수를 주입하지 않으므로 MySQL/API 검사는 SKIP하며 Chromium과 컨테이너 장애를 실행하지 않는다. job Summary에 Go PASS와 SKIP 이름·개수를 구분한다. 서비스 검증 결과를 경량 CI의 성과로 집계하지 않는다.

실제 Docker Compose 검증과 자동 RunId 재실행은 [별도 근거](../artifacts/evidence/readme-command.json)에 있다. [검증 방법](testing.md)은 두 실행 계층과 미검증 경계를 구분한다.

## 남은 한계

패턴 검사 발견 0건이 모든 형식의 비밀 부재를 증명하지는 않는다. 임의의 사설 토큰, 영상 전체 OCR, 외부 복제본은 UNVERIFIED다. Kafka poison-message 격리·Consumer 정지 감지는 미구현 후속 과제다. 현재 포트폴리오 정리에서 핵심 백엔드와 장애 동작은 변경하지 않는다.
## 라이선스

루트 프로젝트 LICENSE는 지정되어 있지 않다. 공개 저장소를 만든다는 이유로 새 라이선스를 추가하지 않는다. 추후 선택 시 간결한 허가 조건의 MIT와 명시적인 특허 허가 조항을 갖는 Apache-2.0을 대안으로 검토할 수 있다. 프로젝트의 선택은 의존성 원문의 조건을 대체하지 않는다. [MIT 원문](https://opensource.org/license/mit), [Apache-2.0 원문](https://www.apache.org/licenses/LICENSE-2.0)

[로컬 의존성 LICENSE 확인 기록](../artifacts/evidence/dependency-license-inventory.json)에 설치 버전과 원문 헤더를 보존한다. go-sql-driver/mysql은 MPL-2.0, kafka-go와 React는 MIT, lucide-react는 ISC이며 기타 Go 의존성에는 BSD·MIT·Apache 조건이 포함된다. 이 인벤토리는 로컬에서 원문을 확보한 패키지의 기록이며 전체 배포물의 법적 적합성 인증은 아니다.

서버 이미지의 조건도 클라이언트 라이브러리와 구분한다. Redis 7.4는 RSALv2/SSPLv1 이중 라이선스이고, MySQL Community 서버에는 GPLv2 조건이 있다. Kafka 4.0.1은 Apache-2.0이다. 저장소는 의존성 바이너리나 Docker 이미지를 재배포하지 않고 빌드·실행 정의를 제공한다. [Redis 7.4 원문](https://github.com/redis/redis/blob/7.4.0/LICENSE.txt), [MySQL 공식 안내](https://www.mysql.com/about/legal/licensing/oem/), [Kafka 4.0.1 원문](https://github.com/apache/kafka/blob/4.0.1/LICENSE)
