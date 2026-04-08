# event-collector 프로젝트 구조

## 프로젝트 개요

Kubernetes 클러스터 내에서 발생하는 Event 리소스를 실시간으로 수집하여 Kafka, Elasticsearch, 또는 Local Volume(파일)으로 저장하는 경량 이벤트 수집기입니다.

- Kubernetes `client-go`의 **Informer 패턴** 기반 (watch 스트리밍, polling 없음)
- 수집된 이벤트를 공통 모델로 변환 후 복수의 Exporter로 병렬 전달
- Go 고루틴/채널 기반 비동기 처리, graceful shutdown 지원

---

## 디렉터리 구조

```
event-collector/
├── cmd/
│   └── collector/
│       ├── main.go                        # 진입점: 환경 변수로 logger 초기화, 설정 로드, Collector 실행
│       └── main_test.go                   # TestMain: 테스트 환경 setup/teardown
├── internal/
│   ├── app/
│   │   ├── collector.go                   # Collector: Exporter/kube.Client 조립, Run() graceful shutdown
│   │   │                                  # wg.Add(1)을 goroutine 시작 전에 호출하여 race condition 방지
│   │   └── handler.go                     # Handler: Informer 콜백 수신 후 handle()로 Exporter 전달
│   ├── config/
│   │   └── config.go                      # YAML 설정 파일 로드 및 필수 항목 검증 (Config 구조체)
│   ├── exporter/
│   │   ├── exporter.go                    # Exporter 인터페이스 정의 (Component: Start, Exporter: Write)
│   │   ├── elasticsearch/
│   │   │   ├── exporter.go                # ElasticsearchExporter: 채널 + 버퍼 기반 Bulk 인덱싱
│   │   │   │                              # ticker는 flushTime 사용, writeBulk에 ctx + 30s timeout 적용
│   │   │   ├── exporter_test.go
│   │   │   └── option.go                  # ES 옵션 (chanSize, flushTime, flushSize, user, pass)
│   │   ├── kafka/
│   │   │   ├── exporter.go                # KafkaExporter: sarama AsyncProducer 기반 비동기 전송
│   │   │   ├── exporter_test.go
│   │   │   └── option.go                  # Kafka 옵션 (timeout, retry, flush, partitioner, compression)
│   │   └── volume/
│   │       ├── exporter.go                # VolumeExporter: 로컬 파일 기록, 크기/개수 기반 로테이션
│   │       │                              # sync.Mutex로 currentFile/currentCount 동시 접근 보호
│   │       ├── exporter_test.go
│   │       ├── option.go                  # Volume 옵션 (maxFileSize, maxFileCount, chanSize=WithChanSize)
│   │       ├── test/                      # 테스트 실행 시 생성되는 출력 파일 디렉터리
│   │       └── util.go                    # 파일명 숫자 suffix 기준 정렬 유틸리티
│   │                                      # reNumericSuffix 전역 변수로 regexp 1회만 컴파일
│   ├── kube/
│   │   ├── client.go                      # Kubernetes Informer Client (네임스페이스별 또는 전체)
│   │   ├── client_test.go
│   │   ├── event.go                       # v1.Event → 공통 Event 모델 변환 + JSON 직렬화
│   │   └── option.go                      # Client 옵션 (kubeconfig, resyncPeriod, namespaces)
│   ├── logger/
│   │   ├── logger.go                      # CreateGlobalLogger / CreateLogger: 환경 aware 초기화
│   │   ├── logger_test.go
│   │   ├── option.go                      # logger 옵션 (level, size, age, backup, compress, encoder)
│   │   ├── test/                          # 테스트 실행 시 생성되는 로그 파일 디렉터리
│   │   └── writer.go                      # 전역 logger 래퍼 (Debug/Info/Warn/Error/Fatal/Panic)
│   ├── pprof/
│   │   ├── pprof_dev.go                   # 개발용: pprof HTTP 서버 0.0.0.0:6060 활성화 (빌드 태그: !prd)
│   │   └── pprof_prd.go                   # 운영용: no-op (빌드 태그: prd)
│   └── testutil/
│       └── dummy.go                       # 테스트용 더미 kube.Event 생성기 (MakeDummy)
├── bootstrap/
│   ├── helm/
│   │   ├── elasticsearch/                 # Elasticsearch Helm chart
│   │   ├── strimzi-kafka-operator/        # Strimzi Kafka Operator Helm chart
│   │   ├── elasticsearch.yaml             # Elasticsearch Helm values
│   │   ├── strimzi.yaml                   # Strimzi Helm values
│   │   ├── kafka_cluster.yaml             # Kafka 클러스터 CR
│   │   ├── kafka_topic.yaml               # Kafka 토픽 CR
│   │   ├── gh-runner.yaml                 # GitHub Actions runner
│   │   ├── elasticsearch_ilm.sh           # ILM 정책 설정 스크립트
│   │   ├── elasticsearch_index.sh         # index 생성 스크립트
│   │   └── elasticsearch_template.sh      # index template 설정 스크립트
│   └── script/
│       ├── elasticsearch_bulk.sh          # Elasticsearch bulk 데이터 삽입 스크립트
│       ├── elasticsearch_delete_index.sh  # Elasticsearch index 삭제 스크립트
│       ├── kcat.yaml                      # kcat Pod manifest
│       ├── kcat_consume.sh                # Kafka 메시지 소비 스크립트
│       ├── kcat_describe.sh               # Kafka topic 정보 조회 스크립트
│       └── kcat_produce.sh                # Kafka 메시지 발행 스크립트
├── manifest/
│   ├── collector-deploy.yaml              # Collector Deployment manifest
│   ├── collector-rbac.yaml                # Collector RBAC (ServiceAccount, ClusterRole 등)
│   └── ghar_rbac.yaml                     # GitHub Actions runner RBAC
├── docs/
│   ├── event_collector.drawio.png         # 아키텍처 다이어그램
│   └── go-test-workflows.drawio.png       # CI 워크플로우 다이어그램
├── .github/
│   └── workflows/
│       └── go-test.yaml                   # Go 테스트 CI (vet, race, timeout 포함)
├── Dockerfile                             # 멀티스테이지 빌드 (golang:1.26-alpine → alpine:3.21)
├── Makefile                               # make build: Docker 이미지 빌드
├── go.mod
├── go.sum
├── OLD_README.md                          # 이전 README 보관본
├── CLAUDE.md                              # 코드 작성 규칙
├── CLAUDE_PLAN.md                         # 리팩토링 계획 및 체크리스트
└── CLAUDE_STRUCT.md                       # 프로젝트 구조 문서 (현재 파일)
```

---

## 핵심 데이터 흐름

```
Kubernetes API Server
        │  (watch 스트리밍)
        ▼
  kube.Client (Informer)
        │  OnAdd / OnUpdate / OnDelete
        ▼
  app.Handler.handle()
        │  kube.ConvertBytes() → JSON []byte
        │
        ├──▶ KafkaExporter.Write()         → sarama AsyncProducer → Kafka
        ├──▶ ElasticsearchExporter.Write() → 채널 → 버퍼 → Bulk API → Elasticsearch
        └──▶ VolumeExporter.Write()        → 채널 → 파일 기록 (로테이션)
```

---

## 주요 의존성

| 패키지 | 용도 |
|--------|------|
| `k8s.io/client-go` | Kubernetes Informer, Clientset |
| `k8s.io/api` | Kubernetes Event v1 타입 |
| `github.com/IBM/sarama` | Kafka 비동기 프로듀서 |
| `github.com/elastic/go-elasticsearch/v8` | Elasticsearch Bulk API |
| `go.uber.org/zap` | 구조화 로깅 |
| `gopkg.in/natefinch/lumberjack.v2` | 운영용 로그 파일 로테이션 |
| `gopkg.in/yaml.v3` | 설정 파일 파싱 |

---

## 설정 구조 (config.yaml)

```yaml
kube:
  config: ""           # kubeconfig 경로 (없으면 in-cluster)
  resync: 0            # resync 주기 (0이면 비활성화)
  namespaces: []       # 수집 대상 네임스페이스 (없으면 전체)

kafka:
  enable: false
  brokers: []          # 필수
  topic: ""            # 필수
  timeout: 0
  retry: 0
  retryBackoff: 0
  flushMsg: 0
  flushTime: 0
  flushByte: 0

elasticsearch:
  enable: false
  addresses: []        # 필수
  index: ""            # 필수
  chanSize: 0
  flushTime: 0
  flushSize: 0

volume:
  enable: false
  fileName: ""         # 필수
  filePath: ""         # 필수
  maxFileSize: 0
  maxFileCount: 0
  # chanSize는 config.yaml 미지원, WithChanSize 옵션으로만 설정 가능
```

---

## 환경 변수

| 환경 변수 | 설명 |
|-----------|------|
| `LOG_LEVEL` | 최소 log 레벨 (DEBUG, INFO, WARN, ERROR 등) |
| `LOG_SIZE` | log 파일 최대 크기 (MB) |
| `LOG_AGE` | log 파일 보관 기간 (일) |
| `LOG_BACK` | 보관할 이전 log 파일 최대 개수 |
| `LOG_COMPRESS` | log 파일 gzip 압축 여부 (true/false) |
| `APP_ENV` | `dev` 설정 시 no-op logger 사용 |
| `ELASTICSEARCH_USER` | Elasticsearch 인증 사용자명 |
| `ELASTICSEARCH_PASSWORD` | Elasticsearch 인증 비밀번호 |

---

## 빌드

```bash
# Docker 이미지 빌드
make build

# 태그 지정 빌드
make build IMAGE_TAG=v1.0.0

# Go 바이너리 빌드 (로컬)
go build ./cmd/collector

# 운영 빌드 (prd 태그)
go build -tags prd ./cmd/collector
```
