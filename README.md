# auth

Turborepo mit zwei Next.js-Apps und einem Go-Backend – jede App ist ein eigenes Vercel-Projekt.

```
apps/
  web/         Next.js 16 – Landing Page, /login, /register  (Port 3000)
  dashboard/   Next.js 16 – Dashboard unter /dashboard        (Port 3001)
  api/         Go – REST-API unter /api                       (Port 8080)
               + Postgres (lokal via Docker Compose, auf Vercel via Neon)
packages/
  ui/                 geteilte React-Komponenten (Button, Logo)
  typescript-config/  geteilte tsconfigs
```

## Architektur

Die **web**-App ist der öffentliche Einstiegspunkt. Per Rewrites (`apps/web/next.config.ts`) leitet sie

- `/dashboard/*` an die **dashboard**-App weiter ([Next.js Multi-Zones](https://nextjs.org/docs/app/guides/multi-zones)) und
- `/api/*` an die **Go-API**.

Dadurch läuft alles unter einer Domain, und das Session-Cookie (`httpOnly`-JWT, von der API gesetzt) funktioniert ohne CORS oder Cross-Domain-Cookies.

### API-Endpunkte

| Methode | Pfad                 | Beschreibung                       |
| ------- | -------------------- | ---------------------------------- |
| GET     | `/api/health`        | Health Check                       |
| POST    | `/api/auth/register` | `{ name, email, password }`        |
| POST    | `/api/auth/login`    | `{ email, password }`              |
| POST    | `/api/auth/logout`   | löscht das Session-Cookie          |
| GET     | `/api/auth/me`       | aktueller User aus dem Session-JWT |

### Datenbank

User werden in **Postgres** gespeichert (`apps/api/pkg/auth/postgres.go`, Treiber [pgx](https://github.com/jackc/pgx)). Das Schema liegt als SQL in `apps/api/pkg/auth/migrations/` und wird beim Start der API automatisch angewendet (idempotent, `IF NOT EXISTS`). Neue Migrationen einfach als `002_….sql` usw. dazulegen.

Ist `DATABASE_URL` lokal nicht gesetzt, fällt die API auf einen In-Memory-Store zurück (mit Warnung). Auf Vercel ist `DATABASE_URL` Pflicht.

## Lokale Entwicklung

Voraussetzungen: Node.js ≥ 20, pnpm 10, Go ≥ 1.24, Docker

```sh
pnpm install
cp apps/api/.env.example apps/api/.env
pnpm db:up        # startet Postgres (docker compose)
pnpm dev          # startet web, dashboard und api parallel
```

Dann <http://localhost:3000> öffnen. Weitere Befehle:

```sh
pnpm build        # baut alle Apps
pnpm check-types  # tsc + go vet
pnpm test         # go test
pnpm db:down      # stoppt Postgres (Daten bleiben im Volume)
```

Die Postgres-Tests laufen nur, wenn `TEST_DATABASE_URL` gesetzt ist:

```sh
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable pnpm test
```

## Deployment auf Vercel

Das Repo wird als **drei Vercel-Projekte** importiert (gleiches Git-Repo, jeweils anderes *Root Directory*):

| Projekt   | Root Directory   | Framework Preset | Environment Variables      |
| --------- | ---------------- | ---------------- | -------------------------- |
| api       | `apps/api`       | Other            | `JWT_SECRET`, `DATABASE_URL` |
| dashboard | `apps/dashboard` | Next.js          | `API_URL`, `WEB_URL`       |
| web       | `apps/web`       | Next.js          | `API_URL`, `DASHBOARD_URL` |

Reihenfolge:

1. **Datenbank anlegen:** im api-Projekt unter *Storage → Create Database* **Neon (Postgres)** wählen und mit dem Projekt verbinden. Vercel setzt `DATABASE_URL` (gepoolte Verbindung) dann automatisch. Alternativ jede andere Postgres-URL (Supabase, RDS, …) manuell als `DATABASE_URL` eintragen.
2. **api** deployen. Die Go-Funktion liegt in `apps/api/api/index.go`; `apps/api/vercel.json` leitet alle Pfade dorthin. `JWT_SECRET` setzen (`openssl rand -base64 32`). Die Tabellen werden beim ersten Request angelegt; `GET /api/health` prüft die DB-Verbindung.
3. **dashboard** deployen mit `API_URL=https://<api-projekt>.vercel.app` und `WEB_URL=https://<web-domain>`.
4. **web** deployen mit `API_URL=https://<api-projekt>.vercel.app` und `DASHBOARD_URL=https://<dashboard-projekt>.vercel.app`.

Danach nur die Domain des **web**-Projekts verwenden – Dashboard und API sind darüber unter `/dashboard` und `/api` erreichbar.

Tipp: In den Projekteinstellungen unter *Git → Ignored Build Step* `npx turbo-ignore` eintragen, damit nur betroffene Apps neu gebaut werden.
