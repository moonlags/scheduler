# scheduler

Web application in Go for automatic shift scheduling of store staff.
Prototype use case: cashiers, whose monthly schedule is currently made by hand.

## Status

Early development. Done so far:
- HTTP server on chi: request logging, panic recovery, per-IP rate limiting,
  cross-origin request protection
- Login with sessions (expired sessions are purged in the background)
- Server-rendered pages (html/template) with partial rendering for HTMX requests
- Database access via sqlx with embedded goose migrations
- Templates and static assets embedded into the binary

Not implemented yet: schedule generation, employees and contract hours,
days-off requests, roles and invite-based registration, organization management.

## Planned
- Generate a monthly schedule from staffing needs per weekday/time of day,
  contract hours (full / 2/3 / half), limits on consecutive working days,
  and employees' requested days off
- Roles: platform admin, organization admin, manager, employee (invite links)
- Substitution suggestions when someone is absent

## Layout

```
cmd/scheduler/main.go        entrypoint: config -> database -> HTTP server
internal/config/             environment configuration
internal/store/              database connection, queries, migrations
internal/store/migrations/   SQL migrations (embedded)
internal/web/                HTTP handlers, routing, middleware
internal/web/templates/      html/template pages and fragments (embedded)
internal/web/static/         CSS and JS assets (embedded)
```

The binary is self-contained: migrations, templates and static files are
compiled in with `go:embed`, so it only needs a database path at runtime.

## Documentation

Thesis material lives in `docs/`: the data model description in
`docs/explaining.txt`, the class diagram and the use case diagram in
`docs/diagrams/`.

## Run

Copy `.env.example` to `.env`, adjust it, then:

```
go run ./cmd/scheduler
```

`DB_URI` is required; `PORT`, `ADMIN_EMAIL` and `ADMIN_PASSWORD` have defaults.
The admin account is created (or its password updated) on startup.

## Test

```
go test ./...
```
