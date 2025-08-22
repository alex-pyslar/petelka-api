# Petelka API

Petelka API - это RESTful API для управления пользователями, продуктами, категориями, заказами и комментариями в интернет-магазине. Проект построен на языке Go с использованием современных инструментов и библиотек для обеспечения высокой производительности и масштабируемости. API развернут и доступен по адресу [https://api.petelka.velesoft.ru](https://api.petelka.velesoft.ru), а интерактивная документация доступна через Swagger UI по адресу [https://api.petelka.velesoft.ru/swagger/index.html](https://api.petelka.velesoft.ru/swagger/index.html).

## Технологии

- **Язык программирования**: Go
- **Фреймворк маршрутизации**: [Gorilla Mux](https://github.com/gorilla/mux)
- **Базы данных**:
  - PostgreSQL (основное хранилище)
  - Redis (кэширование)
- **Логирование**: Собственный логгер
- **Мониторинг**: Prometheus
- **Документация API**: Swagger (OpenAPI)
- **Аутентификация**: JWT

## Основные возможности

- Регистрация и авторизация пользователей
- Управление продуктами (CRUD операции)
- Управление категориями (CRUD операции)
- Создание и управление заказами
- Создание комментариев к продуктам
- Поиск продуктов
- Разграничение доступа (публичные, защищенные и административные маршруты)

## Установка

### Предварительные требования

- Go (версия 1.16 или выше)
- PostgreSQL
- Redis
- Git

### Шаги установки

1. Клонируйте репозиторий:
```bash
git clone https://github.com/alex-pyslar/petelka-api.git
cd petelka-api
```

2. Установите зависимости:
```bash
go mod download
```

3. Настройте переменные окружения:
Создайте файл `.env` в корне проекта со следующими переменными:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=petelka_db
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=your_jwt_secret
```

## Запуск проекта

1. Скомпилируйте и запустите сервер:
```bash
go run main.go
```

2. Локально сервер будет доступен по адресу `http://localhost:8080`. Для доступа к развернутому API используйте [https://api.petelka.velesoft.ru](https://api.petelka.velesoft.ru).

## Документация API

Документация API доступна через Swagger UI по адресу: [https://api.petelka.velesoft.ru/swagger/index.html](https://api.petelka.velesoft.ru/swagger/index.html). Swagger предоставляет интерактивный интерфейс для тестирования всех доступных эндпоинтов API.

## Эндпоинты

### Публичные маршруты
- `POST /api/auth/register` - Регистрация нового пользователя
- `POST /api/auth/login` - Авторизация пользователя
- `GET /api/products` - Список всех продуктов
- `GET /api/products/search` - Поиск продуктов
- `GET /api/products/{id}` - Получение информации о продукте
- `GET /api/categories` - Список всех категорий
- `GET /api/categories/{id}` - Получение информации о категории

### Защищенные маршруты (требуется авторизация)
- `POST /api/comments` - Создание комментария
- `POST /api/orders` - Создание заказа

### Административные маршруты (требуется роль администратора)
- `POST /api/products` - Создание продукта
- `PUT /api/products/{id}` - Обновление продукта
- `DELETE /api/products/{id}` - Удаление продукта
- `POST /api/categories` - Создание категории
- `PUT /api/categories/{id}` - Обновление категории
- `DELETE /api/categories/{id}` - Удаление категории
- `GET /api/users` - Список всех пользователей
- `PUT /api/users/{id}` - Обновление пользователя
- `DELETE /api/users/{id}` - Удаление пользователя

## Мониторинг

Метрики Prometheus доступны по адресу:
```
http://localhost:8080/metrics
```
Для развернутого сервера: [https://api.petelka.velesoft.ru/metrics](https://api.petelka.velesoft.ru/metrics)

## Разработка

Для генерации Swagger документации используйте:
```bash
swag init
```

## Лицензия

[MIT License](LICENSE)