echo "Загрузка последней версии backend'а из репозитория"
git pull
echo "Удаление старого контейнера"
docker rm -f online-store
echo "Удаление старого образа"
docker image rm -f online-store:latest
echo "Создание нового образа"
docker build -t online-store:latest .
echo "Запуск контейнера из нового образа"
docker run --name online-store -d -p 8080:8080 online-store:latest
docker update --restart=always online-store
echo "Deploy завершён"