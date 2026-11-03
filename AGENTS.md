# AGENTS.md

Quick context for AI coding agents working on this repository. Read this before
changing code.

## What this is

Diploma project: a web application in Go for automatic shift scheduling of store
staff (prototype use case: cashiers). The written thesis is in Latvian; code,
comments and commits are in English. The user communicates in English.

Keep it simple: standard library + `chi` + `html/template` + `sqlx` + SQLite.
No ORM, no frontend framework, no JS build step. Prefer extending existing
patterns over introducing new layers/packages.

Deliverables are a working prototype **and** the written documentation, which
must follow the supervisor's methodology in
`docs/documentation_requirements.md`.

## Current state

Works: config from env, SQLite via sqlx, embedded goose migrations, embedded
templates/static, session login, background session purge, rate limiting,
cross-origin protection, index/login/dashboard/logout pages.

Not implemented yet: schedule generation (`internal/schedule` does not exist
yet), employees UI routing, invite-based registration, profile management,
organizations/branches, the four-role model, days-off/sick-leave requests,
history/audit log, email notifications.

Migration `00001_init.sql` already creates `organizations` and `roles`, but no
Go code uses them yet. `users.is_admin` is a placeholder for the real role
model. `internal/web/employee.go` + `employees.html` exist but are **not
routed**; the `employees` table only has `user_id`/`rate` and requires an
existing user.

## Product scope (from `docs/documentation.md`)

The thesis defines the target product; the code is only its authentication
skeleton. Read the spec before implementing a feature.

Roles (hierarchy: Darbinieks < Vecākais darbinieks < Uzņēmuma administrators <
Sistēmas administrators; a higher role can do the lower role's actions within
its own scope, and that check is enforced server-side per request):

- **Darbinieks** (employee): own branch schedule (read-only), own day-off
  requests and sick leave, own profile. Only employees have a workload
  (`slodze`, 0..1).
- **Vecākais darbinieks** (senior employee): one branch — employees, schedule
  settings, generate/edit schedule, approve day-offs, pick substitutes,
  branch edit history. Seniors are not planned into the schedule.
- **Uzņēmuma administrators** (company admin): every branch of one company,
  company audit log. Not tied to a branch, not in schedules.
- **Sistēmas administrators** (system admin): all companies, global audit log,
  no company membership.

Registration is invite-only: system admin invites a company admin, company
admin invites a senior, senior invites an employee. A registration link is
single-use, valid 7 days, and carries role + company/branch (+ workload for
employees). There is no public sign-up. Changing an e-mail uses a separate
24 h confirmation link.

Domain (full entity list in `docs/documentation.md` §1.5): companies → branches
→ memberships (role + workload); branch schedule settings (working hours per
weekday, staffing intervals, max shift length, min rest); monthly schedules
(with a version for optimistic concurrency); shift records (employee or free
shift; status planned/free); day-off requests (pending/approved/rejected);
sick leave (accepted as fact, no reason stored); audit log.

Scheduling rules (§1.2.8, §1.2.9, §1.2.12): month norm = working days × 8 h ×
workload; at most one continuous shift per employee per day; shift inside
branch hours and ≤ max length; min rest between consecutive shifts; planned
hours must not exceed the norm (manual edit only warns); prefer employees with
fewer planned hours; approved day-offs and sick leave are never scheduled.
Hour summary and shortage report are computed from shift records + settings on
every read, never stored. Approving a day-off or filing sick leave marks the
affected future shifts as free (employee cleared, date/time kept). Free shifts
are shown only to seniors/admins and count as unfilled; substitutes are branch
employees (not seniors) without a shift/day-off/sick leave that day, sorted by
planned hours; picking one re-runs the same checks.

Cross-cutting rules: every multi-table write in one transaction; the audit
entry is written in the same transaction as the action; optimistic concurrency
(schedule version → "changed since opened" error, no save); validate everything
server-side; dates `DD.MM.YYYY`, times `HH:MM`; shifts never cross midnight;
shift times in Latvian local time, system event times in UTC; never store
day-off/sick-leave reasons in the audit log or e-mails; bcrypt/Argon2 passwords,
HttpOnly + Secure + SameSite cookies, CSRF protection; SMTP e-mails are
best-effort and must not block or roll back the action.

## Commands

Go is installed but **not on PATH** in the bash tool on this machine:

```bash
export PATH="/c/Program Files/Go/bin:$PATH"
```

Then:

```bash
gofmt -l .                 # should be empty
go build ./...
go vet ./...
go test ./...
go run ./cmd/scheduler
```

Smoke test against a throwaway database:

```bash
DB_URI=/tmp/dev.db PORT=3999 ADMIN_EMAIL=root ADMIN_PASSWORD=root go run ./cmd/scheduler
```

`.env` is optional (`godotenv.Load`); `DB_URI` is the only required variable.
Default admin credentials are `root`/`root` and are logged at startup.

## Layout

```
cmd/scheduler/main.go          entrypoint: config -> store -> web -> run
internal/config/config.go      env config (Config struct, Load)
internal/store/store.go        Store type, connect, embedded goose migrations
internal/store/migrations/     SQL migrations (embedded via go:embed)
internal/store/user.go         User model + user/session-owner queries, EnsureAdmin
internal/store/session.go      session save/delete/purge
internal/store/employee.go     Employee model + queries
internal/web/server.go         Server, New, routes, middleware, Run
internal/web/render.go         //go:embed templates static, template loading/rendering
internal/web/handlers.go       index/dashboard handlers, redirect helper
internal/web/auth.go           login/logout, auth middleware, userFrom(ctx)
internal/web/session.go        cookie + token creation, SHA-256 hashing
internal/web/employee.go       employee handlers (not routed yet)
internal/web/templates/        layout, pages/, fragments/ (embedded)
internal/web/static/           bulma.min.css, htmx.min.js (embedded)
docs/documentation.md          thesis draft (converted from .docx)
docs/documentation_requirements.md  supervisor requirements (converted from .pdf)
docs/media/                    images extracted during conversion
docs/diagrams/                 uml_diagramm.drawio, usecase_diagramm.drawio
```

Request flow: `main` → `config.Load` → `store.Open` (runs migrations) →
`store.EnsureAdmin` → `web.New` (loads templates, builds routes) →
`Server.Run`.

Auth flow: login sets cookie `sid` = random token; DB stores only its SHA-256
hash in `sessions.token`. `withUser` puts `*store.User` in the request context;
`requireUser`/`requireAdmin` guard route groups; handlers read it via
`userFrom(r.Context())`.

Templates: each page is parsed together with `templates/layout.html` and all
fragments, then executed via `{{define "content"}}`. Use `renderPage` for full
pages and `render` for HTMX fragments. Cookies are `Secure` only on HTTPS (or
`X-Forwarded-Proto: https`) so local http development works.

## Database rules

- SQLite driver: `modernc.org/sqlite` (pure Go, no cgo).
- Connections add `foreign_keys(1)` and `busy_timeout(5000)` pragmas.
- Migrations are embedded from `internal/store/migrations` and applied on
  startup. Naming: `000NN_description.sql`, with `-- +goose Up` / `-- +goose Down`.
- Never edit an already-applied migration; add a new numbered one.
- All SQL lives in `internal/store`; handlers call store methods and only deal
  with HTTP. Store methods return raw/`fmt.Errorf`-wrapped errors.
- `users.id` is a ULID string; `sessions.token` is a BLOB; timestamps are Unix
  ints (`expires_at`) or SQLite `datetime('now')` text (`created_at`).
- History survives deletion: audit entries (and shift/request records) keep a
  copy of the display name and only reference IDs, intentionally without FK
  integrity, so deleting a user/company/branch does not lose history. Cascades
  apply to operational data only (schedules, settings, sessions, invites).
- New migrations/templates/static files are picked up automatically by the
  existing `go:embed` directives, but require a rebuild.

## Conventions

- Handlers are methods on `*web.Server`, unexported, named `handleX`.
- UI text/forms use HTMX attributes; return a fragment for HTMX requests.
- Errors: log with `log.Println("handleX:", err)` then
  `http.Error(w, "internal error", http.StatusInternalServerError)`.
- Redirects go through `redirect` (handles `HX-Redirect`).
- Keep store methods small, one query per method; reuse model structs with
  `db:"..."` tags.
- Tests: `internal/store/store_test.go` opens a temp-file DB; follow that
  pattern. Add table-driven tests for scheduling rules once they exist.
- Validation belongs in the store/handler layer (not only the template) and
  mirrors the rule list in `docs/documentation.md` §1.2; user-facing messages
  there are Latvian, keep them consistent.
- Multi-table changes (generating a schedule, approving a request, deleting a
  company/branch, audit log) must use a single `sqlx` transaction.

## Known gaps / things to revisit

- `EnsureAdmin` re-hashes and overwrites the admin password on every startup
  (so DB-side password changes are lost). Revisit when profile password change
  is implemented.
- `users.is_admin` predates the four-role model; migration + code must move to
  `roles`/`organizations`/branch memberships when the first non-admin feature
  is built.
- No CSRF tokens beyond `http.NewCrossOriginProtection`; no email sending.
- Route table intentionally excludes features without handlers; add route and
  handler together.
- Session purge runs hourly inside `Server.Run`.
- `go.mod` module is `github.com/moonlags/ikeaScheduler` (capital S) — imports
  must match exactly.
- `docs/documentation.md` §1.2.15 duplicates the 1.2.13 title in the source; it
  covers the role access matrix (table 3).

## Docs and thesis material

Sources (keep the originals; the Markdown files are regenerable):

- `docs/documentation.docx` → `docs/documentation.md` — the thesis draft:
  chapters 1–3 (task statement, requirements, SWOT, OO data model) and
  appendices 1–2. This is the requirement source for behaviour and naming.
  The draft currently stops at §3.1 (OO conceptual model); the ER model, table
  structure and realisation/testing chapters are still to be written.
- `docs/documentation_requirements.pdf` → `docs/documentation_requirements.md`
  — supervisor's methodology: required thesis structure, content of every
  section, and acceptance rules. Note: the PDF starts at printed page 6, so
  chapter 1 is missing from this source.
- `docs/media/` — images extracted from the two sources.
- `docs/diagrams/uml_diagramm.drawio` — UML class diagram (appendix 2).
- `docs/diagrams/usecase_diagramm.drawio` — use case diagram (appendix 1).
  Keep it review-clean: one solid system boundary, no colours/legend, actors
  outside; independent goals are separate use cases with plain associations
  (not «extend»); only optional behaviour uses «extend»; mandatory behaviour
  uses «include»; actor generalization is a solid line with a hollow triangle.

What the methodology expects, in short: Titullapa, Saturs, Ievads, Uzdevuma
nostādne (input/output data, functional requirements FR x.y.z, non-functional
requirements, users + use case diagram, planned DB content, user scenarios),
Sistēmas koncepcijas izvērtējums (SWOT), DB design (OO conceptual model,
Chen-notation ER diagram, table/relationship structure, 3NF check), realisation
and testing (SQL script, integrity checks, every FR tested with at least a
positive, a boundary and a negative case), conclusions, sources, appendices.
Documentation numbering follows that structure; requirement IDs in
`documentation.md` (§1.2.x) are stable and should be referenced in code review
and tests.

Converting new material:

- `.docx`: `unzip -o file.docx -d /tmp/docx`, then read headings/styles from
  `word/document.xml` (`Heading1`–`Heading3`, `w:numPr` for lists, `w:tbl` for
  tables); images are in `word/media/`. A small Node script is enough; Node is
  available, Python is not (only the Windows Store stub).
- `.pdf`: use poppler's `pdftotext -enc UTF-8 -layout file.pdf out.txt`. UTF-8
  must be requested explicitly, otherwise Latvian diacritics are dropped.
  `-layout` keeps reading order (default mode and `-raw` scramble text boxes).
  Convert to Markdown next to the source; tables/code samples in this material
  end up as `text` fenced blocks, which is intentional.
- `pandoc` is not installed. Put converted Markdown next to the source under
  `docs/` and mention the source file in its header.
- Thesis text is Latvian; keep Latvian terms (Filiāle, Darbinieks, Vecākais
  darbinieks, Uzņēmums, Slodze, Maiņa, Grafiks, Darbību vēsture,
  Reģistrācijas saite) when discussing requirements, but keep Go identifiers
  English.
