# HiTalent Test Task - API для вопросов и ответов

API-сервис для управления вопросами и ответами, построенный на Go с использованием GORM, PostgreSQL и Docker.

## Технологии

- **Go** 1.21+
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - реляционная база данных
- **Goose** - миграции базы данных
- **Chi** - HTTP роутер
- **Zap** - структурированное логирование
- **Docker** - контейнеризация

## Архитектура

Проект следует Clean Architecture принципам:

```
backend/
├── cmd/main.go              # Точка входа
├── config/                  # Конфигурация
├── internal/
│   ├── entity/             # Сущности домена
│   ├── port/               # Интерфейсы (порты)
│   │   ├── repo/           # Интерфейсы репозиториев
│   │   └── service/        # Интерфейсы сервисов
│   ├── cases/              # Бизнес-логика (use cases)
│   ├── adapter/            # Адаптеры
│   │   └── repo/           # Реализация репозиториев
│   └── input/              # Входные точки
│       └── http/           # HTTP handlers
└── pkg/
    └── migration/          # Миграции базы данных
```

## API Endpoints

### Вопросы (Questions)

- `GET /questions/` - получить список всех вопросов
- `POST /questions/` - создать новый вопрос
- `GET /questions/{id}` - получить вопрос и все ответы на него
- `DELETE /questions/{id}` - удалить вопрос (вместе с ответами)

### Ответы (Answers)

- `POST /questions/{id}/answers/` - добавить ответ к вопросу
- `GET /answers/{id}` - получить конкретный ответ
- `DELETE /answers/{id}` - удалить ответ

## Запуск с помощью Docker

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd HiTalent_TestTask
```

2. Запустите приложение через docker-compose:
```bash
docker-compose up --build
```

Приложение будет доступно по адресу `http://localhost:8080`

## Запуск локально

1. Убедитесь, что у вас установлены Go 1.21+ и PostgreSQL

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` в корне проекта:
```env
POSTGRES_CONNECTION_STRING=host=localhost user=your_user password=your_password dbname=your_db sslmode=disable port=5432
HTTP_PORT=8080
```

4. Запустите миграции (они применяются автоматически при старте приложения)

5. Запустите приложение:
```bash
go run backend/cmd/main.go
```

## Примеры использования API

### Создать вопрос
```bash
curl -X POST http://localhost:8080/questions/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Что такое Go?"}'
```

### Получить список вопросов
```bash
curl http://localhost:8080/questions/
```

### Получить вопрос с ответами
```bash
curl http://localhost:8080/questions/1
```

### Добавить ответ к вопросу
```bash
curl -X POST http://localhost:8080/questions/1/answers \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "text": "Go - это язык программирования"}'
```

### Получить ответ
```bash
curl http://localhost:8080/answers/1
```

### Удалить вопрос
```bash
curl -X DELETE http://localhost:8080/questions/1
```

### Удалить ответ
```bash
curl -X DELETE http://localhost:8080/answers/1
```

## Особенности реализации

- Каскадное удаление: при удалении вопроса автоматически удаляются все его ответы
- Валидация: нельзя создать ответ к несуществующему вопросу
- Множественные ответы: один пользователь может оставлять несколько ответов на один вопрос
- Структурированное логирование с использованием Zap
- Автоматические миграции при запуске приложения

## Структура базы данных

### Таблица `questions`
- `id` - первичный ключ (SERIAL)
- `text` - текст вопроса (TEXT, NOT NULL)
- `created_at` - время создания (TIMESTAMP, DEFAULT NOW())

### Таблица `answers`
- `id` - первичный ключ (SERIAL)
- `question_id` - внешний ключ на questions (INTEGER, NOT NULL)
- `user_id` - идентификатор пользователя (VARCHAR(255), NOT NULL)
- `text` - текст ответа (TEXT, NOT NULL)
- `created_at` - время создания (TIMESTAMP, DEFAULT NOW())

## Тестирование

Для тестирования API можно использовать:
- `curl` (примеры выше)
- Postman
- HTTPie
- Любой другой HTTP клиент

## Логирование

Приложение использует структурированное логирование через Zap. Все операции (создание, чтение, удаление) логируются с соответствующим уровнем детализации.

