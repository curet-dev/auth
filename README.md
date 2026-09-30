# auth

Turborepo mit Astro-Landing-Page, React-Dashboard (Vite) und Go-Backend – zusammen als **ein** Vercel-Projekt deploybar ([Vercel Services](https://vercel.com/docs/services)).

```
apps/
  web/         Astro – Landing Page, /login, /register        (Port 3000)
  dashboard/   Vite + React – Dashboard unter /dashboard/     (Port 3001)
  api/         Go – REST-API unter /api                       (Port 8080)
               + Postgres (lokal via Docker Compose, auf Vercel via Neon)
packages/
  ui/                 geteilte React-Komponenten (Button, Logo)
  typescript-config/  geteilte tsconfigs
```

## Architektur

Alle Apps laufen unter **einer Domain**:

| Pfad           | App                                |
| -------------- | ---------------------------------- |
| `/api/*`       | Go-API (`apps/api`)                |
| `/dashboard/*` | Dashboard (`apps/dashboard`, Vite `base: "/dashboard/"`) |
| alles andere   | Landing Page (`apps/web`)          |

Auf Vercel übernimmt das die `vercel.json` im Repo-Root (Services + Rewrites), lokal der Dev-Proxy in `apps/web/astro.config.mjs`. Weil alles auf derselben Origin liegt, funktioniert das Session-Cookie (`httpOnly`-JWT, von der API gesetzt) ohne CORS oder Cross-Domain-Cookies.

- **web** ist eine statische Astro-Seite. Das Login-/Register-Formular ist eine React-Komponente (Astro Island, `client:load`) und nutzt – wie das Dashboard – die gemeinsamen Komponenten aus `packages/ui`.
- **dashboard** ist eine React-SPA. Sie prüft beim Laden per `GET /api/auth/me` die Session und leitet ohne Login auf `/login` um.

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

Das ganze Repo ist **ein Vercel-Projekt** mit drei [Services](https://vercel.com/docs/services) (Beta). Die `vercel.json` im Root beschreibt, welche App wo liegt und welche Pfade sie bekommt – Vercel baut jede App separat (Astro, Vite bzw. Go) und deployt alles zusammen unter einer URL.

```json
{
  "services": {
    "web": { "root": "apps/web/", "framework": "astro" },
    "dashboard": { "root": "apps/dashboard/", "framework": "vite" },
    "api": { "root": "apps/api/", "framework": "go" }
  },
  "redirects": [{ "source": "/dashboard", "destination": "/dashboard/", "permanent": false }],
  "rewrites": [
    { "source": "/api/(.*)", "destination": { "service": "api" } },
    { "source": "/dashboard/(.*)", "destination": { "service": "dashboard" } },
    { "source": "/(.*)", "destination": { "service": "web" } }
  ]
}
```

Das Dashboard wird nach `apps/dashboard/dist/dashboard/` gebaut, damit die Dateipfade den URL-Pfaden entsprechen. Bekommt es später clientseitiges Routing (z. B. React Router mit `/dashboard/settings`), braucht der Service zusätzlich einen SPA-Fallback auf `/dashboard/index.html`.

Einrichtung:

1. *Add New → Project* → Repo importieren. **Root Directory** bleibt das Repo-Root (`./`), **Framework Preset:** `Services`.
2. **Environment Variable** `JWT_SECRET` setzen (`openssl rand -base64 32`).
3. *Deploy*.
4. **Datenbank:** im Projekt unter *Storage → Create Database* **Neon (Postgres)** wählen und verbinden – Vercel setzt `DATABASE_URL` automatisch. Danach einmal *Redeploy*. `GET /api/health` prüft die DB-Verbindung; die Tabellen werden beim Start angelegt.

Eine neue App kommt dazu, indem man sie unter `apps/` anlegt, als Service in `vercel.json` einträgt und eine Rewrite-Regel für ihren Pfad ergänzt.

Das Zusammenspiel lässt sich lokal auch mit der Vercel CLI testen: `vercel dev -L` (braucht `DATABASE_URL` und `JWT_SECRET` in der Umgebung).

Die Go-API wird über das Go-Preset gebaut: Einstiegspunkt ist `apps/api/cmd/server/main.go`, der Server lauscht auf `PORT`.
