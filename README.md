# event-collector
event-collector는 Kubernetes 클러스터 내에서 발생하는 다양한 이벤트를 실시간으로 수집하고, 이를 장기 보관과 분석이 가능한 형태로 Kafka, Elasticsearch, 혹은 Local Volume 파일 형태로 저장하는 목적으로 제작된 경량 수집기(Collector)입니다. Kubernetes는 클러스터 내부에서 발생하는 상태 변화에 대해 Event라는 특별한 리소스를 통해 중요한 정보를 제공합니다. 예를 들어, Pod가 스케줄링에 실패하거나(Node 부족, Taint 적용 등), 컨테이너가 CrashLoopBackOff 상태에 빠지거나, Kubelet이 노드를 Ready 상태로 유지하지 못하는 경우처럼 운영 과정에서 필수적으로 알아야 할 경고 또는 장애 정보를 Event 리소스를 통해 노출합니다.

그러나 Kubernetes의 Event는 기본적으로 수명이 짧고 TTL이 지나면 삭제되는 일시적 저장 구조이기 때문에, 실제 운영 환경에서 장애를 분석하거나 트렌드를 장기적으로 관찰하기에는 적합하지 않습니다. 이벤트가 사라지기 전에 빠르게 대응하지 않으면 문제 원인을 파악하는 데 어려움이 생기며, 장애 발생의 흐름을 재현하거나 클러스터 사용 패턴을 분석하는 것도 매우 어려워집니다. 이렇게 중요한 정보를 지속적으로 수집해 체계적으로 보존하고 분석할 수 있는 별도의 모듈이 필요하다는 점에서 이 프로젝트는 출발했습니다.

event-collector는 Kubernetes의 client-go 라이브러리에서 제공하는 Informer 패턴을 기반으로 동작하며, API Server를 지속적으로 폴링(polling)하는 방식이 아닌 watch 기반의 스트리밍 방식으로 이벤트를 수집합니다. 이 구조를 통해 클러스터에 부하를 주지 않으면서도 이벤트가 생성된 시점에 거의 실시간에 가깝게 이벤트를 받아볼 수 있습니다. Informer는 내부적으로 로컬 캐시를 유지하며, Add / Update / Delete 이벤트가 발생할 때마다 등록된 핸들러를 호출하도록 되어 있어, 이벤트 스트림을 안정적으로 수집하고 후처리할 수 있는 구조를 자연스럽게 제공합니다.

본 프로젝트에서는 Informer를 통해 들어온 Event 데이터를 공통 모델로 변환한 뒤, 다양한 Exporter로 전달하여 저장할 수 있도록 설계했습니다. Exporter는 인터페이스 기반으로 설계되어 있어 사용자는 원하는 저장소(Kafka, Elasticsearch, Local File)를 선택하거나, 필요하다면 새로운 저장소를 위한 Exporter를 직접 개발하여 추가할 수도 있습니다.

event-collector의 전체 구조는 매우 경량이며, Go의 고루틴 및 채널 기반 비동기 방식으로 동작하기 때문에 높은 처리량을 요구하지 않는 환경에서도 최소한의 자원으로 안정적으로 운영될 수 있습니다. 실제로 CPU 사용률은 매우 낮게 유지되고, 메모리 역시 수십 MiB 정도면 충분합니다. 이는 클러스터 내부에서 동작하는 애플리케이션 특성상 매우 중요한 요소입니다. Collector가 과도한 리소스를 사용하면 오히려 클러스터의 안정성에 악영향을 미칠 수 있기 때문입니다.

# event-collector 실행 조건
event-collector가 Kubernetes Event 리소스를 안정적으로 수집하기 위해서는 클러스터 내부 또는 외부에서 Kubernetes API Server에 접근할 수 있는 적절한 인증·인가 구성이 필요합니다. Collector는 두 가지 방식 중 하나를 통해 클러스터와 연결됩니다. 첫 번째는 로컬 환경이나 외부 시스템에서 실행될 때 사용되는 kubeconfig 기반 인증 방식이며, 두 번째는 Pod 형태로 클러스터 내부에서 실행될 때 자동으로 주어지는 in-cluster config 방식입니다.

또한 Collector를 Pod 형태로 배포하는 경우, 이 ClusterRole을 특정 ServiceAccount와 묶어주기 위해 적절한 RoleBinding 또는 ClusterRoleBinding이 필요합니다. 이를 통해 Collector는 지정된 네임스페이스나 클러스터 전역에서 Event 리소스에 접근할 수 있는 공식적인 권한을 부여받게 됩니다. 운영 환경에서는 보안 강화를 위해 불필요한 권한을 포함하지 않도록 최소 권한(Least Privilege) 원칙에 따라 ClusterRole을 설계하는 것이 중요합니다. 예를 들어, Event 리소스 외의 다른 리소스(Pod, Deployment 등)에 대한 권한은 Collector가 필요로 하지 않으므로 부여하지 않는 것이 바람직합니다.

# Exporter 종류
## Kafka Exporter
Kafka Exporter는 수집된 이벤트 데이터를 Apache Kafka로 전송하는 역할을 합니다.
Kafka Exporter는 대규모 이벤트 처리나 스트리밍 파이프라인을 구성하는 데 적합합니다. Kubernetes 이벤트는 정상적인 클러스터에서는 매우 대량으로 발생하지 않지만, 규모가 큰 환경(예: 수백 개 이상의 노드, 수천 개의 Pod)에서는 이벤트가 짧은 시간 내에 폭발적으로 증가할 수 있습니다. Kafka Exporter는 이러한 상황에서도 안정적으로 이벤트를 전달할 수 있으며, 이후 Flink, Spark, Kafka Streams, ksqlDB와 같은 분석 시스템과 연동하여 실시간 분석 환경을 쉽게 구성할 수 있습니다.

## Elasticsearch Exporter
Elasticsearch Exporter는 운영 모니터링과 장애 분석을 좀 더 직관적으로 수행하기 위한 용도로 설계되었습니다. 이벤트를 Elasticsearch에 저장하면 Kibana를 통해 문제 시점별로 이벤트를 검색하고 필터링할 수 있으며, 특정 워크로드 또는 특정 노드에서 반복적으로 발생하는 문제를 한눈에 확인할 수 있습니다. 예를 들어, 특정 Deployment가 지속적으로 CrashLoopBackOff 상태에 빠진다거나 특정 노드에서만 ImagePullBackOff가 반복된다면 Elasticsearch 분석 환경에서 이를 손쉽게 시각화하고 패턴을 찾을 수 있습니다.

## Volume Exporter
Volume Exporter는 이벤트 데이터를 로컬 파일 시스템에 저장하는 간단한 방법을 제공합니다. 개발 환경에서의 테스트나 간단한 Raw 데이터 보관용으로 사용할 수 있으며, AWS S3 같은 오브젝트 스토리지에 업로드하기 전에 임시로 데이터를 저장하는 용도로도 활용할 수 있습니다. Volume Exporter는 이벤트를 손실 없이 모두 파일로 저장하고자 할 때 유용하며, 이벤트 원본을 그대로 유지하면서 후처리를 진행해야 하는 경우에도 적합합니다. 특히 개발 환경에서는 Kafka나 Elasticsearch 같은 인프라를 구성하기 어렵기 때문에, 단순히 파일 형태로 저장할 수 있는 기능이 큰 도움이 됩니다.

# Kubernetes event 수집방법
해당 프로젝트는 Kubernetes API Server에 Informer 방식을 사용하여 이벤트를 실시간으로 수집합니다. Informer는 Kubernetes 클라이언트 라이브러리에서 제공하는 기능으로, 특정 리소스(이 경우에는 이벤트 리소스)의 변경 사항을 감지하고 이를 처리할 수 있도록 도와줍니다.

Informer는 다음과 같은 방식으로 작동합니다:
1. 리소스 감시: Informer는 Kubernetes API Server에 특정 리소스(이 경우에는 이벤트 리소스)를 감시하도록 설정됩니다. 이를 통해 해당 리소스에 대한 생성, 업데이트, 삭제 등의 변경 사항을 실시간으로 감지할 수 있습니다.
2. 캐싱: Informer는 감시하는 리소스의 상태를 로컬 캐시에 저장합니다. 이를 통해 API Server에 대한 반복적인 요청을 줄이고, 빠른 액세스를 제공합니다.
3. 이벤트 핸들링: Informer는 리소스의 변경 사항이 감지되면 등록된 핸들러 함수를 호출합니다. 이 함수를 통해 변경된 이벤트 데이터를 처리하고, 필요한 작업(예: 데이터 저장소로 전송)을 수행할 수 있습니다.
4. 동기화: Informer는 주기적으로 API Server와 동기화하여 로컬 캐시의 상태를 최신 상태로 유지하여 일관성을 보장합니다.


# event-collector 아키텍처
![event-collector architecture](./docs/event_collector.drawio.png)
1. Kubernetes API Server에서 이벤트를 수집합니다.
2. Event Collector는 수집한 이벤트를 공통 모델로 변환합니다.
3. Event Collector는 Kafka, Elasticsearch, Volume에 이벤트를 저장합니다.



# 프로젝트 소스 코드 구조
* cmd : 프로그램 application 소스 코드
  * collector : 이벤트 수집기 소스 코드
* pkg : 프로그램 application에서 사용하는 패키지 소스 코드
  * kube : Kubernetes API Server와 통신하는 패키지 소스 코드
  * logger : 로그를 수집하는 패키지 소스 코드
* exporter : 프로그램 application에서 사용하는 패키지를 외부로 전달하는 패키지 소스 코드
  * kafka : Kafka로 데이터를 전송하는 패키지 소스 코드
  * elasticsearch : Elasticsearch로 데이터를 전송하는 패키지 소스 코드
  * volume : Volume으로 데이터를 전송하는 패키지 소스 코드
* dev_setup : Kafka, Elasticsearch, Volume 등의 데이터 저장소 설치를 위한 Helm 차트 및 매니페스트 소스 코드
* manifest : Kubernetes 배포 파일 

---
# 개발 환경 구축
개발 및 테스트 환경에서 event-collector를 안정적으로 실행하기 위해서는 몇 가지 핵심 컴포넌트를 사전에 구축해야 합니다. 본 프로젝트는 이벤트 데이터를 저장하고 분석하기 위한 백엔드로 Elasticsearch와 Kafka를 활용하므로, 로컬 또는 테스트용 Kubernetes 클러스터에 이 두 가지 시스템을 설치하는 과정이 필요합니다. 이를 통해 수집된 이벤트를 검색, 분석하거나 스트리밍 파이프라인으로 전달하는 기능을 충분히 검증할 수 있습니다.

## Elasticsearch 설치
개인 로컬 Kubernetes 클러스터 환경에 다음과 같은 방법으로 Elasticsearch를 3노드로 구성합니다.
```bash
# elasticsearch 설치
NAMESPACE=event-collector
helm upgrade --install elasticsearch dev_setup/helm/elasticsearch \
    --namespace $NAMESPACE --create-namespace \
    -f dev_setup/helm/elasticsearch.yaml

# 배포 확인
NAMESPACE=event-collector
kubectl get pods -n $NAMESPACE
```

설치된 Elasticsearch에 ILM 정책, 인덱스 템플릿, 초기 인덱스를 생성합니다.
```bash
./dev_setup/helm/elasticsearch_ilm.sh
./dev_setup/helm/elasticsearch_template.sh
./dev_setup/helm/elasticsearch_index.sh
```

## Kafka cluster install (using Strimzi)
Kubernetes 클러스터에 Strimzi Kafka Operator를 사용하여 Kafka 클러스터를 설치합니다.
```bash
# strimzi kafka operator 설치
NAMESPACE=event-collector
helm upgrade --install strimzi dev_setup/helm/strimzi-kafka-operator \
    --namespace $NAMESPACE --create-namespace \
    -f dev_setup/helm/strimzi.yaml

# kafka cluster 및 topic 생성
NAMESPACE=event-collector
kubectl apply -f dev_setup/helm/kafka_cluster.yaml -n $NAMESPACE
kubectl apply -f dev_setup/helm/kafka_topic.yaml -n $NAMESPACE

# 배포 확인
kubectl get pods -n $NAMESPACE
```

# 테스트 환경 구축
저는 실제 데이터를 다루지 않는 단순한 Mock 테스트보다, 실제 컴포넌트들이 동작하는 환경에서 테스트가 수행되는 방식을 선호합니다. 이러한 이유로 event-collector 프로젝트에서는 GitHub Actions의 Self-Hosted Runner를 활용하여 실제 Kubernetes 클러스터에서 이벤트를 수집하고, 그 데이터를 Elasticsearch와 Kafka에 실제로 저장하는 통합 테스트 환경을 구축하였습니다. 이를 위해 Self-Hosted Runner를 Kubernetes 클러스터 내에 직접 배포하였으며, 해당 Runner는 GitHub Actions 워크플로우를 통해 go test를 실행하면서 실제 클러스터 환경에서 이벤트 수집 로직이 올바르게 동작하는지 검증하도록 구성되어 있습니다. 이러한 방식은 단순한 유닛 테스트를 넘어, 운영 환경과 동일한 조건에서 collector의 동작을 검증할 수 있다는 점에서 신뢰성과 안정성을 크게 높여줍니다.

![event-collector architecture](./docs/go-test-workflows.drawio.png)

## github action self-hosted runner service account permission setup
TDD 환경 구축을 위해 GitHub Action Self-Hosted Runner를 설정합니다. 먼저, 필요한 권한을 가진 서비스 계정을 생성합니다.
```bash
NAMESPACE="gh-runner"
kubectl create namespace ${NAMESPACE}
kubectl apply -f manifest/gh_rbac.yaml
```

# github action self-hosted runner setup
GitHub Actions의 self-hosted runner를 Kubernetes 클러스터에 배포합니다.
```bash
# gha-runner-controller 설치
NAMESPACE="gh-runner"
helm upgrade --install gh-arc \
	oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller \
	--namespace "${NAMESPACE}" --create-namespace

# self-hosted runner 배포
NAMESPACE="gh-runner"
helm upgrade --install gh-runner \
    oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set \
    --namespace ${NAMESPACE} --create-namespace \
    -f dev_setup/helm/gh-runner.yaml \
    --set githubConfigSecret.github_token=${GITHUB_TOKEN} \
    --set githubConfigUrl=${GITHUB_REPO}
```

## Github workflow 작성
.github/workflows/go-test.yaml 파일을 생성하여 Go 테스트를 자동화합니다.
[go-test.yaml](https://github.com/k8shuginn/event-collector/blob/develop/.github/workflows/go-test.yaml)
