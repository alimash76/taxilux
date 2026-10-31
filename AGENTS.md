# AGENTS.md — Taxi Company Website (Netherlands)

> Read this entire file before doing anything. It is the single source of truth for scope, stack, phases and working rules.
> If you use Claude Code, copy this file to `CLAUDE.md` as well.

---

## 0. WORKING RULES (most important — read twice)

1. **Work in phases. Phases are strictly sequential.** The phases are defined in section 5.
2. **HARD STOP AT THE END OF EVERY PHASE.** When a phase is finished you must:
   - stop writing code,
   - produce the **Phase Report** (template in section 6),
   - wait.
3. **Never start the next phase until the human owner writes exactly `APPROVED PHASE N`** (N = the phase number you just reported). Silence, "looks good", "ok" or a thumbs-up is NOT approval. If unsure, ask.
4. If the owner requests changes, fix them **inside the current phase** and produce a new Phase Report. Do not move on.
5. **Do not build anything outside the current phase**, even if it seems helpful or "quick". Put ideas in `docs/BACKLOG.md` instead.
6. **Keep it simple.** Prefer the standard library and boring, well-known tools. Every dependency must be justified in one line in `docs/DECISIONS.md`. If a simpler approach exists, use it.
7. **Never invent business facts** (company name, phone, email, KvK/BTW numbers, address, prices, license numbers, service-area cities). Use clearly marked placeholders such as `{{COMPANY_PHONE}}` stored in ONE content/config file, and list every unresolved placeholder in each Phase Report.
8. **Ask instead of guessing** when a requirement is ambiguous. Ask at most the minimum number of questions, all in one message.
9. **Git discipline:** one branch per phase (`phase-N-short-name`), small meaningful commits. After approval the owner merges and tags `phase-N-done`.
10. **Everything must run from a clean checkout** with `docker compose up --build` (plus a documented `.env` copied from `.env.example`). No manual steps that are not in the README.
11. Code, comments, commit messages, docs: **English**. Website UI text: **Dutch (default) and English**. Dutch text you write must be flagged in the Phase Report as "needs native-speaker review".

---

## 1. PROJECT CONTEXT

A small taxi company in the Netherlands needs a public website and a simple system to receive and manage ride requests.

**Audience (design for them):**
- Elderly people with little IT knowledge.
- Business people and their secretaries.
- Both groups are busy and impatient. In the first seconds on the home page they must understand: **where we are, what we do, how to order.**

**Primary goal:** the **"Request a ride" (Rit aanvragen)** button.
- It is visible on the first screen of the home page AND reachable from every important page (persistent button/bar, thumb-reachable on mobile).
- Clicking it offers exactly two options:
  1. **Call us** — a big `tel:` link with the company phone number (for short trips / urgent requests).
  2. **Send a request by email** — opens the ride-request form (for long trips / planned trips; customer gets a confirmation).

**Design principles:** simple, minimal, calm, high contrast, large text (base ≥ 18px), large buttons (≥ 48px tap target), no carousels, no auto-playing anything, no jargon, no hidden menus, no animations beyond subtle state changes. Must work well on a cheap phone and on desktop.


---

## 2. SITEMAP / CONTENT (from the requirements file)

- **Home:** short intro, main services, service area, "Request a ride" button.
- **About us:** company introduction; Vision, Mission, Goals; legal & financial details (company registration number = **KvK-nummer**, VAT number = **BTW-nummer**; tax code if applicable).
- **Services:** service area (geographic coverage) + services overview. *(The source file has 3 empty bullets here — propose sensible items in the Phase 1 report, do not silently invent them.)*
- **Trip types:** city/local trips; airport transfers; pre-booked trips; business/corporate (contract) trips; intercity/long distance.
- **How to book:** direct phone call for short/urgent trips; email request + confirmation for long/planned trips.
- **Feedback & complaints:**
  - General terms & conditions for passenger transport
  - Privacy statement
  - Complaints & suggestions handling overview
  - Complaint-handling procedure
  - **Complaint form** (goes to company email; automatic acknowledgement of receipt to the sender)
  - **Suggestion form** (goes to company email)
- **Contact.**

---

## 3. TECH STACK (decided — do not change without owner approval)

| Part | Choice | Notes |
|---|---|---|
| Backend | **Go** (latest stable) | `net/http` with Go 1.22+ routing; optionally `chi`. Standard project layout: `cmd/server`, `internal/...`. No heavy frameworks, no ORM. |
| Database | **SQLite** | Pure-Go driver `modernc.org/sqlite` (no CGO → simple static Docker image). `database/sql` + plain SQL. Embedded migrations (e.g. `goose`). Pragmas: `journal_mode=WAL`, `foreign_keys=ON`, `busy_timeout`. Single backend instance. |
| Frontend (public site) | **Astro** (static output) + minimal vanilla JS | Chosen for easy i18n (NL/EN) and shared layout (persistent CTA). Served by **nginx** in its own container. Zero client JS except forms and the ride-request dialog. If you believe something simpler is better, propose it in Phase 0 — do not switch silently. |
| Admin panel | Server-rendered Go `html/template` under `/admin` | Part of the backend container. Keep it plain and fast. |
| Reverse proxy / TLS | **Caddy** container | Automatic HTTPS in production. Routes `/api/*` and `/admin/*` → backend, everything else → frontend. |
| Containers | **Docker + Docker Compose** | Separate images for `backend` and `frontend`; `caddy` as third service; `mailpit` (dev only) for catching emails. Multi-stage builds, non-root users, small images, healthchecks. |
| Email | SMTP via env config | Provider must be EU-based/GDPR-compliant. Dev uses Mailpit. Email sending sits behind a small interface. |
| Backups | SQLite online backup (`VACUUM INTO` or `.backup`) | Scheduled, rotated, stored on a separate volume. |

**Repository layout (create in Phase 0):**

```
/
├── AGENTS.md
├── README.md
├── docker-compose.yml
├── docker-compose.dev.yml
├── Makefile
├── .env.example
├── backend/            # Go
├── frontend/           # Astro
├── caddy/              # Caddyfile(s)
├── docs/               # DECISIONS.md, BACKLOG.md, RUNBOOK.md, requirements docx
└── content/            # site.config (company data/placeholders), i18n strings
```

---

## 4. NON-FUNCTIONAL REQUIREMENTS (apply to every phase)

- **Accessibility:** WCAG 2.2 AA minimum. Semantic HTML, visible focus, full keyboard use, labels on all fields, error messages in plain language next to the field, colour is never the only signal, contrast ≥ 4.5:1 (prefer 7:1 for body text), respects `prefers-reduced-motion`.
- **Mobile-first:** perfect at 360px width.
- **Performance:** fast on slow connections; no large images; no third-party scripts, fonts or trackers by default (self-host fonts).
- **Privacy / GDPR (AVG):** collect only necessary data; EU hosting; no analytics or cookies unless owner approves (then add consent banner); retention policy; privacy statement consistent with what the system really does.
- **Security:** HTTPS, secure headers, server-side validation, parameterized SQL, rate limiting, spam protection, no secrets in the repo.
- **Language:** Dutch default (`/nl/...`), English (`/en/...`), with a visible, simple language switch. `hreflang` set.
- **Testing:** backend unit/integration tests for all handlers and DB code; `go vet` and `gofmt` clean; a smoke test script that checks the main flows against the running containers.

---

## 5. PHASES

Each phase has: **Goal**, **Tasks**, **Out of scope**, **Acceptance checklist** (the human will verify these), then **STOP**.

### Phase 0 — Foundations

**Goal:** a working skeleton that runs in Docker, with decisions documented.

**Tasks**
- Create the repository layout (section 3), `README.md` (how to run, how to test), `.env.example`, `Makefile` (`up`, `down`, `test`, `lint`, `logs`).
- Backend: minimal Go server with `GET /api/health` (checks DB connectivity), structured logging (`log/slog`), config from env, graceful shutdown, SQLite opened with the required pragmas, migration runner with an initial empty migration.
- Frontend: Astro project with one placeholder page in NL and EN, built and served by nginx.
- Caddy routing, `docker-compose.yml` (prod-like) and `docker-compose.dev.yml` (with Mailpit and hot-reload if simple).
- `docs/DECISIONS.md` (each dependency + one-line reason), `docs/BACKLOG.md`.
- `content/site.config` with all business placeholders listed in one place.

**Out of scope:** real pages, forms, emails, admin.

**Acceptance checklist**
- [ ] `docker compose up --build` works from a clean clone.
- [ ] `http://localhost/` shows the placeholder page; `/api/health` returns OK with DB status.
- [ ] Backend and frontend run in separate containers, non-root, with healthchecks.
- [ ] `make test` and `make lint` pass.
- [ ] README is enough for a newcomer to run the project.

### Phase 1 — Public website (frontend only, no real submissions)

**Goal:** the complete public site, polished, accessible, bilingual — forms are present but not yet connected.

**Tasks**
- All pages from section 2, NL + EN, with shared layout, header, footer.
- Home page: intro, main services, service area, big "Request a ride" button above the fold.
- Persistent "Request a ride" button/bar on all pages; opens a dialog/page with the two options (call / email form). Phone number is a `tel:` link.
- Form pages (ride request, complaint, suggestion, contact) with full UI, client-side validation and plain-language errors — submit shows a "not connected yet" notice for now.
- Legal pages (terms, privacy, complaint procedure) with **clearly marked placeholder text** to be replaced by the owner/legal advisor.
- Design system in CSS variables (colours, type scale, spacing); light, calm palette (propose 2 palette options in the report).
- Proposal for the 3 empty "Services" bullets.

**Out of scope:** backend endpoints, database writes, emails, admin, SEO extras.

**Acceptance checklist**
- [ ] Opening the home page on a 360px phone: I immediately see who/where/what + the Request button without scrolling.
- [ ] The Request button is present on every page and offers call/email.
- [ ] NL ⇄ EN switch works on every page.
- [ ] Keyboard-only navigation works; focus is always visible.
- [ ] Lighthouse (mobile): Accessibility ≥ 95, Best Practices ≥ 95, Performance ≥ 90.
- [ ] No third-party requests (check the browser network tab).
- [ ] Every placeholder is listed in the report.

### Phase 2 — Backend API + SQLite (forms become real, saved to the database)

**Goal:** submissions are validated, stored, and acknowledged on screen.

**Tasks**
- Migrations and tables (initial model below).
- Endpoints (JSON): `POST /api/ride-requests`, `POST /api/complaints`, `POST /api/suggestions`, `POST /api/contact`.
- Server-side validation (required fields, lengths, email/phone format, date not in the past for ride requests, privacy consent checkbox required).
- Spam protection: honeypot field + per-IP rate limiting (`golang.org/x/time/rate` or equivalent). Optional CAPTCHA only if owner approves.
- Human-friendly **reference numbers** (e.g. `RR-2026-000123`, `CP-2026-000045`).
- Frontend forms wired to the API; success page shows the reference number; error states handled gracefully (network down, validation errors, rate limited).
- Tests for every handler and the DB layer.

**Initial data model (adjust only with owner approval)**
- `ride_requests`: id, reference, created_at, locale, name, phone, email, pickup_address, destination, pickup_at, trip_type (`city|airport|prebooked|corporate|intercity`), passengers, notes, privacy_consent, status (`new|confirmed|assigned|completed|cancelled`), internal_notes, updated_at.
- `complaints`: id, reference, created_at, locale, name, email, phone (optional), ride_reference (optional), message, privacy_consent, status (`new|in_progress|resolved`), internal_notes, updated_at.
- `suggestions`: id, created_at, locale, name (optional), email (optional), message, privacy_consent.
- `contact_messages`: id, created_at, locale, name, email, message, privacy_consent.
- `email_outbox` (used in Phase 3): id, created_at, to, subject, body, status, attempts, last_error, sent_at.
- `admin_users`, `sessions` (used in Phase 4).

**Out of scope:** sending emails, admin panel.

**Acceptance checklist**
- [ ] Submitting each of the 4 forms creates a row in SQLite (show me how to inspect it).
- [ ] Invalid input shows clear Dutch/English messages next to the field.
- [ ] Bot test: a filled honeypot or >N requests/minute is rejected without losing legitimate users.
- [ ] Reference number is shown to the user after success.
- [ ] `make test` passes; coverage for handlers is meaningful.
- [ ] Data survives `docker compose down && up` (volume works).

### Phase 3 — Email notifications

**Goal:** the business is notified of every submission; customers get confirmations. **A submission must never be lost because email failed.**

**Tasks**
- Outbox pattern: handler saves submission + outbox row in one transaction; a background worker sends emails with retry/backoff.
- Emails: (a) to the business — new ride request / complaint / suggestion / contact message with all details; (b) to the customer — acknowledgement with reference number (ride request, complaint; for suggestions/contact optional short thanks).
- Plain-text + simple HTML versions, NL/EN according to the customer's locale.
- Configurable recipient addresses per form type via env.
- Mailpit in dev; documentation for production SMTP setup including **SPF, DKIM, DMARC** DNS records (written as a checklist for the owner).
- Protect against email header injection.

**Out of scope:** admin panel, status-change emails (Phase 4).

**Acceptance checklist**
- [ ] Each form triggers the business email and the customer acknowledgement (visible in Mailpit).
- [ ] If SMTP is down, submissions are still saved and emails are sent later automatically.
- [ ] Emails render well in plain text and HTML, in NL and EN.
- [ ] The DNS/deliverability checklist is clear enough for a non-developer.

### Phase 4 — Admin panel

**Goal:** the owner/dispatcher can manage requests without touching the database.

**Tasks**
- Login at `/admin` with email + password (argon2id), server-side sessions in SQLite, secure cookies, CSRF protection, login rate limiting, **TOTP 2FA**.
- Roles: `owner` (everything, including user management) and `dispatcher` (ride requests and complaints only).
- Ride requests: list with filters (today / new / status / date range / search), detail page, change status, internal notes, **email to customer on confirm/cancel** (editable template text).
- Complaints: list, detail, status change, internal notes.
- Suggestions and contact messages: read-only lists.
- CSV export for ride requests and complaints (date range).
- CLI command to create the first owner account (no default passwords).
- Admin UI is simple, large-text, works on phone.

**Out of scope:** drivers, vehicles, payments, live tracking.

**Acceptance checklist**
- [ ] I can create the owner via CLI, log in, enable 2FA.
- [ ] A dispatcher cannot access owner-only pages (try direct URLs).
- [ ] Changing a ride status updates the DB and sends the right email.
- [ ] CSV export opens correctly in Excel/LibreOffice (UTF-8, correct delimiters).
- [ ] Unauthenticated access to any `/admin` page redirects to login.

### Phase 5 — Security, GDPR and operations hardening

**Goal:** safe to expose to the internet and compliant in practice.

**Tasks**
- Security headers (CSP, HSTS, X-Content-Type-Options, Referrer-Policy, Permissions-Policy), request size limits, timeouts.
- Data retention: automatic deletion/anonymisation of old ride requests/complaints after configurable periods (propose defaults; owner decides). Admin function to delete/export all data of one person (data-subject request).
- Privacy statement text updated to match real processing (draft for legal review, flagged as such).
- Cookie/consent: confirm no non-essential cookies; if any were added, add consent banner.
- Scheduled SQLite backups with rotation + **a tested restore procedure** documented in `docs/RUNBOOK.md`.
- Structured logs without personal data in plain logs; log rotation.
- Dependency vulnerability scan (`govulncheck`, `npm audit`) and fix/justify findings.
- Basic uptime/health monitoring instructions.

**Acceptance checklist**
- [ ] Security header scan (e.g. securityheaders.com or local equivalent) shows A-level.
- [ ] I performed a backup → delete data → restore drill following the RUNBOOK and it worked.
- [ ] Retention job demonstrated with test data.
- [ ] No personal data in application logs (spot check).
- [ ] `govulncheck` / `npm audit` clean or exceptions documented.

### Phase 6 — SEO, production deployment and handover

**Goal:** live, findable, maintainable.

**Tasks**
- Local SEO: unique titles/descriptions, `LocalBusiness`/`TaxiService` JSON-LD, `sitemap.xml`, `robots.txt`, canonical + hreflang, Open Graph basics, clean URLs.
- Production Compose + Caddy with automatic HTTPS; deployment guide for a small EU VPS (step-by-step for a non-expert), update procedure, rollback procedure.
- Final Lighthouse + accessibility pass; real-device checklist.
- `docs/RUNBOOK.md` finished: start/stop, logs, backups, restore, updating, creating users, common problems.
- Owner handover document: where content lives and how to edit it (phone, address, texts).

**Acceptance checklist**
- [ ] Site reachable over HTTPS on the real domain (or a staging domain).
- [ ] Structured data validates; sitemap lists all NL/EN pages.
- [ ] Following only the deployment guide, a fresh server can be set up.
- [ ] Updating content (e.g. phone number) is documented and works.

### Backlog (do NOT start; only if owner explicitly asks)

SMS/Telegram/WhatsApp alert for urgent requests, corporate clients module, drivers & vehicles, address autocomplete, online payment, live tracking, mobile app, analytics (with consent).

---

## 6. PHASE REPORT TEMPLATE (output this at the end of every phase, then STOP)

```
# Phase N Report — <name>

## Summary
(3–6 lines: what was built)

## How to run & verify
(exact commands + URLs)

## Acceptance checklist
(copy the phase checklist; mark each item PASS / FAIL / NOT TESTED with a short note)

## Decisions made
(anything not dictated by AGENTS.md, one line each; also added to docs/DECISIONS.md)

## Unresolved placeholders / information needed from the owner
(list)

## Known issues / limitations
(list, honest)

## Dutch text needing native review
(list of files/keys)

## Next
Waiting for: APPROVED PHASE N
```

**After sending this report: do nothing else until the owner replies.**

---

## 7. CODING CONVENTIONS

- Go: `gofmt`, `go vet`, small packages under `internal/`, context-aware handlers, errors wrapped with context, no global state, config via env, `log/slog`.
- SQL: plain, parameterized, migrations only (never edit an applied migration).
- Frontend: semantic HTML first, CSS variables, no CSS framework unless justified in DECISIONS.md, i18n strings in one structured place (no hard-coded text in components).
- Secrets only via environment variables; `.env` is git-ignored; `.env.example` documents everything.
- Every phase leaves `main` in a working, runnable state.

---

## 8. HOW THE OWNER STARTS YOU

The owner will send: **"Read AGENTS.md. Start Phase 0 only."**
Begin Phase 0, finish it, send the Phase 0 Report, and wait for `APPROVED PHASE 0`.