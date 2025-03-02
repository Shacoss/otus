### Запуск сервисов

1. Создать namespace
```shell
kubectl create namespace otus-hw10
```
2. Установить RabbitMQ
```shell
helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq -f ./kuber/helm/rabbitmq_values.yaml --namespace otus-hw10
```
3. Установить PostgreSQL
```shell
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql -f ./kuber/helm/postgres_values.yaml --namespace otus-hw10
```
4. Установить Redis
```shell
helm install redis oci://registry-1.docker.io/bitnamicharts/redis -f ./kuber/helm/redis_values.yaml --namespace otus-hw10
```
5. Применить манифесты
```shell
kubectl apply -f ./kuber
```
6. Установить Prometheus
```shell
helm install prometheus oci://registry-1.docker.io/bitnamicharts/prometheus -f ./kuber/helm/prometheus_values.yaml --namespace otus-hw10
```
7. Установить Grafana
```shell
helm install grafana oci://registry-1.docker.io/bitnamicharts/grafana -f ./kuber/helm/grafana_values.yaml --namespace otus-hw10
```