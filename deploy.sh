#!/bin/bash

set -e

echo "Загрузка последней версии бэкенда из репозитория"
git fetch
git checkout main
git pull origin main

echo "Удаление старого образа (если существует)"
docker image rm -f petelka-api:latest || true

echo "Создание нового образа"
docker build -t petelka-api:latest .

echo "Экспорт образа в k8s"
docker save "petelka-api:latest" > "petelka-api.tar"
sudo ctr -n k8s.io images import "petelka-api.tar"
rm "petelka-api.tar"

echo "Применение манифестов Kubernetes"
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

echo "Перезапуск Deployment для применения нового образа"
kubectl rollout restart deployment petelka-api