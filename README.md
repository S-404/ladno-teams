# Ladno Teams

Небольшой сервер для обмена workspace с командами.

## Стек

- Go 1.26
- PostgreSQL
- sqlx + pgx
- JWT (refresh cookie + access header)
- REST API
- HTMX admin panel
- Swagger (docs/rest/swagger/swagger.yaml)

## Запуск

1. Создайте `.env` по примеру:

```bash
cp .env.example .env
```

2. Убедитесь, что PostgreSQL доступен и создана база `ladno_teams` (или указана в `DB_NAME`).

3. Примените миграции:

```bash
goose up
```

Переменные для goose задаются в `.env`:

```env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}
GOOSE_MIGRATION_DIR=./infrastructure/database/migrations
```

4. Запустите сервер:

```bash
go run ./cmd
```

## Сайт (вход и регистрация по инвайту)

- `http://localhost:8080/` — главная
- `http://localhost:8080/login` — вход
- `http://localhost:8080/register?guid=<invite-guid>` — регистрация по инвайту

Регистрация без инвайта недоступна. Клиент использует `POST /api/auth/invite-register` с полем `guid` инвайта.

После входа или регистрации на сайте refresh token сохраняется в httpOnly cookie `token` (как и при API-логине).

## Админ-панель

Откройте в браузере `http://localhost:8080/admin` — при отсутствии сессии произойдёт редирект на страницу входа `http://localhost:8080/admin/login`.

Учётные данные по умолчанию (из `.env`):

- login: `admin`
- password: `admin`

Первый админ создается автоматически на основе переменных `ADMIN_LOGIN`/`ADMIN_PASSWORD`/`ADMIN_NAME`.

## API

- `POST /api/auth/login` — логин (возвращает access token, refresh token в httpOnly cookie)
- `POST /api/auth/logout` — выход
- `POST /api/auth/refresh` — обновление access token
- `POST /api/auth/invite-register` — регистрация по инвайту
- `POST /api/auth/accept-invite` — принятие инвайта авторизованным пользователем
- `POST /api/teams` — создать команду
- `GET /api/teams` — мои команды
- `GET|PUT|DELETE /api/teams/:guid` — управление командой (только лидер)
- `POST /api/teams/:guid/invites` — создать инвайт (лидер)
- `GET /api/teams/:guid/teammates` — список участников
- `PUT|DELETE /api/teams/:guid/teammates/:user_guid` — управление участниками (лидер)
- `POST /api/teams/:guid/workspaces` — создать workspace (лидер)
- `GET /api/teams/:guid/workspaces` — список workspace в команде
- `GET|PUT|DELETE /api/workspaces/:guid` — доступ к workspace по ролям
- `POST /api/workspaces/:guid/roles` — назначить роль в workspace (лидер/maintainer)
- `GET /api/workspaces/:guid/roles` — список ролей
- `PUT|DELETE /api/workspaces/:guid/roles/:teammate_guid` — управление ролями

## Роли в workspace

- `MAINTAINER` — полное управление workspace
- `DEVELOPER` — обновление workspace
- `GUEST` — только чтение

Лидеры команды автоматически имеют права владельца на все workspace своей команды.
