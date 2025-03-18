kafka 설치
```bash
# Operator 설치
kubectl create ns mytest
helm repo add strimzi https://strimzi.io/charts/
helm repo update
helm upgrade --install -n mytest kafka-operator strimzi/strimzi-kafka-operator --version 0.45.0 -f ./middleware/strimzi.yaml

# Kafka Cluster 생성
kubectl apply -f ./middleware/kafka_cluster.yaml
# Kafka Topic 생성
kubectl apply -f ./middleware/kafka_topic.yaml

# kcat 사용
kubectl apply -f ./test/kcat.yaml
./test/kcat_consume.sh
```

elasticsearch 설치
```bash
# ElasticSearch 설치
helm repo add elastic https://helm.elastic.co
helm repo update
helm upgrade --install --version 8.5.1 -n mytest elasticsearch elastic/elasticsearch -f ./middleware/elasticsearch.yaml

# elasticsearch ilm 설정
./middleware/elasticsearch_ilm.sh
./middleware/elasticsearch_template.sh
./middleware/elasticsearch_index.sh
```

