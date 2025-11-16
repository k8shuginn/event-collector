# 테스트 개발환경 구축하기
개발 및 테스트를 위한 쿠버네티스 클러스터에 필요한 컴포넌트들을 설치하는 방법입니다. event-collector 시스템은 Elasticsearch와 Kafka를 사용하므로, 이 두 가지 컴포넌트를 설치하는 방법을 안내합니다. 또한, GitHub Actions의 self-hosted runner를 설정하여 TDD 환경을 구축하여 안전한 개발이 가능하도록 합니다.

# elasticsearch cluster install
개인 로컬 Kubernetes 클러스터 환경에 작은 크기의 Elasticsearch를 3노드로 구성합니다.
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


# Kafka cluster install (using Strimzi)
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

# github action self-hosted runner service account permission setup
```bash
NAMESPACE="gh-runner"
kubectl create namespace ${NAMESPACE}
kubectl apply -f manifest/gh_rbac.yaml
```

# github action self-hosted runner setup
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

