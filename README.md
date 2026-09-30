# auth

Turborepo mit zwei Next.js-Apps und einem Go-Backend – jede App ist ein eigenes Vercel-Projekt.

```
apps/
  web/         Next.js 16 – Landing Page, /login, /register  (Port 3000)
  dashboard/   Next.js 16 – Dashboard unter /dashboard        (Port 3001)
  api/         Go – REST-API unter /api                       (Port 8080)
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

> ⚠️ User werden aktuell **im Speicher** gehalten (`apps/api/pkg/auth/store.go`). Auf Vercel (Serverless) gehen sie bei jedem Cold Start verloren und werden nicht zwischen Instanzen geteilt. Für den echten Betrieb das `UserStore`-Interface mit einer Datenbank (z. B. Postgres/Neon) implementieren.

## Lokale Entwicklung

Voraussetzungen: Node.js ≥ 20, pnpm 10, Go ≥ 1.24

```sh
pnpm install
pnpm dev          # startet web, dashboard und api parallel
```

Dann <http://localhost:3000> öffnen. Weitere Befehle:

```sh
pnpm build        # baut alle Apps
pnpm check-types  # tsc + go vet
pnpm test         # go test
```

## Deployment auf Vercel

Das Repo wird als **drei Vercel-Projekte** importiert (gleiches Git-Repo, jeweils anderes *Root Directory*):

| Projekt   | Root Directory   | Framework Preset | Environment Variables      |
| --------- | ---------------- | ---------------- | -------------------------- |
| api       | `apps/api`       | Other            | `JWT_SECRET`               |
| dashboard | `apps/dashboard` | Next.js          | `API_URL`, `WEB_URL`       |
| web       | `apps/web`       | Next.js          | `API_URL`, `DASHBOARD_URL` |

Reihenfolge:

1. **api** deployen. Die Go-Funktion liegt in `apps/api/api/index.go`; `apps/api/vercel.json` leitet alle Pfade dorthin. `JWT_SECRET` setzen (`openssl rand -base64 32`).
2. **dashboard** deployen mit `API_URL=https://<api-projekt>.vercel.app` und `WEB_URL=https://<web-domain>`.
3. **web** deployen mit `API_URL=https://<api-projekt>.vercel.app` und `DASHBOARD_URL=https://<dashboard-projekt>.vercel.app`.

Danach nur die Domain des **web**-Projekts verwenden – Dashboard und API sind darüber unter `/dashboard` und `/api` erreichbar.

Tipp: In den Projekteinstellungen unter *Git → Ignored Build Step* `npx turbo-ignore` eintragen, damit nur betroffene Apps neu gebaut werden.
