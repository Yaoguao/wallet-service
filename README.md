# wallet-service


## Быстрый старт

1. Склонируйте репозиторий:

```bash
git clone https://github.com/Yaoguao/wallet-service.git
cd wallet-service
```

2. Опционально. Скопируйте файл окружения или отредактируйте значения (см. `config.env` и `.env`):

```bash
cp config.env .env
# отредактируйте .env под вашу среду
```

3. Запустите через Docker Compose (локально поднимается БД и сервисы, если настроено):

```bash
docker-compose up --build
```

---

## Что внутри репозитория

* `cmd/wallet-service` — точка входа сервиса.
* `internal/` — бизнес-логика, репозитории, сервисы.
* `pkg/` — утилитарные пакеты, которые могут быть для работы с postgres с использованием pgx, ну просто обертка.
* `migrations/` — SQL-миграции для PostgreSQL.
* `config.env`, `.env` — примеры переменных окружения.
* `docker-compose.yml` и `Dockerfile` — для контейнеризации.
* `tests/` — тесты.

---


## Миграции

Миграции хранятся в `migrations/`. Для применения миграций используйте предпочитаемый вам инструмент (например, `go-migrate`):

```bash
migrate -path migrations -database "postgres://user:pass@host:5432/dbname?sslmode=disable" up
```

---

## Архитектура и рекомендации

* Проект организован по стандартной структуре Go: `cmd/`, `internal/`, `pkg/`.
* Репозитории и слои сервиса разделены: обработчики HTTP -> сервисный слой -> репозитории -> БД.
* Используется PostgreSQL (см. `migrations/`) и Docker для локальной разработки.


## Contributing

PR и issue приветствуются. Оставьте в описании PR краткое резюме изменений или рекомендаций.

---

