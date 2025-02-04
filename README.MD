### Запуск сервисов

1. Создать namespace
```shell
kubectl create namespace otus-hw9
```
2. Установить RabbitMQ
```shell
helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq -f ./kuber/helm/rabbitmq_values.yaml --namespace otus-hw9
```
3. Установить PostgreSQL
```shell
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql -f ./kuber/helm/postgres_values.yaml --namespace otus-hw9
```
4. Установить Redis
```shell
helm install redis oci://registry-1.docker.io/bitnamicharts/redis -f ./kuber/helm/redis_values.yaml --namespace otus-hw9
```
5. Применить манифесты
```shell
kubectl apply -f ./kuber
```

## Примечание 
Ключ идемпотентности действует в течение 10 секунд. В тестах между вторым и третьим созданием заказа добавлена задержка в 10 секунд, 
чтобы можно было повторно создать заказ с теми же параметрами.