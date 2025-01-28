### Запуск сервисов

1. Создать namespace
```shell
kubectl create namespace otus-hw8
```
2. Установить RabbitMQ
```shell
helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq -f ./kuber/helm/rabbitmq_values.yaml --namespace otus-hw8
```
3. Установить PostgreSQL
```shell
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql -f ./kuber/helm/postgres_values.yaml --namespace otus-hw8
```
4. Применить манифесты
```shell
kubectl apply -f ./kuber
```

## Примечание 
Дата доставки должна быть не раньше текущей даты плюс один день и не позже одного месяца с текущей даты.