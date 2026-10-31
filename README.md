# TaxiLux

Taxi company website and ride-request system for the Netherlands. Public site (Astro), Go backend with SQLite, Caddy reverse proxy.

## Stack

- Backend: Go (`net/http`, SQLite via `modernc.org/sqlite`)
- Frontend: Astro (static), served by nginx
- Proxy: Caddy
- Orchestration: Docker Compose

## Run (production-like)

```bash
cp .env.example .env
# edit .env with your values (all business data are placeholders)
docker compose up --build
```

- Site: http://localhost/ (Caddy → frontend)
- API health: http://localhost/api/health (Caddy → backend)
- Mailpit (dev only): http://localhost:8025/

## Run (development)

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```
Backend runs with `go run` (source mounted), frontend with `astro dev`, plus Mailpit on port 8025.

## Test

```bash
make test
```
- Backend: `go vet` + `go test` (inside the backend container)
- Frontend: `npm ci` + `npm run build`

## Lint

```bash
make lint
```
- Backend: `gofmt` + `go vet`
- Frontend: `astro check`

## Logs

```bash
make logs
```

## Repository layout

See AGENTS.md section 3. Business placeholders live in `content/site.config` — never invent real company data.
