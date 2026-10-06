# scheduler

Web application in Go for automatic shift scheduling of store staff.
Prototype use case: cashiers, whose monthly schedule is currently made by hand.

## Status
Early development. Done so far:
- HTTP server on chi: request logging, panic recovery, per-IP rate limiting,
  cross-origin request protection
- Login with sessions (expired sessions are purged in the background)
- Server-rendered pages (html/template) with partial rendering for HTMX requests
- Database access via sqlx

Not implemented yet: schedule generation, employees and contract hours,
days-off requests, roles and invite-based registration, organization management.

## Planned
- Generate a monthly schedule from staffing needs per weekday/time of day,
  contract hours (full / 2/3 / half), limits on consecutive working days,
  and employees' requested days off
- Roles: platform admin, organization admin, manager, employee (invite links)
- Substitution suggestions when someone is absent

## Run
`go run .`
