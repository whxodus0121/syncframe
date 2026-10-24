# 검증 범위와 재현 방법

## 경량 CI와 로컬 통합 검증

[GitHub Actions](https://github.com/whxodus0121/syncframe/actions/workflows/ci.yml)는 Go 테스트·go vet, 프론트엔드 단위 테스트·TypeScript/Vite 빌드, 추적된 문서 링크를 검사한다. 외부 서비스와 비밀번호를 설정하지 않으므로 MySQL/API 통합 테스트는 SKIP하며 Chromium과 Redis/Kafka/Go 장애 테스트를 실행하지 않는다. CI의 녹색 상태를 전체 통합 검증으로 해석하지 않는다. 실행한 Go 테스트와 건너뛴 테스트 이름·개수는 각 job Summary에 기록한다.

실제 서비스 검증은 Windows PowerShell에서 별도로 실행한다.

```powershell
./scripts/browser-qa.ps1
```

RunId는 UTC 밀리초 시각과 GUID로 자동 생성한다. 기존 결과 폴더는 덮어쓰지 않으며, 동일한 명시적 RunId가 존재하면 실패한다. 위 명령을 인수 없이 실제 실행해 Go 9개·frontend 5개·Chromium 13개와 실제 컨테이너 장애 회귀를 확인했다. [명령 실행 근거](../artifacts/evidence/readme-command.json)에 RunId·실행 범위·결과를 기록한다. 이 결과는 GitHub Actions 실행과 별개다.

전체 출력은 로컬의 Git 제외 경로에 보관한다. 공개본은 [대표 증거](../artifacts/README.md)만 유지한다. 상세 대표 로그는 publication-final-20261024 시점의 [Go](../artifacts/evidence/browser/go-tests.txt), [frontend](../artifacts/evidence/browser/frontend-tests.txt), [빌드](../artifacts/evidence/browser/frontend-build.txt), [Chromium](../artifacts/evidence/browser/browser-tests.txt), [실제 장애](../artifacts/evidence/browser/infrastructure-regression.txt), [E2E JSON](../artifacts/evidence/browser/results.json)이다. 두 실행의 시간을 섞어 집계하지 않는다.

## 실행 계층과 판정

| 계층 | 입력과 assertion |
|---|---|
| Go model | 필드 충돌·삭제·입력 제한, gap 거부와 중복 무시 |
| 실제 MySQL | 20개 동시 생성·같은 필드 수정, 동일 요청 20회, ID 변조 거부, snapshot 50 이후 복원 |
| SQL 실패 주입 | 마지막 Outbox INSERT trigger 실패 뒤 전체 rollback, 다음 승인 Revision 1 |
| 실제 HTTP/WS | 권한·origin·16 KiB 제한, 연결 종료 뒤 누락 이력, 동시 소켓 종료 정리 |
| 프론트엔드 | 역순·중복·누락 및 4개 이벤트의 24개 순열 수렴 |
| Chromium | 공유·조작·3인/2서버·초안 충돌·Offline·서버 restart·Time Machine·반응형·axe |
| Redis/Kafka 장애 | 다른 서버의 SQL 복구, Outbox 대기·drain, 중복 produce와 offset 확인 |
| 부하 | 20/50/100 sockets, 모든 socket == API == Replay, 별도 SQL 상태·감사 수 비교 |

E2E는 실제 백엔드를 사용한다. 로딩 검증은 요청을 보류한 후 원래 서버로 continue하며 응답 본문을 대체하지 않는다. 예상 밖 Console/Network 오류는 실패로 판정한다. 환경·판정·한계는 [브라우저 QA](browser-qa.md), [성능 실험](performance.md)에 있다.

## 독립 명령

기본값만으로 실행하는 다음 Go 명령은 서비스 테스트를 SKIP한다. CI에서도 같은 경계를 적용한다.

```sh
go test -count=1 -v ./...
go vet ./...
cd web
npm ci
npm test
npm run build
```

실제 MySQL/API 검사에는 로컬 데모 환경을 명시한다.

```powershell
$env:TEST_MYSQL_DSN='syncframe:syncframe@tcp(127.0.0.1:33060)/syncframe?parseTime=true&timeout=3s&readTimeout=5s&writeTimeout=5s'
$env:TEST_ADMIN_DSN='root:local-root-only@tcp(127.0.0.1:33060)/syncframe?parseTime=true'
$env:TEST_ALLOW_FAULTS='1'
$env:TEST_BASE_URL='http://localhost:8080'
go test -count=1 -v ./...
```

전체 서비스 장애는 부하나 다른 검사와 동시에 실행하지 않는다.

```powershell
$env:TEST_SECOND_URL='http://localhost:8081'
$env:TEST_KAFKA_BROKER='127.0.0.1:19092'
$env:TEST_COMPOSE_PROJECT='syncframe'
$env:TEST_DOCKER_FAULTS='1'
go test -count=1 -timeout 5m -v ./tests -run TestMultiServerAndInfrastructureFaults
Remove-Item Env:TEST_DOCKER_FAULTS
```

취약점 검사와 Linux race 검사는 경량 CI의 범위에 포함하지 않는다. 과거 실행은 [고정된 공개 기록](https://github.com/whxodus0121/syncframe/tree/678c58e0a0c9bce89390bbe836b4514b3d1198e4/artifacts)에 남아 있으며 재실행 명령은 다음과 같다.

```sh
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
docker build -f Dockerfile.verify -t syncframe-verify .
docker run --rm --network syncframe_default \
  -e 'TEST_MYSQL_DSN=syncframe:syncframe@tcp(mysql:3306)/syncframe?parseTime=true' \
  -e 'TEST_ADMIN_DSN=root:local-root-only@tcp(mysql:3306)/syncframe?parseTime=true' \
  -e TEST_ALLOW_FAULTS=1 -e TEST_REDIS_ADDR=redis:6379 syncframe-verify
```

race 컨테이너는 실제 DB와 in-process WS를 검사한다. 외부 API/fault suite는 별도 호스트 검사가 담당한다. 모든 운영 경로가 race 계측되었다는 뜻은 아니다.

## 한계

- Kafka ACK/SQL commit 사이의 정확한 kill, 전원 손실·디스크 full·운영 HA는 UNVERIFIED다.
- Consumer poison message 격리와 정지 감시는 미구현이며 [후속 개선](failure-recovery.md#후속-개선)에 분류한다.
- Chromium Offline·서비스 재시작과 WAN 패킷 손실은 다른 조건이다. reducer 순열은 패킷 프록시 실험이 아니다.
- 자동 axe는 해당 화면의 규칙만 확인한다. 실물 모바일·가상 키보드·스크린리더 평가는 UNVERIFIED다.
- 5분 제한 부하는 수일의 안정성이나 데이터 보존 정책의 검증이 아니다.
- scripts/verify.ps1의 -Fresh는 데모 볼륨을 삭제한다. 기존 데이터 보존이 필요한 환경에서는 사용하지 않는다.