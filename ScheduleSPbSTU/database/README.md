# database — база данных ScheduleSPbSTU

PostgreSQL 16 в Docker. В этой папке: запуск БД, настройки и SQL-миграции.

## Структура
```
database/
├── docker-compose.yml   запуск PostgreSQL
├── .env.example         пример настроек (копируем в .env)
├── migrations/
│   ├── 001_groups.sql   таблица groups
│   └── 002_users.sql    таблица users (внешний ключ на groups)
└── README.md
```

## Запуск
Нужны Docker Desktop и терминал в папке `database/`.
```bash
cp .env.example .env        # один раз
docker compose up -d        # запустить БД в фоне
docker compose ps           # статус должен быть healthy
```

## Проверка, что таблицы созданы
```bash
docker compose exec db psql -U app -d schedule -c "\dt"
docker compose exec db psql -U app -d schedule -c "\d users"
```
В списке `\dt` должны быть `groups` и `users`.

## Параметры подключения
| Параметр | Значение (по умолчанию) |
|---|---|
| Хост | `localhost` |
| Порт | `5432` (меняется переменной `DB_PORT`) |
| База | `schedule` |
| Пользователь / пароль | `app` / `app` |
| Строка | `postgres://app:app@localhost:5432/schedule?sslmode=disable` |

## Остановка и сброс
```bash
docker compose down         # остановить, данные сохраняются
docker compose down -v      # остановить и УДАЛИТЬ данные (чистая БД)
```

## Важно
- Файлы из `migrations/` выполняются **только при первом создании** БД. Изменили SQL — выполните `docker compose down -v`, затем `docker compose up -d`.
- Миграции выполняются по порядку имён файлов, поэтому нумерация `001`, `002`... обязательна.
- Порт 5432 не должен быть занят (другой запущенный PostgreSQL или корневой `docker-compose.yml` проекта). Либо остановите их, либо поставьте другой `DB_PORT` в `.env`.
- Бэкенд пока использует собственную миграцию `backend/migrations/001_init.sql`. Объединение схем — отдельная задача команды.
