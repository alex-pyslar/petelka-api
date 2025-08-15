echo "Загрузка последней версии backend'а из репозитория"
git pull
echo "Удаление старого контейнера"
docker rm -f petelka-api
echo "Удаление старого образа"
docker image rm -f petelka-api:latest
echo "Создание нового образа"
docker build -t petelka-api:latest .
echo "Запуск контейнера из нового образа"
docker run --name petelka-api -d -p 8080:8080 petelka-api:latest
docker update --restart=always petelka-api
echo "Deploy завершён"