# 검증 방법과 실행 근거

## 범위

공개 전 반복 실행은 publication-final-20261024다. Go 9개·frontend 5개·Chromium 13개(65.8s, skip/flaky 0)와 Redis/Kafka/Go 장애 회귀(70.80s)가 통과했다. [공개 전 실행 기록](../artifacts/publication/validation.json), [통합 로그](../artifacts/publication/verification.txt), [브라우저 JSON](../artifacts/browser-qa/runs/publication-final-20261024/results.json)에 근거를 보존한다. 아래 portfolio-final 실행은 성능 측정 시점의 별도 기록이다.

단위 테스트만으로 협업을 판정하지 않는다. 실제 MySQL 트랜잭션과 HTTP/WebSocket, 독립 Chromium 세션, Redis·Kafka·Go 컨테이너 장애를 각각 검사한다. 최신 재현 명령은 다음과 같다. 기존 볼륨을 유지하며 임의 ID의 테스트 보드를 생성한다.

```powershell
./scripts/browser-qa.ps1 -RunId portfolio-final-20261024
```

원본: [통합 실행 로그](../artifacts/portfolio/final-verification.txt), [단계별 결과](../artifacts/browser-qa/runs/portfolio-final-20261024/), [E2E JSON](../artifacts/browser-qa/runs/portfolio-final-20261024/results.json). 실행 판정과 한계는 [브라우저 QA](browser-qa.md), [최종 검증 기록](../artifacts/portfolio/final-results.json)에 있다.

## 환경

Windows amd64, Intel Core i5-10400F / 12 logical CPUs, Docker Desktop Linux engine 29.2.0 / Compose v5.0.2. Docker VM은 12 CPU / 약 7.69 GiB다. Go 1.26.10, Node 24.19.0, Playwright 1.64.1 / Chromium 156.0.8078.4를 사용한다. 서비스는 MySQL 8.4, Redis 7.4, Kafka 4.0.1이다. 브라우저 실제 버전은 각 실행의 environment.json에 보존한다. 부하 측정 조건은 [별도 실험](performance.md)이다.

## 실행 계층과 판정

| 계층 | 실제 입력과 assertion |
|---|---|
| Go model | 필드 충돌·삭제·입력 제한, gap 거부와 중복 무시 |
| 실제 MySQL | 20개 동시 생성·같은 필드 수정, 동일 요청 20회, ID 변조 거부, snapshot 50 이후 복원 |
| SQL 실패 주입 | 마지막 Outbox INSERT trigger 실패 뒤 state/revision/event/outbox 전체 rollback, 다음 승인 Revision 1 |
| 실제 HTTP/WS | 권한·origin·16 KiB 제한, 연결 종료 뒤 누락 이력, 동시 소켓 종료 정리 |
| 프론트엔드 | reducer의 역순·중복·누락 및 4개 이벤트의 24개 순열 수렴 |
| Chromium 전체 E2E | 공유·조작·3인/2서버·초안 충돌·Offline·서버 restart·Time Machine·3개 viewport·axe |
| Redis 장애 | stop/start 중 다른 서버의 영속 편집 전파, Presence degraded/restored |
| Kafka 장애 | 중단 중 편집 승인·pending Outbox, 복구 뒤 audit drain, 실제 중복 produce와 offset 확인 |
| 부하 | 20/50/100 sockets, 모든 socket state == API == full Replay, Operation ID 유일성 및 별도 SQL 확인 |

E2E는 실제 백엔드를 사용한다. 로딩 검증은 요청을 잠시 보류한 후 원래 서버로 continue하며 응답 본문을 대체하지 않는다. 화면·연속 Revision·이벤트·SQL을 assertion으로 확인하고 예상하지 않은 Console/Network 오류는 실패로 판정한다.

## 독립 명령

```powershell
$env:TEST_MYSQL_DSN='syncframe:syncframe@tcp(127.0.0.1:33060)/syncframe?parseTime=true&timeout=3s&readTimeout=5s&writeTimeout=5s'
$env:TEST_ADMIN_DSN='root:local-root-only@tcp(127.0.0.1:33060)/syncframe?parseTime=true'
$env:TEST_ALLOW_FAULTS='1'
$env:TEST_BASE_URL='http://localhost:8080'
go test -count=1 -v ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
cd web
npm ci
npm test
npm run build
npm audit
npx playwright install chromium
npm run browser:qa
cd ..
```

인프라 환경 변수가 없는 go test는 해당 테스트를 SKIP한다. 전체 서비스 장애는 다른 검사나 부하와 동시에 실행하지 않는다.

```powershell
$env:TEST_SECOND_URL='http://localhost:8081'
$env:TEST_KAFKA_BROKER='127.0.0.1:19092'
$env:TEST_COMPOSE_PROJECT='syncframe'
$env:TEST_DOCKER_FAULTS='1'
go test -count=1 -timeout 5m -v ./tests -run TestMultiServerAndInfrastructureFaults
Remove-Item Env:TEST_DOCKER_FAULTS
```

Linux race 검사:

```sh
docker build -f Dockerfile.verify -t syncframe-verify .
docker run --rm --network syncframe_default \
  -e 'TEST_MYSQL_DSN=syncframe:syncframe@tcp(mysql:3306)/syncframe?parseTime=true' \
  -e 'TEST_ADMIN_DSN=root:local-root-only@tcp(mysql:3306)/syncframe?parseTime=true' \
  -e TEST_ALLOW_FAULTS=1 -e TEST_REDIS_ADDR=redis:6379 syncframe-verify
```

race 컨테이너는 실제 DB와 in-process WS를 검사한다. TEST_BASE_URL이 없으므로 외부 API/fault suite는 SKIP하며 별도 호스트 실행이 담당한다. 외부 Compose API binary의 모든 경로가 race 계측되었다는 뜻은 아니다.

## 해석과 한계

- 최초 구현의 3개 E2E와 이전 QA 13개 결과는 해당 시점의 원본 증거다. 최신 실행과 합산해 더 많은 시나리오를 검증한 것처럼 세지 않는다.
- 중복 Kafka 전달은 실제 검증하지만 정확한 Kafka ACK/SQL commit 사이 producer kill은 UNVERIFIED다. SQL trigger 실패도 전원 손실이나 디스크 full 검증을 대신하지 않는다.
- 브라우저 Offline·실제 컨테이너 장애와 WAN 패킷 손실·장기 불안정 연결은 다르다. reducer 순열 테스트는 네트워크 프록시의 패킷 재배열 실험이 아니다.
- 자동 axe 결과는 해당 viewport·화면의 규칙만 확인한다. 실물 모바일·가상 키보드·스크린리더 평가는 UNVERIFIED다.
- 5분 제한 부하는 수일 이상의 안정성·보존 정책·HA 검증이 아니다.
- scripts/verify.ps1의 -Fresh는 데모 볼륨을 삭제한다. 데이터 보존이 필요한 환경에서는 사용하지 않는다. 공개 전 검증은 이 옵션을 사용하지 않았다.