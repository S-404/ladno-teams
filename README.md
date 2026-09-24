# Ladno Teams

A small server for sharing workspaces with teams.

## Stack

- Go 1.26
- PostgreSQL
- sqlx + pgx
- JWT (refresh cookie + access header)
- REST API
- HTMX admin panel
- Swagger (docs/rest/swagger/swagger.yaml)

## Getting started

1. Create `.env` from the example:

```bash
cp .env.example .env
```

2. Make sure PostgreSQL is running and the `ladno_teams` database exists (or set `DB_NAME` accordingly).

3. Run migrations:

```bash
goose up
```

Goose variables are defined in `.env`:

```env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}
GOOSE_MIGRATION_DIR=./infrastructure/database/migrations
```

4. Start the server:

```bash
go run ./cmd
```

## Website (sign-in and invite-only registration)

- `http://localhost:8080/` — home
- `http://localhost:8080/login` — sign in
- `http://localhost:8080/register?guid=<invite-guid>` — register via invite

Registration without an invite is not available. The client uses `POST /api/auth/invite-register` with the invite `guid` field.

After sign-in or registration on the website, the refresh token is stored in the httpOnly `token` cookie (same as API login).

## Admin panel

Open `http://localhost:8080/admin` in your browser — if there is no session, you will be redirected to `http://localhost:8080/admin/login`.

Default credentials (from `.env`):

- login: `admin`
- password: `admin`

The first admin is created automatically from `ADMIN_LOGIN`, `ADMIN_PASSWORD`, and `ADMIN_NAME`.

## API

- `POST /api/auth/login` — sign in (returns access token, refresh token in httpOnly cookie)
- `POST /api/auth/logout` — sign out
- `POST /api/auth/refresh` — refresh access token
- `POST /api/auth/invite-register` — register via invite
- `POST /api/auth/accept-invite` — accept invite as an authenticated user
- `POST /api/teams` — create a team
- `GET /api/teams` — my teams
- `GET|PUT|DELETE /api/teams/:guid` — manage team (leader only)
- `POST /api/teams/:guid/invites` — create invite (leader)
- `GET /api/teams/:guid/teammates` — list teammates
- `PUT|DELETE /api/teams/:guid/teammates/:user_guid` — manage teammates (leader)
- `POST /api/teams/:guid/workspaces` — create workspace (leader)
- `GET /api/teams/:guid/workspaces` — list workspaces in team
- `GET|PUT|DELETE /api/workspaces/:guid` — workspace access by role
- `POST /api/workspaces/:guid/roles` — assign workspace role (leader/maintainer)
- `GET /api/workspaces/:guid/roles` — list roles
- `PUT|DELETE /api/workspaces/:guid/roles/:teammate_guid` — manage roles

## Workspace git sync

Each team workspace has a bare git repository on disk. The Fyne client clones it and uses pull/push instead of overwriting JSONB snapshots.

- Config: `GIT_REPOS_ROOT` (default `./data/git-repos`) — see `.env.example`
- Path layout: `{GIT_REPOS_ROOT}/{workspace-guid}.git`
- Git HTTP (JWT Bearer, same as REST):
  - `GET /git/workspaces/:guid/info/refs?service=git-upload-pack|git-receive-pack`
  - `POST /git/workspaces/:guid/git-upload-pack`
  - `POST /git/workspaces/:guid/git-receive-pack`
- Roles: `GUEST` can fetch; `DEVELOPER`+ can push. Team leaders have full access.

Requires a system `git` binary on the server PATH (smart HTTP uses `git --stateless-rpc`).

## Workspace roles

- `MAINTAINER` — full workspace management
- `DEVELOPER` — update workspace
- `GUEST` — read-only

Team leaders automatically have owner permissions on all workspaces in their team.
